package nano

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/FloatTech/ttl"
)

var (
	triggeredMessages   = ttl.NewCache[string, []string](time.Minute * 5)
	triggeredMessagesMu = sync.Mutex{}
)

func logtriggeredmessages(id, reply string) {
	triggeredMessagesMu.Lock()
	defer triggeredMessagesMu.Unlock()
	triggeredMessages.Set(id, append(triggeredMessages.Get(id), reply))
}

// GetTriggeredMessages 获取被 id 消息触发的回复消息 id
func GetTriggeredMessages(id string) []string {
	triggeredMessagesMu.Lock()
	defer triggeredMessagesMu.Unlock()
	return triggeredMessages.Get(id)
}

type MessageType int

const (
	MessageTypeText MessageType = iota
	MessageTypeImage
	MessageTypeImageBytes
	MessageTypeReply
)

// Message impl the array form of message
type Messages []MessageSegment

// MessageSegment impl the single message
// MessageSegment 消息数组
type MessageSegment struct {
	Type MessageType
	Data string
}

// String impls the interface fmt.Stringer
func (m MessageSegment) String() string {
	return m.Data
}

// Text 纯文本
func Text(text ...interface{}) MessageSegment {
	return MessageSegment{
		Type: MessageTypeText,
		Data: MessageEscape(fmt.Sprint(text...)),
	}
}

// Face QQ表情
// https://bot.q.qq.com/wiki/develop/api/openapi/message/message_format.html#%E6%94%AF%E6%8C%81%E7%9A%84%E6%A0%BC%E5%BC%8F
func Face(id int) MessageSegment {
	return MessageSegment{
		Type: MessageTypeText,
		Data: "<emoji:" + strconv.Itoa(id) + ">",
	}
}

// Image 普通图片
func Image(file string) MessageSegment {
	return MessageSegment{
		Type: MessageTypeImage,
		Data: file,
	}
}

// ImageBytes 普通图片
func ImageBytes(data []byte) MessageSegment {
	return MessageSegment{
		Type: MessageTypeImageBytes,
		Data: BytesToString(data),
	}
}

// At @某人
// https://bot.q.qq.com/wiki/develop/api/openapi/message/message_format.html#%E6%94%AF%E6%8C%81%E7%9A%84%E6%A0%BC%E5%BC%8F
func At(id string) MessageSegment {
	if id == "all" {
		return AtAll()
	}
	return MessageSegment{
		Type: MessageTypeText,
		Data: "<@!" + id + ">",
	}
}

// AtAll @全体成员
// https://bot.q.qq.com/wiki/develop/api/openapi/message/message_format.html#%E6%94%AF%E6%8C%81%E7%9A%84%E6%A0%BC%E5%BC%8F
func AtAll() MessageSegment {
	return MessageSegment{
		Type: MessageTypeText,
		Data: "@everyone",
	}
}

// AtChannel #频道
// https://bot.q.qq.com/wiki/develop/api/openapi/message/message_format.html#%E6%94%AF%E6%8C%81%E7%9A%84%E6%A0%BC%E5%BC%8F
func AtChannel(id string) MessageSegment {
	return MessageSegment{
		Type: MessageTypeText,
		Data: "<#channel_id>",
	}
}

// Reply 回复
// https://github.com/botuniverse/onebot-11/tree/master/message/segment.md#%E5%9B%9E%E5%A4%8D
func ReplyTo(id string) MessageSegment {
	return MessageSegment{
		Type: MessageTypeReply,
		Data: id,
	}
}

// ReplyWithMessage returns a reply message
func ReplyWithMessage(messageID string, m ...MessageSegment) Messages {
	return append(Messages{ReplyTo(messageID)}, m...)
}
