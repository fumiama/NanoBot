package nano

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
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
	URL string `json:"url"`
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
	AuditID string `json:"audit_id"`
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
	Embed            *MessageEmbed     `json:"embed,omitempty"`
	Ark              *MessageArk       `json:"ark,omitempty"`
	MessageReference *MessageReference `json:"message_reference,omitempty"`
	Image            string            `json:"image,omitempty"`
	ImageFile        string            `json:"-"` // ImageFile 为图片路径 file:/// or base64:// or base16384:// , 与 Image 参数二选一, 优先 Image
	ReplyMessageID   string            `json:"msg_id,omitempty"`
	ReplyEventID     string            `json:"event_id,omitempty"`
	Markdown         *MessageMarkdown  `json:"markdown,omitempty"`
	KeyBoard         *MessageKeyboard  `json:"keyboard,omitempty"`
}

// MessageEscape 消息转义
//
// https://bot.q.qq.com/wiki/develop/api/openapi/message/message_format.html
func MessageEscape(text string) string {
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	return text
}

// MessageUnescape 消息解转义
//
// https://bot.q.qq.com/wiki/develop/api/openapi/message/message_format.html
func MessageUnescape(text string) string {
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	return text
}

// PostMessageToChannel 向 channel_id 指定的子频道发送消息
//
// https://bot.q.qq.com/wiki/develop/api/openapi/message/post_messages.html
func (bot *Bot) PostMessageToChannel(id string, content *MessagePost) (*Message, error) {
	if content.ImageFile == "" {
		return bot.postOpenAPIofMessage("/channels/"+id+"/messages", "", WriteBodyFromJSON(content))
	}
	x := reflect.ValueOf(content).Elem()
	t := x.Type()
	msg := []any{}
	for i := 0; i < x.NumField(); i++ {
		xi := x.Field(i)
		if xi.IsZero() {
			continue
		}
		tag := t.Field(i).Tag.Get("json")
		if tag == "-" {
			tag = "file_image"
		}
		msg = append(msg, tag)
		if xi.Kind() == reflect.Struct {
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
	return bot.postOpenAPIofMessage("/channels/"+id+"/messages", contenttype, body)
}

// DeleteMessageInChannel 回子频道 channel_id 下的消息 message_id
//
// https://bot.q.qq.com/wiki/develop/api/openapi/message/delete_message.html
func (bot *Bot) DeleteMessageInChannel(channelid, messageid string, hidetip bool) error {
	return bot.DeleteOpenAPI(WriteHTTPQueryIfNotNil("/channels/"+channelid+"/messages/"+messageid,
		"hidetip", hidetip,
	), "", nil)
}
