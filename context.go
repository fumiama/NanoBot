package nano

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

//go:generate go run codegen/context/main.go

type Ctx struct {
	Event
	State
	Message *Message
	IsToMe  bool
	IsQQ    bool

	caller *Bot
	ma     *Matcher
}

// decoder 反射获取的数据
type decoder []dec

type dec struct {
	index int
	key   string
}

// decoder 缓存
var decoderCache = sync.Map{}

// Parse 将 Ctx.State 映射到结构体
func (ctx *Ctx) Parse(model interface{}) (err error) {
	var (
		rv       = reflect.ValueOf(model).Elem()
		t        = rv.Type()
		modelDec decoder
	)
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parse state error: %v", r)
		}
	}()
	d, ok := decoderCache.Load(t)
	if ok {
		modelDec = d.(decoder)
	} else {
		modelDec = decoder{}
		for i := 0; i < t.NumField(); i++ {
			t1 := t.Field(i)
			if key, ok := t1.Tag.Lookup("zero"); ok {
				modelDec = append(modelDec, dec{
					index: i,
					key:   key,
				})
			}
		}
		decoderCache.Store(t, modelDec)
	}
	for _, d := range modelDec { // decoder类型非小内存，无法被编译器优化为快速拷贝
		rv.Field(d.index).Set(reflect.ValueOf(ctx.State[d.key]))
	}
	return nil
}

// CheckSession 判断会话连续性
func (ctx *Ctx) CheckSession() Rule {
	msg := ctx.Value.(*Message)
	return func(ctx2 *Ctx) bool {
		msg2, ok := ctx.Value.(*Message)
		if !ok || msg.Author == nil || msg2.Author == nil { // 确保无空
			return false
		}
		return msg.Author.ID == msg2.Author.ID && msg.ChannelID == msg2.ChannelID
	}
}

// Send 发送一批消息
func (ctx *Ctx) Send(messages Messages) (m []*Message, err error) {
	isnextreply := false
	textlist := []any{}
	var reply *Message
	for _, msg := range messages {
		switch msg.Type {
		case MessageSegmentTypeText:
			textlist = append(textlist, msg.Data)
		case MessageSegmentTypeImage:
			reply, err = ctx.SendImage(msg.Data, isnextreply, textlist...)
			if isnextreply {
				isnextreply = false
			}
			textlist = textlist[:0]
			m = append(m, reply)
			if err != nil {
				return
			}
		case MessageSegmentTypeImageBytes:
			if ctx.IsQQ {
				continue
			}
			reply, err = ctx.SendImageBytes(StringToBytes(msg.Data), isnextreply, textlist...)
			if isnextreply {
				isnextreply = false
			}
			textlist = textlist[:0]
			m = append(m, reply)
			if err != nil {
				return
			}
		case MessageSegmentTypeReply:
			isnextreply = true
		case MessageSegmentTypeAudio, MessageSegmentTypeVideo:
			if !ctx.IsQQ {
				continue
			}
			fp := &FilePost{
				URL: msg.Data,
			}
			if msg.Type == MessageSegmentTypeAudio {
				fp.Type = FileTypeAudio
			} else if msg.Type == MessageSegmentTypeVideo {
				fp.Type = FileTypeVideo
			}
			if OnlyQQGroup(ctx) {
				reply, err = ctx.PostFileToQQGroup(ctx.Message.ChannelID, fp)
			} else if OnlyQQPrivate(ctx) {
				reply, err = ctx.PostFileToQQUser(ctx.Message.Author.ID, fp)
			}
			m = append(m, reply)
			if err != nil {
				return
			}
		}
	}
	if len(textlist) > 0 {
		reply, err = ctx.SendPlainMessage(isnextreply, textlist...)
		m = append(m, reply)
	}
	return
}

// SendChain 链式发送
func (ctx *Ctx) SendChain(message ...MessageSegment) (m []*Message, err error) {
	return ctx.Send(message)
}

// Post 发送消息到对方
func (ctx *Ctx) Post(replytosender bool, post *MessagePost) (reply *Message, err error) {
	msg := ctx.Message
	if msg != nil {
		post.ReplyMessageID = msg.ID
		if OnlyGuild(ctx) && replytosender {
			post.MessageReference = &MessageReference{
				MessageID: msg.ID,
			}
		}
	} else {
		post.ReplyMessageID = "MESSAGE_CREATE"
	}

	if OnlyDirect(ctx) { // dms
		reply, err = ctx.PostMessageToUser(msg.GuildID, post)
	} else if OnlyChannel(ctx) {
		reply, err = ctx.PostMessageToChannel(msg.ChannelID, post)
	} else { // v2
		switch {
		case post.Markdown != nil:
			post.Type = MessageTypeMarkdown
		case post.Ark != nil:
			post.Type = MessageTypeArk
		case post.Embed != nil:
			post.Type = MessageTypeEmbed
		default:
			post.Type = MessageTypeText
		}
		post.Seq = len(GetTriggeredMessages(msg.ID)) + 1
		if OnlyQQGroup(ctx) {
			reply, err = ctx.PostMessageToQQGroup(msg.ChannelID, post)
		} else if OnlyQQPrivate(ctx) {
			reply, err = ctx.PostMessageToQQUser(msg.ChannelID, post)
		}
	}
	if err != nil && msg != nil && reply != nil && reply.ID != "" {
		logtriggeredmessages(msg.ID, reply.ID)
	}
	return
}

// SendPlainMessage 发送纯文本消息到对方
func (ctx *Ctx) SendPlainMessage(replytosender bool, printable ...any) (*Message, error) {
	return ctx.Post(replytosender, &MessagePost{
		Content: HideURL(fmt.Sprint(printable...)),
	})
}

// SendImage 发送带图片消息到对方
func (ctx *Ctx) SendImage(file string, replytosender bool, caption ...any) (reply *Message, err error) {
	if OnlyQQ(ctx) {
		fp := &FilePost{
			Type: FileTypeImage,
			URL:  file,
		}
		if len(caption) > 0 {
			_, _ = ctx.SendPlainMessage(replytosender, caption...)
		}
		if OnlyQQGroup(ctx) {
			reply, err = ctx.PostFileToQQGroup(ctx.Message.ChannelID, fp)
		} else if OnlyQQPrivate(ctx) {
			reply, err = ctx.PostFileToQQUser(ctx.Message.Author.ID, fp)
		}
		return
	}

	post := &MessagePost{
		Content: HideURL(fmt.Sprint(caption...)),
	}

	if strings.HasPrefix(file, "http") {
		post.Image = file
	} else {
		post.ImageFile = file
	}

	return ctx.Post(replytosender, post)
}

// SendImageBytes 发送带图片消息到对方
func (ctx *Ctx) SendImageBytes(data []byte, replytosender bool, caption ...any) (*Message, error) {
	if OnlyQQ(ctx) {
		return nil, errors.New("QQ暂不支持直接发送图片数据")
	}

	post := &MessagePost{
		Content: HideURL(fmt.Sprint(caption...)),
	}

	post.ImageBytes = data

	return ctx.Post(replytosender, post)
}

// Echo 向自身分发虚拟事件
func (ctx *Ctx) Echo(payload *WebsocketPayload) {
	ctx.caller.processEvent(payload)
}

// FutureEvent ...
func (ctx *Ctx) FutureEvent(Type string, rule ...Rule) *FutureEvent {
	return ctx.ma.FutureEvent(Type, rule...)
}

// Get 从 promt 获得回复
func (ctx *Ctx) Get(prompt string) string {
	if prompt != "" {
		_, _ = ctx.SendPlainMessage(false, prompt)
	}
	return (<-ctx.FutureEvent("Message", ctx.CheckSession()).Next()).Event.Value.(*Message).Content
}

// ExtractPlainText 提取消息中的纯文本
func (ctx *Ctx) ExtractPlainText() string {
	if ctx == nil || ctx.Value == nil {
		return ""
	}
	if msg, ok := ctx.Value.(*Message); ok {
		return msg.Content
	}
	return ""
}

// MessageString 字符串消息便于Regex
func (ctx *Ctx) MessageString() string {
	return ctx.ExtractPlainText()
}

// Block 匹配成功后阻止后续触发
func (ctx *Ctx) Block() {
	ctx.ma.SetBlock(true)
}

// Block 在 pre, rules, mid 阶段阻止后续触发
func (ctx *Ctx) Break() {
	ctx.ma.Break = true
}

// GroupID 唯一的发送者所属组 ID
func (ctx *Ctx) GroupID() uint64 {
	grp := uint64(0)
	if ctx.IsQQ {
		if OnlyQQGroup(ctx) {
			grp = DigestID(ctx.Message.ChannelID)
		} else if OnlyQQPrivate(ctx) {
			grp = DigestID(ctx.Message.Author.ID)
		} else {
			return 0
		}
	} else {
		var err error
		grp, err = strconv.ParseUint(ctx.Message.ChannelID, 10, 64)
		if err != nil {
			return 0
		}
	}
	return grp
}

// GroupID 唯一的发送者 ID
func (ctx *Ctx) UserID() uint64 {
	if ctx.IsQQ {
		return DigestID(ctx.Message.Author.ID)
	}
	grp, _ := strconv.ParseUint(ctx.Message.Author.ID, 10, 64)
	return grp
}
