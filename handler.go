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

}
