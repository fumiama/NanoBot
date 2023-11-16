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

type MessageSegmentType int

const (
	MessageSegmentTypeText MessageSegmentType = iota
	MessageSegmentTypeImage
	MessageSegmentTypeImageBytes
	MessageSegmentTypeReply
	MessageSegmentTypeAudio
	MessageSegmentTypeVideo
)

// Message impl the array form of message
type Messages []MessageSegment

// MessageSegment impl the single message
// MessageSegment 消息数组
type MessageSegment struct {
	Type MessageSegmentType
	Data string
}

// String impls the interface fmt.Stringer
func (m MessageSegment) String() string {
	return m.Data
}

// Text 纯文本
func Text(text ...interface{}) MessageSegment {
	return MessageSegment{
		Type: MessageSegmentTypeText,
		Data: HideURL(MessageEscape(fmt.Sprint(text...))),
	}
}

// Face QQ表情
// https://bot.q.qq.com/wiki/develop/api/openapi/message/message_format.html#%E6%94%AF%E6%8C%81%E7%9A%84%E6%A0%BC%E5%BC%8F
func Face(id int) MessageSegment {
	return MessageSegment{
		Type: MessageSegmentTypeText,
		Data: "<emoji:" + strconv.Itoa(id) + ">",
	}
}

// Image 普通图片
func Image(file string) MessageSegment {
	return MessageSegment{
		Type: MessageSegmentTypeImage,
		Data: file,
	}
}

// ImageBytes 普通图片
func ImageBytes(data []byte) MessageSegment {
	return MessageSegment{
		Type: MessageSegmentTypeImageBytes,
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
		Type: MessageSegmentTypeText,
		Data: "<@!" + id + ">",
	}
}

// AtAll @全体成员
// https://bot.q.qq.com/wiki/develop/api/openapi/message/message_format.html#%E6%94%AF%E6%8C%81%E7%9A%84%E6%A0%BC%E5%BC%8F
func AtAll() MessageSegment {
	return MessageSegment{
		Type: MessageSegmentTypeText,
		Data: "@everyone",
	}
}

// AtChannel #频道
// https://bot.q.qq.com/wiki/develop/api/openapi/message/message_format.html#%E6%94%AF%E6%8C%81%E7%9A%84%E6%A0%BC%E5%BC%8F
func AtChannel(id string) MessageSegment {
	return MessageSegment{
		Type: MessageSegmentTypeText,
		Data: "<#channel_id>",
	}
}

// Record QQ 语音
// https://bot.q.qq.com/wiki/develop/api-231017/server-inter/message/send-receive/rich-text-media.html
func Record(url string) MessageSegment {
	return MessageSegment{
		Type: MessageSegmentTypeAudio,
		Data: url,
	}
}

// Video QQ 视频
// https://bot.q.qq.com/wiki/develop/api-231017/server-inter/message/send-receive/rich-text-media.html
func Video(url string) MessageSegment {
	return MessageSegment{
		Type: MessageSegmentTypeVideo,
		Data: url,
	}
}

// Reply 回复
// https://github.com/botuniverse/onebot-11/tree/master/message/segment.md#%E5%9B%9E%E5%A4%8D
func ReplyTo(id string) MessageSegment {
	return MessageSegment{
		Type: MessageSegmentTypeReply,
		Data: id,
	}
}

// ReplyWithMessage returns a reply message
func ReplyWithMessage(messageID string, m ...MessageSegment) Messages {
	return append(Messages{ReplyTo(messageID)}, m...)
}
