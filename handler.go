package nano

import (
	"reflect"
	"unsafe"
)

// generalHandleType 作为通用的 handler 函数调用约定使用
type generalHandleType func(uint32, *Bot, unsafe.Pointer)

// eventHandlerType 一个事件函数调用的必须信息
type eventHandlerType struct {
	h generalHandleType
	t reflect.Type
}

var types map[string]reflect.Type // types 便于反射初始化的 types

func init() {
	h := reflect.ValueOf(&Handler{}).Elem()
	t := h.Type()
	types = make(map[string]reflect.Type, h.NumField()*4)
	for i := 0; i < h.NumField(); i++ {
		tp := t.Field(i).Name[2:] // skip On
		types[tp] = t.Field(i).Type.In(2).Elem()
	}
}

// Handler 事件订阅
//
// https://bot.q.qq.com/wiki/develop/api/gateway/intents.html
type Handler struct {
	// GUILDS (1 << 0)

	OnGuildCreate   func(s uint32, bot *Bot, d *Guild)
	OnGuildUpdate   func(s uint32, bot *Bot, d *Guild)
	OnGuildDelete   func(s uint32, bot *Bot, d *Guild)
	OnChannelCreate func(s uint32, bot *Bot, d *Channel)
	OnChannelUpdate func(s uint32, bot *Bot, d *Channel)
	OnChannelDelete func(s uint32, bot *Bot, d *Channel)
	// GUILD_MEMBERS (1 << 1)

	OnGuildMemberAdd    func(s uint32, bot *Bot, d *Member)
	OnGuildMemberUpdate func(s uint32, bot *Bot, d *Member)
	OnGuildMemberRemove func(s uint32, bot *Bot, d *Member)
	// GUILD_MESSAGES (1 << 9)    // 消息事件，仅 *私域* 机器人能够设置此 intents。

	OnMessageCreate func(s uint32, bot *Bot, d *Message)
	OnMessageDelete func(s uint32, bot *Bot, d *MessageDelete)
	// GUILD_MESSAGE_REACTIONS (1 << 10)

	OnMessageReactionAdd    func(s uint32, bot *Bot, d *MessageReaction)
	OnMessageReactionRemove func(s uint32, bot *Bot, d *MessageReaction)
	// DIRECT_MESSAGE (1 << 12)

	OnDirectMessageCreate func(s uint32, bot *Bot, d *Message)
	OnDirectMessageDelete func(s uint32, bot *Bot, d *MessageDelete)
	// OPEN_FORUMS_EVENT (1 << 18)      // 论坛事件, 此为公域的论坛事件

	OnOpenForumThreadCreate func(s uint32, bot *Bot, d *Thread)
	OnOpenForumThreadUpdate func(s uint32, bot *Bot, d *Thread)
	OnOpenForumThreadDelete func(s uint32, bot *Bot, d *Thread)
	OnOpenForumPostCreate   func(s uint32, bot *Bot, d *Post)
	OnOpenForumPostDelete   func(s uint32, bot *Bot, d *Post)
	OnOpenForumReplyCreate  func(s uint32, bot *Bot, d *Reply)
	OnOpenForumReplyDelete  func(s uint32, bot *Bot, d *Reply)
	// AUDIO_OR_LIVE_CHANNEL_MEMBER (1 << 19)  // 音视频/直播子频道成员进出事件

	OnAudioOrLiveChannelMemberEnter func(s uint32, bot *Bot, d *AudioLiveChannelUsersChange)
	OnAudioOrLiveChannelMemberExit  func(s uint32, bot *Bot, d *AudioLiveChannelUsersChange)
	// QQ (1<<25) QQ 的一堆事件

	OnC2cMessageCreate     func(s uint32, bot *Bot, d *MessageV2)
	OnGroupAtMessageCreate func(s uint32, bot *Bot, d *MessageV2)
	OnGroupAddRobot        func(s uint32, bot *Bot, d *QQRobotStatus)
	OnGroupDelRobot        func(s uint32, bot *Bot, d *QQRobotStatus)
	OnGroupMsgReject       func(s uint32, bot *Bot, d *QQRobotStatus)
	OnGroupMsgReceive      func(s uint32, bot *Bot, d *QQRobotStatus)
	OnFriendAdd            func(s uint32, bot *Bot, d *QQRobotStatus)
	OnFriendDel            func(s uint32, bot *Bot, d *QQRobotStatus)
	OnC2cMsgReject         func(s uint32, bot *Bot, d *QQRobotStatus)
	OnC2cMsgReceive        func(s uint32, bot *Bot, d *QQRobotStatus)

	// INTERACTION (1 << 26) 事件结构不明

	// MESSAGE_AUDIT (1 << 27)

	OnMessageAuditPass   func(s uint32, bot *Bot, d *MessageAudited)
	OnMessageAuditReject func(s uint32, bot *Bot, d *MessageAudited)
	// FORUMS_EVENT (1 << 28)  // 论坛事件，仅 *私域* 机器人能够设置此 intents。

	OnForumThreadCreate       func(s uint32, bot *Bot, d *Thread)
	OnForumThreadUpdate       func(s uint32, bot *Bot, d *Thread)
	OnForumThreadDelete       func(s uint32, bot *Bot, d *Thread)
	OnForumPostCreate         func(s uint32, bot *Bot, d *Post)
	OnForumPostDelete         func(s uint32, bot *Bot, d *Post)
	OnForumReplyCreate        func(s uint32, bot *Bot, d *Reply)
	OnForumReplyDelete        func(s uint32, bot *Bot, d *Reply)
	OnForumPublishAuditResult func(s uint32, bot *Bot, d *AuditResult)
	// AUDIO_ACTION (1 << 29)

	OnAudioStart  func(s uint32, bot *Bot, d *AudioAction)
	OnAudioFinish func(s uint32, bot *Bot, d *AudioAction)
	OnAudioOnMic  func(s uint32, bot *Bot, d *AudioAction)
	OnAudioOffMic func(s uint32, bot *Bot, d *AudioAction)
	// PUBLIC_GUILD_MESSAGES (1 << 30) // 消息事件，此为公域的消息事件

	OnAtMessageCreate     func(s uint32, bot *Bot, d *Message)
	OnPublicMessageDelete func(s uint32, bot *Bot, d *MessageDelete)
}
