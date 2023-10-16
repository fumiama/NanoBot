package nano

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	ErrEmptyMessagePost = errors.New("empty message post")
)

// Message 消息对象
//
// https://bot.q.qq.com/wiki/develop/api/openapi/message/model.html#%E6%B6%88%E6%81%AF%E5%AF%B9%E8%B1%A1-message
type Message struct {
	ID               string              `json:"id"`
	ChannelID        string              `json:"channel_id"`
	GuildID          string              `json:"guild_id"`
	Content          string              `json:"content"`
	Timestamp        time.Time           `json:"timestamp"`
	EditedTimestamp  time.Time           `json:"edited_timestamp"`
	MentionEveryone  bool                `json:"mention_everyone"`
	Author           *User               `json:"author"`
	Attachments      []MessageAttachment `json:"attachments"`
	Embeds           []MessageEmbed      `json:"embeds"`
	Member           *Member             `json:"member"`
	Ark              *MessageArk         `json:"ark"`
	SeqInChannel     string              `json:"seq_in_channel"`
	MessageReference *MessageReference   `json:"message_reference"`
	SrcGuildID       string              `json:"src_guild_id"`
	Data             *struct {
		MessageAudit *MessageAudited `json:"message_audit,omitempty"`
	} `json:"data,omitempty"`
}

// "=> ｷﾞﾙﾄﾞ:", ctx.Message.GuildID+", 频道:", ctx.Message.ChannelID+", 用户:", ctx.Message.Author.Username+"("+ctx.Message.Author.ID+"), 内容:", ctx.Message.Content
func (m *Message) String() string {
	sb := strings.Builder{}
	sb.WriteByte('[')
	sb.WriteString(m.SeqInChannel)
	sb.WriteByte(']')
	sb.WriteString(" ｷﾞﾙﾄﾞ: ")
	sb.WriteString(m.GuildID)
	if m.SrcGuildID != "" {
		sb.WriteString(", 元ｷﾞﾙﾄﾞ: ")
		sb.WriteString(m.SrcGuildID)
	}
	sb.WriteString(", 频道: ")
	sb.WriteString(m.ChannelID)
	if m.Author != nil {
		sb.WriteString(", 用户: ")
		sb.WriteString(m.Author.Username)
		sb.WriteByte('(')
		sb.WriteString(m.Author.ID)
		sb.WriteByte(')')
		if m.Author.Bot {
			sb.WriteString("(机器人)")
		}
	} else {
		sb.WriteString(", 用户: 未知")
	}
	if m.Content == "" {
		sb.WriteString(", 无文本")
	} else {
		sb.WriteString(", 文本: ")
		if m.MentionEveryone {
			sb.WriteString("[@全体]")
		}
		sb.WriteString(m.Content)
	}
	if len(m.Attachments) > 0 {
		sb.WriteString(", 附加: ")
		for _, a := range m.Attachments {
			sb.WriteString("<ID:")
			sb.WriteString(a.ID)
			sb.WriteString(",URL:")
			sb.WriteString(a.URL)
			sb.WriteByte('>')
		}
	}
	if len(m.Embeds) > 0 {
		for _, e := range m.Embeds {
			sb.WriteString(", 嵌入: <标题:")
			sb.WriteString(e.Title)
			sb.WriteString(",提示:")
			sb.WriteString(e.Prompt)
			sb.WriteByte('>')
		}
	}
	if m.Ark != nil {
		sb.WriteString(", 模版: ")
		sb.WriteString(strconv.Itoa(m.Ark.TemplateID))
	}
	if m.MessageReference != nil {
		sb.WriteString(", 回复: ")
		sb.WriteString(m.MessageReference.MessageID)
	}
	return sb.String()
}

// MessageEmbed https://bot.q.qq.com/wiki/develop/api/openapi/message/model.html#messageembed
type MessageEmbed struct {
	Title     string                 `json:"title"`
	Prompt    string                 `json:"prompt"`
	Thumbnail *MessageEmbedThumbnail `json:"thumbnail"`
	Fields    []MessageEmbedField    `json:"fields"`
}

// MessageEmbedThumbnail https://bot.q.qq.com/wiki/develop/api/openapi/message/model.html#messageembedthumbnail
type MessageEmbedThumbnail struct {
	URL string `json:"url"`
}

// MessageEmbedField https://bot.q.qq.com/wiki/develop/api/openapi/message/model.html#messageembedfield
type MessageEmbedField struct {
	Name string `json:"name"`
}

// MessageAttachment https://bot.q.qq.com/wiki/develop/api/openapi/message/model.html#messageattachment
type MessageAttachment struct {
	ContentType string `json:"content_type,omitempty"`
	Filename    string `json:"filename,omitempty"`
	Height      int    `json:"height,omitempty"`
	ID          string `json:"id,omitempty"`
	Size        int    `json:"size,omitempty"`
	URL         string `json:"url,omitempty"`
	Width       int    `json:"width,omitempty"`
}

// MessageArk https://bot.q.qq.com/wiki/develop/api/openapi/message/model.html#messageark
type MessageArk struct {
	TemplateID int            `json:"template_id"`
	KV         []MessageArkKV `json:"kv"`
}

// MessageArkKV https://bot.q.qq.com/wiki/develop/api/openapi/message/model.html#messagearkkv
type MessageArkKV struct {
	Key   string          `json:"key"`
	Value string          `json:"value"`
	Obj   []MessageArkObj `json:"obj"`
}

// MessageArkObj https://bot.q.qq.com/wiki/develop/api/openapi/message/model.html#messagearkobj
type MessageArkObj struct {
	ObjKV []MessageArkObjKV `json:"obj_kv"`
}

// MessageArkObjKV https://bot.q.qq.com/wiki/develop/api/openapi/message/model.html#messagearkobjkv
type MessageArkObjKV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// MessageReference https://bot.q.qq.com/wiki/develop/api/openapi/message/model.html#messagereference
type MessageReference struct {
	MessageID             string `json:"message_id"`
	IgnoreGetMessageError bool   `json:"ignore_get_message_error"`
}

// MessageDelete https://bot.q.qq.com/wiki/develop/api/openapi/message/model.html#messagedelete
type MessageDelete struct {
	Message *Message `json:"message"`
	OpUser  *User    `json:"op_user"`
}

// MessageAudited 消息审核对象
//
// https://bot.q.qq.com/wiki/develop/api/openapi/message/model.html#%E6%B6%88%E6%81%AF%E5%AE%A1%E6%A0%B8%E5%AF%B9%E8%B1%A1-messageaudited
type MessageAudited struct {
	AuditID    string    `json:"audit_id"`
	AuditTime  time.Time `json:"audit_time"`
	ChannelID  string    `json:"channel_id"`
	CreateTime time.Time `json:"create_time"`
	GuildID    string    `json:"guild_id"`
	MessageID  string    `json:"message_id"`
}

// GetMessageFromChannel 获取子频道 channel_id 下的消息 message_id 的详情
//
// https://bot.q.qq.com/wiki/develop/api/openapi/message/get_message_of_id.html
func (bot *Bot) GetMessageFromChannel(messageid, channelid string) (*Message, error) {
	return bot.getOpenAPIofMessage("/channels/" + channelid + "/messages/" + messageid)
}

// MessagePost 发送消息所需参数
//
// https://bot.q.qq.com/wiki/develop/api/openapi/message/post_messages.html#%E9%80%9A%E7%94%A8%E5%8F%82%E6%95%B0
type MessagePost struct {
	Content          string            `json:"content,omitempty"`
	Embed            *MessageEmbed     `json:"embed,omitempty"` // https://bot.q.qq.com/wiki/develop/api/openapi/message/template/embed_message.html
	Ark              *MessageArk       `json:"ark,omitempty"`   // https://bot.q.qq.com/wiki/develop/api/openapi/message/message_template.html
	MessageReference *MessageReference `json:"message_reference,omitempty"`
	Image            string            `json:"image,omitempty"`
	ImageFile        string            `json:"-"` // ImageFile 为图片路径 file:/// or base64:// or base16384:// , 与 Image 参数二选一, 优先 Image
	ReplyMessageID   string            `json:"msg_id,omitempty"`
	ReplyEventID     string            `json:"event_id,omitempty"`
	Markdown         *MessageMarkdown  `json:"markdown,omitempty"`
	KeyBoard         *MessageKeyboard  `json:"keyboard,omitempty"`
}

func (mp *MessagePost) String() string {
	sb := strings.Builder{}
	if mp.Content == "" {
		sb.WriteString("无文本")
	} else {
		sb.WriteString("文本: ")
		sb.WriteString(mp.Content)
	}
	if mp.ReplyMessageID != "" {
		sb.WriteString(", 回应消息: ")
		sb.WriteString(mp.ReplyMessageID)
	}
	if mp.ReplyEventID != "" {
		sb.WriteString(", 回应事件: ")
		sb.WriteString(mp.ReplyEventID)
	}
	if mp.Embed != nil {
		sb.WriteString(", 嵌入: <标题:")
		sb.WriteString(mp.Embed.Title)
		sb.WriteString(",提示:")
		sb.WriteString(mp.Embed.Prompt)
		sb.WriteByte('>')
	}
	if mp.Ark != nil {
		sb.WriteString(", 模版: ")
		sb.WriteString(strconv.Itoa(mp.Ark.TemplateID))
	}
	if mp.MessageReference != nil {
		sb.WriteString(", 回复: ")
		sb.WriteString(mp.MessageReference.MessageID)
	}
	if mp.Image != "" {
		sb.WriteString(", 图片URL:")
		sb.WriteString(mp.Image)
	}
	if mp.ImageFile != "" {
		sb.WriteString(", 图片内容:")
		x := mp.ImageFile
		if len(x) > 64 {
			x = x[:64] + "..."
		}
		sb.WriteString(x)
	}
	if mp.Markdown != nil {
		sb.WriteString(", MD模版: ")
		sb.WriteString(strconv.Itoa(mp.Markdown.TemplateID))
	}
	if mp.KeyBoard != nil {
		sb.WriteString(", KB模版: ")
		sb.WriteString(mp.KeyBoard.ID)
	}
	return sb.String()
}

func (bot *Bot) postMessageTo(ep string, content *MessagePost) (*Message, error) {
	if content.ImageFile == "" {
		return bot.postOpenAPIofMessage(ep, "", WriteBodyFromJSON(content))
	}
	x := reflect.ValueOf(content).Elem()
	t := x.Type()
	msg := []any{}
	for i := 0; i < x.NumField(); i++ {
		xi := x.Field(i)
		if xi.IsZero() {
			continue
		}
		tag, _, _ := strings.Cut(t.Field(i).Tag.Get("json"), ",")
		if tag == "-" {
			tag = "file_image"
		}
		msg = append(msg, tag)
		if xi.Kind() == reflect.Pointer && xi.Elem().Kind() == reflect.Struct {
			data, err := json.Marshal(xi.Interface())
			if err != nil {
				return nil, err
			}
			msg = append(msg, data)
			continue
		}
		msg = append(msg, xi.String())
	}
	if len(msg) < 2 {
		return nil, ErrEmptyMessagePost
	}
	body, contenttype, err := WriteBodyByMultipartFormData(msg...)
	if err != nil {
		return nil, errors.Wrap(err, getThisFuncName())
	}
	return bot.postOpenAPIofMessage(ep, contenttype, body)
}

// PostMessageToChannel 向 channel_id 指定的子频道发送消息
//
// https://bot.q.qq.com/wiki/develop/api/openapi/message/post_messages.html
func (bot *Bot) PostMessageToChannel(id string, content *MessagePost) (*Message, error) {
	logrus.Infoln(getLogHeader(), "<= 频道:", id+",", content)
	return bot.postMessageTo("/channels/"+id+"/messages", content)
}

// DeleteMessageInChannel 回子频道 channel_id 下的消息 message_id
//
// https://bot.q.qq.com/wiki/develop/api/openapi/message/delete_message.html
func (bot *Bot) DeleteMessageInChannel(channelid, messageid string, hidetip bool) error {
	logrus.Infoln(getLogHeader(), "<x 频道:", channelid+", 消息:", messageid)
	return bot.DeleteOpenAPI(WriteHTTPQueryIfNotNil("/channels/"+channelid+"/messages/"+messageid,
		"hidetip", hidetip,
	), "", nil)
}

// MessageSetting 频道消息频率设置对象
//
// https://bot.q.qq.com/wiki/develop/api/openapi/setting/model.html
type MessageSetting struct {
	DisableCreateDm   bool     `json:"disable_create_dm"`
	DisablePushMsg    bool     `json:"disable_push_msg"`
	ChannelIDs        []string `json:"channel_ids"`
	ChannelPushMaxNum uint32   `json:"channel_push_max_num"`
}

// GetGuildMessageSetting 获取机器人在频道 guild_id 内的消息频率设置
//
// https://bot.q.qq.com/wiki/develop/api/openapi/setting/message_setting.html
func (bot *Bot) GetGuildMessageSetting(id string) (*MessageSetting, error) {
	return bot.getOpenAPIofMessageSetting("/guilds/" + id + "/message/setting")
}
