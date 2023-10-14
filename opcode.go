package nano

import "encoding/json"

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

// OpCodeIdentifyMessage https://bot.q.qq.com/wiki/develop/api/gateway/reference.html#_2-%E9%89%B4%E6%9D%83%E8%BF%9E%E6%8E%A5
type OpCodeIdentifyMessage struct {
	Token      string          `json:"token"`
	Intents    uint32          `json:"intents"`
	Shard      [2]byte         `json:"shard"`
	Properties json.RawMessage `json:"properties"`
}
