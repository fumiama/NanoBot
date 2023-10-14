package nano

import (
	"encoding/base64"
	"encoding/json"
	"net"
	"net/http"
	"reflect"
	"sync"
	"time"
	"unsafe"

	"github.com/RomiChan/syncx"
	"github.com/RomiChan/websocket"
	log "github.com/sirupsen/logrus"
)

var clients = syncx.Map[string, *Bot]{}

// Bot 一个机器人实例的配置
type Bot struct {
	AppID      string        // AppID is BotAppID（开发者ID）
	Token      string        // Token is 机器人令牌
	Secret     string        // Secret is 机器人密钥
	SuperUsers []string      // SuperUsers 超级用户
	Timeout    time.Duration // Timeout is API 调用超时
	Handler    *Handler      // Handler 注册对各种事件的处理

	gateway    string  // gateway 获得的网关
	shard      [2]byte // shard 分片
	seq        uint32
	handlers   map[string]GeneralHandleType // handlers 方便调用的 handler
	mu         sync.Mutex                   // 写锁
	conn       *websocket.Conn
	heartbeat  uint32 // heartbeat 心跳周期, 单位毫秒
	intents    uint32
	properties json.RawMessage
	ready      EventReady
}

// Init 初始化, 只需执行一次
func (b *Bot) Init(gateway string, shard [2]byte, intents uint32, properties json.RawMessage) *Bot {
	b.gateway = gateway
	b.shard = shard
	b.intents = intents
	b.properties = properties
	if b.Handler != nil {
		h := reflect.ValueOf(b.Handler).Elem()
		t := h.Type()
		b.handlers = make(map[string]GeneralHandleType, h.NumField()*4)
		for i := 0; i < h.NumField(); i++ {
			f := h.Field(i)
			if f.IsZero() {
				continue
			}
			tp := t.Field(i).Name[2:] // skip On
			log.Infoln("[bot] 注册处理函数", tp)
			handler := f.Interface()
			b.handlers[tp] = *(*GeneralHandleType)(unsafe.Add(unsafe.Pointer(&handler), unsafe.Sizeof(uintptr(0))))
		}
	}
	return b
}

// Authorization 返回 Authorization Header value
func (bot *Bot) Authorization() string {
	return "Bot " + bot.AppID + "." + bot.Token
}

// receive 收一个 payload
func (bot *Bot) reveive() (payload WebsocketPayload, err error) {
	err = bot.conn.ReadJSON(&payload)
	return
}

// Connect 连接到 Gateway + 鉴权连接
//
// https://bot.q.qq.com/wiki/develop/api/gateway/reference.html#_1-%E8%BF%9E%E6%8E%A5%E5%88%B0-gateway
func (bot *Bot) Connect() {
	network, address := resolveURI(bot.gateway)
	log.Infoln("[bot] 开始尝试连接到网关:", network, ", AppID:", bot.AppID)
	dialer := websocket.Dialer{
		NetDial: func(_, addr string) (net.Conn, error) {
			if network == "unix" {
				host, _, err := net.SplitHostPort(addr)
				if err != nil {
					host = addr
				}
				filepath, err := base64.RawURLEncoding.DecodeString(host)
				if err == nil {
					addr = BytesToString(filepath)
				}
			}
			return net.Dial(network, addr) // support unix socket transport
		},
	}
	for {
		conn, resp, err := dialer.Dial(address, http.Header{})
		if err != nil {
			log.Warnf("[bot] 连接到网关 %v 时出现错误: %v", bot.gateway, err)
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		bot.conn = conn
		_ = resp.Body.Close()
		payload, err := bot.reveive()
		if err != nil {
			log.Warnln("[bot] 获取心跳间隔时出现错误:", err)
			_ = conn.Close()
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		bot.heartbeat, err = payload.GetHeartbeatInterval()
		if err != nil {
			log.Warnln("[bot] 解析心跳间隔时出现错误:", err)
			_ = conn.Close()
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		payload.Op = OpCodeIdentify
		err = payload.WrapData(&OpCodeIdentifyMessage{
			Token:      bot.Authorization(),
			Intents:    bot.intents,
			Shard:      bot.shard,
			Properties: bot.properties,
		})
		if err != nil {
			log.Warnln("[bot] 包装 Identify 时出现错误:", err)
			_ = conn.Close()
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		err = bot.SendPayload(&payload)
		if err != nil {
			log.Warnln("[bot] 发送 Identify 时出现错误:", err)
			_ = conn.Close()
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		payload, err = bot.reveive()
		if err != nil {
			log.Warnln("[bot] 获取 EventReady 时出现错误:", err)
			_ = conn.Close()
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		bot.ready, err = payload.GetEventReady()
		if err != nil {
			log.Warnln("[bot] 解析 EventReady 时出现错误:", err)
			_ = conn.Close()
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		break
	}
	clients.Store(bot.ready.User.ID, bot)
	log.Infoln("[bot] 连接到网关成功, 用户名:", bot.ready.User.Username)
	go bot.doheartbeat()
}

// doheartbeat 按指定间隔进行心跳包发送
func (bot *Bot) doheartbeat() {
	payload := struct {
		Op OpCode  `json:"op"`
		D  *uint32 `json:"d"`
	}{Op: OpCodeHeartbeat}
	t := time.NewTicker(time.Duration(bot.heartbeat) * time.Millisecond)
	defer t.Stop()
	time.Sleep(time.Minute)
	for range t.C {
		if bot.seq == 0 {
			payload.D = nil
		} else {
			payload.D = &bot.seq
		}
		bot.mu.Lock()
		err := bot.conn.WriteJSON(&payload)
		bot.mu.Unlock()
		if err != nil {
			log.Warnln("[bot] 发送心跳时出现错误:", err)
		}
	}
}
