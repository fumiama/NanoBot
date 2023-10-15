package nano

import (
	"encoding/base64"
	"encoding/json"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/RomiChan/syncx"
	"github.com/RomiChan/websocket"
	log "github.com/sirupsen/logrus"
)

var clients = syncx.Map[string, *Bot]{}

// Bot 一个机器人实例的配置
type Bot struct {
	AppID      string          // AppID is BotAppID（开发者ID）
	Token      string          // Token is 机器人令牌
	Secret     string          // Secret is 机器人密钥
	SuperUsers []string        // SuperUsers 超级用户
	Timeout    time.Duration   // Timeout is API 调用超时
	Handler    *Handler        // Handler 注册对各种事件的处理
	Intents    uint32          // Intents 欲接收的事件
	Properties json.RawMessage // Properties 一些环境变量, 目前没用

	gateway   string                       // gateway 获得的网关
	shard     [2]byte                      // shard 分片
	seq       uint32                       // seq 最新的 s
	handlers  map[string]GeneralHandleType // handlers 方便调用的 handler
	mu        sync.Mutex                   // 写锁
	conn      *websocket.Conn              // conn 目前的 wss 连接
	heartbeat uint32                       // heartbeat 心跳周期, 单位毫秒
	hbonce    sync.Once                    // hbonce 保证仅执行一次 heartbeat

	ready EventReady //
}

// Init 初始化, 只需执行一次
func (b *Bot) Init(gateway string, shard [2]byte) *Bot {
	b.gateway = gateway
	b.shard = shard
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
			log.Infoln(getLogHeader(), "注册处理函数", tp)
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
func (bot *Bot) Connect() *Bot {
	network, address := resolveURI(bot.gateway)
	log.Infoln(getLogHeader(), "开始尝试连接到网关:", address, ", AppID:", bot.AppID)
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
			log.Warnf(getLogHeader(), "连接到网关 %v 时出现错误: %v", bot.gateway, err)
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		bot.conn = conn
		_ = resp.Body.Close()
		payload, err := bot.reveive()
		if err != nil {
			log.Warnln(getLogHeader(), "获取心跳间隔时出现错误:", err)
			_ = conn.Close()
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		hb, err := payload.GetHeartbeatInterval()
		if err != nil {
			log.Warnln(getLogHeader(), "解析心跳间隔时出现错误:", err)
			_ = conn.Close()
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		payload.Op = OpCodeIdentify
		err = payload.WrapData(&OpCodeIdentifyMessage{
			Token:      bot.Authorization(),
			Intents:    bot.Intents,
			Shard:      bot.shard,
			Properties: bot.Properties,
		})
		if err != nil {
			log.Warnln(getLogHeader(), "包装 Identify 时出现错误:", err)
			_ = conn.Close()
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		err = bot.SendPayload(&payload)
		if err != nil {
			log.Warnln(getLogHeader(), "发送 Identify 时出现错误:", err)
			_ = conn.Close()
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		payload, err = bot.reveive()
		if err != nil {
			log.Warnln(getLogHeader(), "获取 EventReady 时出现错误:", err)
			_ = conn.Close()
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		bot.ready, err = payload.GetEventReady()
		if err != nil {
			log.Warnln(getLogHeader(), "解析 EventReady 时出现错误:", err)
			_ = conn.Close()
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		atomic.StoreUint32(&bot.heartbeat, hb)
		break
	}
	clients.Store(bot.Token+"_"+strconv.Itoa(int(bot.shard[0])), bot)
	log.Infoln(getLogHeader(), "连接到网关成功, 用户名:", bot.ready.User.Username)
	bot.hbonce.Do(func() {
		go bot.doheartbeat()
	})
	return bot
}

// doheartbeat 按指定间隔进行心跳包发送
func (bot *Bot) doheartbeat() {
	payload := struct {
		Op OpCode  `json:"op"`
		D  *uint32 `json:"d"`
	}{Op: OpCodeHeartbeat}
	for {
		if atomic.LoadUint32(&bot.heartbeat) == 0 {
			time.Sleep(time.Second)
			log.Warnln(getLogHeader(), "等待服务器建立连接...")
			continue
		}
		time.Sleep(time.Duration(bot.heartbeat) * time.Millisecond)
		if bot.seq == 0 {
			payload.D = nil
		} else {
			payload.D = &bot.seq
		}
		bot.mu.Lock()
		err := bot.conn.WriteJSON(&payload)
		bot.mu.Unlock()
		if err != nil {
			log.Warnln(getLogHeader(), "发送心跳时出现错误:", err)
		}
	}
}

// Resume 恢复连接
//
// https://bot.q.qq.com/wiki/develop/api/gateway/reference.html#_4-%E6%81%A2%E5%A4%8D%E8%BF%9E%E6%8E%A5
func (bot *Bot) Resume() error {
	network, address := resolveURI(bot.gateway)
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
	conn, resp, err := dialer.Dial(address, http.Header{})
	if err != nil {
		return err
	}
	bot.conn = conn
	_ = resp.Body.Close()
	payload := WebsocketPayload{Op: OpCodeResume}
	payload.WrapData(&struct {
		T string `json:"token"`
		S string `json:"session_id"`
		Q uint32 `json:"seq"`
	}{bot.Authorization(), bot.ready.SessionID, bot.seq})
	return bot.SendPayload(&payload)
}

// Listen 监听事件
func (bot *Bot) Listen() {
	log.Infoln(getLogHeader(), "开始监听", bot.ready.User.Username, "的事件")
	payload := WebsocketPayload{}
	lastheartbeat := time.Now()
	for {
		payload.Reset()
		err := bot.conn.ReadJSON(&payload)
		if err != nil { // reconnect
			atomic.StoreUint32(&bot.heartbeat, 0)
			k := bot.Token + "_" + strconv.Itoa(int(bot.shard[0]))
			clients.Delete(k)
			log.Warnln(getLogHeader(), bot.ready.User.Username, "的网关连接断开, 尝试恢复:", err)
			for {
				time.Sleep(time.Second)
				err = bot.Resume()
				if err == nil {
					break
				}
				log.Warnln(getLogHeader(), bot.ready.User.Username, "的网关连接恢复失败:", err)
			}
			clients.Store(k, bot)
			continue
		}
		log.Debugln(getLogHeader(), "接收到第", payload.S, "个事件:", payload.Op, ", 类型:", payload.T, ", 数据:", BytesToString(payload.D))
		bot.seq = payload.S
		switch payload.Op {
		case OpCodeDispatch: // Receive
			switch payload.T {
			case "RESUMED":
				log.Infoln(getLogHeader(), bot.ready.User.Username, "的网关连接恢复完成")
			}
		case OpCodeHeartbeat: // Send/Receive
			log.Debugln(getLogHeader(), "收到服务端推送心跳, 间隔:", time.Since(lastheartbeat))
			lastheartbeat = time.Now()
		case OpCodeReconnect: // Receive
			log.Warnln(getLogHeader(), "收到服务端通知重连")
			atomic.StoreUint32(&bot.heartbeat, 0)
			bot.Connect()
		case OpCodeInvalidSession: // Receive
			log.Warnln(getLogHeader(), bot.ready.User.Username, "的网关连接恢复失败: InvalidSession, 尝试重连...")
			atomic.StoreUint32(&bot.heartbeat, 0)
			bot.Connect()
		case OpCodeHello: // Receive
			intv, err := payload.GetHeartbeatInterval()
			if err != nil {
				log.Warnln(getLogHeader(), "解析心跳间隔时出现错误:", err)
				continue
			}
			atomic.StoreUint32(&bot.heartbeat, intv)
		case OpCodeHeartbeatACK: // Receive/Reply
			log.Debugln(getLogHeader(), "收到心跳返回, 间隔:", time.Since(lastheartbeat))
			lastheartbeat = time.Now()
		case OpCodeHTTPCallbackACK: // Reply
		}
	}
}
