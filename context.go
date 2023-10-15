package nano

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type Ctx struct {
	Event
	State
	Caller  *Bot
	Message *Message
	ma      *Matcher
	IsToMe  bool
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

// Send 发送消息到对方
func (ctx *Ctx) Send(replytosender bool, post *MessagePost) (*Message, error) {
	msg := ctx.Value.(*Message)
	post.ReplyMessageID = msg.ID
	if replytosender {
		post.MessageReference = &MessageReference{
			MessageID: msg.ID,
		}
	}
	return ctx.Caller.PostMessageToChannel(msg.ChannelID, post)
}

// SendPlainMessage 发送纯文本消息到对方
func (ctx *Ctx) SendPlainMessage(replytosender bool, printable ...any) (*Message, error) {
	msg := ctx.Value.(*Message)
	post := &MessagePost{
		ReplyMessageID: msg.ID,
	}
	if replytosender {
		post.MessageReference = &MessageReference{
			MessageID: msg.ID,
		}
	}
	post.Content = fmt.Sprint(printable...)
	return ctx.Caller.PostMessageToChannel(msg.ChannelID, post)
}

// SendImage 发送带图片消息到对方
func (ctx *Ctx) SendImage(file string, replytosender bool, caption ...any) (*Message, error) {
	msg := ctx.Value.(*Message)
	post := &MessagePost{
		ReplyMessageID: msg.ID,
	}
	if strings.HasPrefix(file, "http") {
		post.Image = file
	} else {
		post.ImageFile = file
	}
	if replytosender {
		post.MessageReference = &MessageReference{
			MessageID: msg.ID,
		}
	}
	post.Content = fmt.Sprint(caption...)
	return ctx.Caller.PostMessageToChannel(msg.ChannelID, post)
}

// Block 匹配成功后阻止后续触发
func (ctx *Ctx) Block() {
	ctx.ma.SetBlock(true)
}

// Block 在 pre, rules, mid 阶段阻止后续触发
func (ctx *Ctx) Break() {
	ctx.ma.Break = true
}
