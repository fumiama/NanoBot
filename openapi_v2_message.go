package nano

import (
	"strconv"

	"github.com/sirupsen/logrus"
)

type MessageType int

const (
	MessageTypeText MessageType = iota
	MessageTypeTextImage
	MessageTypeMarkdown
	MessageTypeArk
	MessageTypeEmbed
)

func (mt2 MessageType) String() string {
	switch mt2 {
	case MessageTypeText:
		return "文本"
	case MessageTypeTextImage:
		return "图文混排"
	case MessageTypeMarkdown:
		return "MD"
	case MessageTypeArk:
		return "模版"
	case MessageTypeEmbed:
		return "嵌入"
	default:
		return "未知类型" + strconv.Itoa(int(mt2))
	}
}

// PostMessageToQQUser 向 openid 指定的用户发送消息
//
// https://bot.q.qq.com/wiki/develop/api-231017/server-inter/message/send-receive/send.html#%E5%8D%95%E8%81%8A
func (bot *Bot) PostMessageToQQUser(id string, content *MessagePost) (*Message, error) {
	logrus.Infoln(getLogHeader(), "<= [Q]单:", id+",", content)
	return bot.postMessageTo("/v2/users/"+id+"/messages", content)
}

// PostMessageToQQGroup 向 openid 指定的群发送消息
//
// https://bot.q.qq.com/wiki/develop/api-231017/server-inter/message/send-receive/send.html#%E7%BE%A4%E8%81%8A
func (bot *Bot) PostMessageToQQGroup(id string, content *MessagePost) (*Message, error) {
	logrus.Infoln(getLogHeader(), "<= [Q]群:", id+",", content)
	return bot.postMessageTo("/v2/groups/"+id+"/messages", content)
}
