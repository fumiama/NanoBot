package nano

import (
	"encoding/json"
	"errors"
	"strconv"
)

// WebsocketPayload payload 指的是在 websocket 连接上传输的数据，网关的上下行消息采用的都是同一个结构
//
// https://bot.q.qq.com/wiki/develop/api/gateway/reference.html
type WebsocketPayload struct {
	Op OpCode          `json:"op"`
	D  json.RawMessage `json:"d,omitempty"`
	S  int             `json:"s,omitempty"`
	T  string          `json:"t,omitempty"`
}

// GetHeartbeatInterval OpCodeHello 获得心跳周期 单位毫秒
func (wp *WebsocketPayload) GetHeartbeatInterval() (uint32, error) {
	if wp.Op != OpCodeHello {
		return 0, errors.New(getThisFuncName() + " unexpected OpCode " + strconv.Itoa(int(wp.Op)) + ", T: " + wp.T + ", D: " + BytesToString(wp.D))
	}
	data := &struct {
		H uint32 `json:"heartbeat_interval"`
	}{}
	err := json.Unmarshal(wp.D, data)
	return data.H, err
}

// SendPayload 发送 ws 包
func (bot *Bot) SendPayload(wp *WebsocketPayload) error {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	return bot.conn.WriteJSON(wp)
}

// WrapData 将结构体序列化到 wp.D
func (wp *WebsocketPayload) WrapData(v any) (err error) {
	wp.D, err = json.Marshal(v)
	return
}

// EventReady https://bot.q.qq.com/wiki/develop/api/gateway/reference.html#_2-%E9%89%B4%E6%9D%83%E8%BF%9E%E6%8E%A5
type EventReady struct {
	Version   int     `json:"version"`
	SessionID string  `json:"session_id"`
	User      *User   `json:"user"`
	Shard     [2]byte `json:"shard"`
}

// GetEventReady OpCodeDispatch READY
func (wp *WebsocketPayload) GetEventReady() (er EventReady, err error) {
	if wp.Op != OpCodeDispatch {
		err = errors.New(getThisFuncName() + " unexpected OpCode " + strconv.Itoa(int(wp.Op)) + ", T: " + wp.T + ", D: " + BytesToString(wp.D))
		return
	}
	if wp.T != "READY" {
		err = errors.New(getThisFuncName() + " unexpected event type " + wp.T + ", T: " + wp.T + ", D: " + BytesToString(wp.D))
		return
	}
	err = json.Unmarshal(wp.D, &er)
	return
}
