package nano

import "unsafe"

// GeneralHandleType 作为通用的 handler 函数调用约定使用
type GeneralHandleType func(int, *Bot, unsafe.Pointer)

// Handler 事件订阅
//
// https://bot.q.qq.com/wiki/develop/api/gateway/intents.html
type Handler struct {
	// GUILDS (1 << 0)

	OnGuildCreate   func(s int, bot *Bot, d *Guild)
	OnGuildUpdate   func(s int, bot *Bot, d *Guild)
	OnGuildDelete   func(s int, bot *Bot, d *Guild)
	OnChannelCreate func(s int, bot *Bot, d *Channel)
	OnChannelUpdate func(s int, bot *Bot, d *Channel)
	OnChannelDelete func(s int, bot *Bot, d *Channel)
	// GUILD_MEMBERS (1 << 1)

	OnGuildMemberAdd    func(s int, bot *Bot, d *Member)
	OnGuildMemberUpdate func(s int, bot *Bot, d *Member)
	OnGuildMemberRemove func(s int, bot *Bot, d *Member)
	// GUILD_MESSAGES (1 << 9)    // 消息事件，仅 *私域* 机器人能够设置此 intents。

	OnMessageCreate func(s int, bot *Bot, d *Message)
	OnMessageDelete func(s int, bot *Bot, d *Message)
	// GUILD_MESSAGE_REACTIONS (1 << 10)

	OnMessageReactionAdd    func(s int, bot *Bot, d *MessageReaction)
	OnMessageReactionRemove func(s int, bot *Bot, d *MessageReaction)
	// DIRECT_MESSAGE (1 << 12)

	OnDirectMessageCreate func(s int, bot *Bot, d *Message)
	OnDirectMessageDelete func(s int, bot *Bot, d *Message)
	// OPEN_FORUMS_EVENT (1 << 18)      // 论坛事件, 此为公域的论坛事件

	OnOpenForumThreadCreate func(s int, bot *Bot, d *Thread)
	OnOpenForumThreadUpdate func(s int, bot *Bot, d *Thread)
	OnOpenForumThreadDelete func(s int, bot *Bot, d *Thread)
	OnOpenForumPostCreate   func(s int, bot *Bot, d *Post)
	OnOpenForumPostDelete   func(s int, bot *Bot, d *Post)
	OnOpenForumReplyCreate  func(s int, bot *Bot, d *Reply)
	OnOpenForumReplyDelete  func(s int, bot *Bot, d *Reply)
	// AUDIO_OR_LIVE_CHANNEL_MEMBER (1 << 19)  // 音视频/直播子频道成员进出事件

	OnAudioOrLiveChannelMemberEnter func(s int, bot *Bot, d *AudioLiveChannelUsersChange)
	OnAudioOrLiveChannelMemberExit  func(s int, bot *Bot, d *AudioLiveChannelUsersChange)
	// INTERACTION (1 << 26) 事件结构不明

	// MESSAGE_AUDIT (1 << 27)

	OnMessageAuditPass   func(s int, bot *Bot, d *MessageAudited)
	OnMessageAuditReject func(s int, bot *Bot, d *MessageAudited)
	// FORUMS_EVENT (1 << 28)  // 论坛事件，仅 *私域* 机器人能够设置此 intents。

	OnForumThreadCreate       func(s int, bot *Bot, d *Thread)
	OnForumThreadUpdate       func(s int, bot *Bot, d *Thread)
	OnForumThreadDelete       func(s int, bot *Bot, d *Thread)
	OnForumPostCreate         func(s int, bot *Bot, d *Post)
	OnForumPostDelete         func(s int, bot *Bot, d *Post)
	OnForumReplyCreate        func(s int, bot *Bot, d *Reply)
	OnForumReplyDelete        func(s int, bot *Bot, d *Reply)
	OnForumPublishAuditResult func(s int, bot *Bot, d *AuditResult)
	// AUDIO_ACTION (1 << 29)

	OnAudioStart  func(s int, bot *Bot, d *AudioAction)
	OnAudioFinish func(s int, bot *Bot, d *AudioAction)
	OnAudioOnMic  func(s int, bot *Bot, d *AudioAction)
	OnAudioOffMic func(s int, bot *Bot, d *AudioAction)
	// PUBLIC_GUILD_MESSAGES (1 << 30) // 消息事件，此为公域的消息事件

	OnAtMessageCreate     func(s int, bot *Bot, d *Message)
	OnPublicMessageDelete func(s int, bot *Bot, d *Message)
}
