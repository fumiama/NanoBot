package nano

// https://bot.q.qq.com/wiki/develop/api/gateway/intents.html
const (
	IntentGuilds                   = 1 << 0
	IntentGuildMembers             = 1 << 1
	IntentGuildMessages            = 1 << 9
	IntentGuildMessageReactions    = 1 << 10
	IntentDirectMessage            = 1 << 12
	IntentOpenForumsEvent          = 1 << 18
	IntentAudioOrLiveChannelMember = 1 << 19
	IntentInteraction              = 1 << 26
	IntentMessageAudit             = 1 << 27
	IntentForumsEvent              = 1 << 28
	IntentAudioAction              = 1 << 29
	IntentPublicGuildMessages      = 1 << 30

	// IntentAll 监听全部事件
	IntentAll = IntentGuilds | IntentGuildMembers | IntentGuildMessages | IntentGuildMessageReactions |
		IntentDirectMessage | IntentOpenForumsEvent | IntentAudioOrLiveChannelMember | IntentInteraction |
		IntentMessageAudit | IntentForumsEvent | IntentAudioAction | IntentPublicGuildMessages
	// IntentPublic 监听公域事件
	IntentPublic = IntentGuilds | IntentGuildMembers | IntentGuildMessageReactions |
		IntentDirectMessage | IntentOpenForumsEvent | IntentAudioOrLiveChannelMember | IntentInteraction |
		IntentMessageAudit | IntentAudioAction | IntentPublicGuildMessages
	// IntentPrivate 监听私域事件
	IntentPrivate = IntentGuilds | IntentGuildMembers | IntentGuildMessages | IntentGuildMessageReactions |
		IntentDirectMessage | IntentAudioOrLiveChannelMember | IntentInteraction |
		IntentMessageAudit | IntentForumsEvent | IntentAudioAction
)
