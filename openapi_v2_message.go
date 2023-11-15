package nano

import (
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type MessageTypeV2 int

const (
	MessageTypeV2Text MessageTypeV2 = iota
	MessageTypeV2TextImage
	MessageTypeV2Markdown
	MessageTypeV2Ark
	MessageTypeV2Embed
)

func (mt2 MessageTypeV2) String() string {
	switch mt2 {
	case MessageTypeV2Text:
		return "文本"
	case MessageTypeV2TextImage:
		return "图文混排"
	case MessageTypeV2Markdown:
		return "MD"
	case MessageTypeV2Ark:
		return "模版"
	case MessageTypeV2Embed:
		return "嵌入"
	default:
		return "未知类型" + strconv.Itoa(int(mt2))
	}
}

type MessageV2 struct {
	Author struct {
		UserOpenID   string `json:"user_openid"`
		MemberOpenID string `json:"member_openid"`
	} `json:"author"`
	Content     string              `json:"content"`
	ID          string              `json:"id"`
	GroupOpenID string              `json:"group_openid"`
	Timestamp   time.Time           `json:"timestamp"`
	Attachments []MessageAttachment `json:"attachments"`
}

// MessagePostV2 V2 发消息结构体
//
// https://bot.q.qq.com/wiki/develop/api-231017/server-inter/message/send-receive/send.html
type MessagePostV2 struct {
	Type           MessageTypeV2 `json:"msg_type"`
	Seq            int           `json:"msg_seq,omitempty"` // 回复消息的序号，与 msg_id 联合使用，避免相同消息id回复重复发送，不填默认是1。相同的 msg_id + msg_seq 重复发送会失败。
	Content        string        `json:"content,omitempty"`
	ReplyEventID   string        `json:"event_id,omitempty"` // 前置收到的事件ID，用于发送被动消息
	ReplyMessageID string        `json:"msg_id,omitempty"`

	// image		否	【暂不支持】
	MessageReference *MessageReference `json:"message_reference,omitempty"` // 【暂未支持】消息引用

	Markdown *MessageMarkdown `json:"markdown,omitempty"`
	KeyBoard *MessageKeyboard `json:"keyboard,omitempty"`
	Ark      *MessageArk      `json:"ark,omitempty"`
	Embed    *MessageEmbed    `json:"embed,omitempty"`
}

func (mp *MessagePostV2) String() string {
	sb := strings.Builder{}
	sb.WriteString("[v2.")
	sb.WriteString(mp.Type.String())
	sb.WriteString(".")
	sb.WriteString(strconv.Itoa(mp.Seq))
	sb.WriteString("]")
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
	/*if mp.Image != "" {
		sb.WriteString(", 图片URL: ")
		sb.WriteString(mp.Image)
	}
	if mp.ImageFile != "" {
		sb.WriteString(", 图片内容: ")
		x := mp.ImageFile
		if len(x) > 64 {
			x = x[:64] + "..."
		}
		sb.WriteString(x)
	}
	if len(mp.ImageBytes) > 0 {
		sb.WriteString(", 图片大小: ")
		sb.WriteString(strconv.Itoa(len(mp.ImageBytes)))
	}*/
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

func (bot *Bot) postV2MessageTo(ep string, content *MessagePostV2) (*IDTimestampMessageResult, error) {
	return bot.postOpenAPIofIDTimestampMessageResult(ep, "", WriteBodyFromJSON(content))
}

// PostMessageToQQUser 向 openid 指定的用户发送消息
//
// https://bot.q.qq.com/wiki/develop/api-231017/server-inter/message/send-receive/send.html#%E5%8D%95%E8%81%8A
func (bot *Bot) PostMessageToQQUser(id string, content *MessagePostV2) (*IDTimestampMessageResult, error) {
	logrus.Infoln(getLogHeader(), "<= [Q]单:", id+",", content)
	return bot.postV2MessageTo("/v2/users/"+id+"/messages", content)
}

// PostMessageToQQGroup 向 openid 指定的群发送消息
//
// https://bot.q.qq.com/wiki/develop/api-231017/server-inter/message/send-receive/send.html#%E7%BE%A4%E8%81%8A
func (bot *Bot) PostMessageToQQGroup(id string, content *MessagePostV2) (*IDTimestampMessageResult, error) {
	logrus.Infoln(getLogHeader(), "<= [Q]群:", id+",", content)
	return bot.postV2MessageTo("/v2/groups/"+id+"/messages", content)
}
