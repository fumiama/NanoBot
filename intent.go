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

	IntentAll = IntentGuilds | IntentGuildMembers | IntentGuildMessages | IntentGuildMessageReactions |
		IntentDirectMessage | IntentOpenForumsEvent | IntentAudioOrLiveChannelMember | IntentInteraction |
		IntentMessageAudit | IntentForumsEvent | IntentAudioAction | IntentPublicGuildMessages
)
