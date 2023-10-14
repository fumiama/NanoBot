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
}
