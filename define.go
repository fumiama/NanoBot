package nano

import "encoding/json"

const (
	// StandardAPI 正式环境接口域名
	StandardAPI = `https://api.sgroup.qq.com`
	// SandboxAPI 沙箱环境接口域名
	SandboxAPI = `https://sandbox.api.sgroup.qq.com`
)

var (
	OpenAPI = StandardAPI // OpenAPI 实际使用的 API, 默认 StandardAPI, 可自行赋值配置
)

// CodeMessageBase 各种消息都有的 code + message 基类
type CodeMessageBase struct {
	C int    `json:"code"`
	M string `json:"message"`
}

// OpCode https://bot.q.qq.com/wiki/develop/api/gateway/opcode.html
type OpCode int

const (
	OpCodeDispatch  OpCode = iota // Receive
	OpCodeHeartbeat               // Send/Receive
	OpCodeIdentify                // Send
	OpCodeEmpty1
	OpCodeEmpty2
	OpCodeEmpty3
	OpCodeResume    // Send
	OpCodeReconnect // Receive
	OpCodeEmpty4
	OpCodeInvalidSession  // Receive
	OpCodeHello           // Receive
	OpCodeHeartbeatACK    // Receive/Reply
	OpCodeHTTPCallbackACK // Reply
)

// WebsocketPayload payload 指的是在 websocket 连接上传输的数据，网关的上下行消息采用的都是同一个结构
//
// https://bot.q.qq.com/wiki/develop/api/gateway/reference.html
type WebsocketPayload struct {
	Op OpCode          `json:"op"`
	D  json.RawMessage `json:"d"`
	S  int             `json:"s"`
	T  string          `json:"t"`
}

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
