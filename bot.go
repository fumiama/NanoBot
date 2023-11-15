package nano

import (
	"encoding/base64"
	"encoding/json"
	"errors"
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

var (
	clients   = syncx.Map[string, *Bot]{}
	isrunning uintptr
)

// Bot 一个机器人实例的配置
type Bot struct {
	AppID      string          `yaml:"AppID"` // AppID is BotAppID（开发者ID）
	Token      string          `yaml:"Token"` // Token is 机器人令牌 有 Secret 则使用新版 API
	token      string          // token 是通过 secret 获得的残血 token
	Secret     string          `yaml:"Secret"`     // Secret is 机器人令牌 V2 (AppSecret/ClientSecret) 沙盒目前虽然能登录但无法收发消息
	SuperUsers []string        `yaml:"SuperUsers"` // SuperUsers 超级用户
	Timeout    time.Duration   `yaml:"Timeout"`    // Timeout is API 调用超时
	Handler    *Handler        `yaml:"-"`          // Handler 注册对各种事件的处理
	Intents    uint32          `yaml:"Intents"`    // Intents 欲接收的事件
	ShardIndex uint8           `yaml:"ShardIndex"` // ShardIndex 本连接为第几个分片, 默认 1, 0 为不使用分片
	ShardCount uint8           `yaml:"ShardCount"` // ShardCount 分片总数
	shard      [2]byte         // shard 分片
	Properties json.RawMessage `yaml:"Properties"` // Properties 一些环境变量, 目前没用

	gateway   string                      // gateway 获得的网关
	seq       uint32                      // seq 最新的 s
	heartbeat uint32                      // heartbeat 心跳周期, 单位毫秒
	expiresec int64                       // expiresec Token 有效时间
	handlers  map[string]eventHandlerType // handlers 方便调用的 handler
	mu        sync.Mutex                  // 写锁
	conn      *websocket.Conn             // conn 目前的 wss 连接
	hbonce    sync.Once                   // hbonce 保证仅执行一次 heartbeat
	exonce    sync.Once                   // exonce 保证仅执行一次刷新 token
	client    *http.Client                // client 主要配置 timeout

	ready EventReady // ready 连接成功后下发的 bot 基本信息
}

// GetReady 获得 bot 基本信息
func (ctx *Ctx) GetReady() *EventReady {
	return &ctx.caller.ready
}

// getinitinfo 获得 gateway 和 shard
func (bot *Bot) getinitinfo() (secret, gw string, shard [2]byte, err error) {
	shard[1] = 1
	if bot.client == nil {
		bot.client = http.DefaultClient
	}
	secret = bot.Secret
	if bot.Secret != "" {
		bot.Secret = ""
	}
	if bot.ShardIndex == 0 {
		gw, err = bot.GetGeneralWSSGatewayNoContext()
		if err != nil {
			return
		}
	} else {
		var sgw *ShardWSSGateway
		sgw, err = bot.GetShardWSSGatewayNoContext()
		if err != nil {
			return
		}
		if bot.ShardCount == 0 {
			log.Infoln(getLogHeader(), "使用网关推荐Shards数:", sgw.Shards)
			bot.ShardCount = uint8(sgw.Shards)
		}
		if bot.ShardCount <= bot.ShardIndex {
			err = errors.New("shard index " + strconv.Itoa(int(bot.ShardIndex)) + " >= suggested size " + strconv.Itoa(sgw.Shards))
			return
		}
		gw = sgw.URL
		shard[0] = byte(bot.ShardIndex)
		shard[1] = byte(bot.ShardCount)
	}
	return
}

// Start clients without blocking
func Start(bots ...*Bot) error {
	if !atomic.CompareAndSwapUintptr(&isrunning, 0, 1) {
		log.Warnln(getLogHeader(), "已忽略重复调用的", getThisFuncName())
	}
	for _, b := range bots {
		s, gw, shard, err := b.getinitinfo()
		if err != nil {
			return err
		}
		go b.Init(s, gw, shard).Connect().Listen()
	}
	return nil
}

// Run clients and block self in listening last one
func Run(preblock func(), bots ...*Bot) error {
	if !atomic.CompareAndSwapUintptr(&isrunning, 0, 1) {
		log.Warnln(getLogHeader(), "已忽略重复调用的", getThisFuncName())
	}
	var b *Bot
	switch len(bots) {
	case 0:
		return nil
	case 1:
		b = bots[0]
		s, gw, shard, err := b.getinitinfo()
		if err != nil {
			return err
		}
		b.Init(s, gw, shard)
	default:
		for _, b := range bots[:len(bots)-1] {
			s, gw, shard, err := b.getinitinfo()
			if err != nil {
				return err
			}
			go b.Init(s, gw, shard).Connect().Listen()
		}
		b = bots[len(bots)-1]
		s, gw, shard, err := b.getinitinfo()
		if err != nil {
			return err
		}
		b.Init(s, gw, shard)
	}
	b.Connect()
	if preblock != nil {
		preblock()
	}
	b.Listen()
	return nil
}

// Init 初始化, 只需执行一次
func (bot *Bot) Init(secret, gateway string, shard [2]byte) *Bot {
	bot.gateway = gateway
	bot.shard = shard
	if bot.Timeout == 0 {
		bot.Timeout = time.Minute
	}
	bot.client = &http.Client{
		Timeout: bot.Timeout,
	}
	if bot.Handler != nil {
		h := reflect.ValueOf(bot.Handler).Elem()
		t := h.Type()
		bot.handlers = make(map[string]eventHandlerType, h.NumField()*4)
		for i := 0; i < h.NumField(); i++ {
			f := h.Field(i)
			if f.IsZero() {
				continue
			}
			tp := t.Field(i).Name[2:] // skip On
			log.Infoln(getLogHeader(), "注册处理函数", tp)
			handler := f.Interface()
			bot.handlers[tp] = eventHandlerType{
				h: *(*generalHandleType)(unsafe.Add(unsafe.Pointer(&handler), unsafe.Sizeof(uintptr(0)))),
				t: t.Field(i).Type.In(2).Elem(),
			}
		}
	}
	bot.Secret = secret
	if bot.IsV2() {
		for {
			err := bot.GetAppAccessTokenNoContext()
			if err == nil {
				log.Infoln(getLogHeader(), "获得 Token: "+bot.token+", 超时:", bot.expiresec, "秒")
				bot.exonce.Do(func() {
					go bot.refreshtoken()
				})
				break
			}
			log.Infoln(getLogHeader(), "获得 Token 失败:", err)
			time.Sleep(time.Second * 3)
		}
	}
	return bot
}

// IsV2 判断是否运行于 V2 API 下
func (bot *Bot) IsV2() bool {
	return bot.Secret != ""
}

// Authorization 返回 Authorization Header value
func (bot *Bot) Authorization() string {
	if bot.IsV2() {
		return "QQBot " + bot.token
	}
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
		bot.ready, bot.seq, err = payload.GetEventReady()
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

// refreshtoken 以 Expire 为间隔刷新 Token
func (bot *Bot) refreshtoken() {
	for {
		time.Sleep(time.Second * 10)
		if atomic.LoadUint32(&bot.heartbeat) == 0 {
			log.Warnln(getLogHeader(), "等待服务器建立连接...")
			continue
		}
		time.Sleep(time.Duration(bot.expiresec) * time.Second)
		err := bot.GetAppAccessTokenNoContext()
		if err != nil {
			log.Warnln(getLogHeader(), "刷新 Token 时出现错误:", err)
		} else {
			log.Infoln(getLogHeader(), "刷新 Token: "+bot.token+", 超时:", bot.expiresec, "秒")
		}
	}
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
		log.Debug(getLogHeader(), " 接收到第 ", payload.S, " 个事件: ", payload.Op, ", 类型: ", payload.T, ", 数据: ", BytesToString(payload.D))
		switch payload.Op {
		case OpCodeDispatch: // Receive
			if payload.S <= bot.seq {
				log.Warn(getLogHeader(), " 忽略重复编号: ", payload.S, ", 事件: ", payload.Op, ", 类型: ", payload.T)
				continue
			}
			switch payload.T {
			case "RESUMED":
				log.Infoln(getLogHeader(), bot.ready.User.Username, "的网关连接恢复完成")
			default:
				bot.processEvent(&payload)
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
		default:
			log.Warn(getLogHeader(), " 忽略未知事件, 序号: ", payload.S, ", Op: ", payload.Op, ", 类型: ", payload.T, ", 数据: ", BytesToString(payload.D))
		}
		if payload.S > bot.seq {
			bot.seq = payload.S
		}
	}
}

// GetBot 获取指定的bot (Ctx)实例
func GetBot(id string) *Ctx {
	caller, ok := clients.Load(id)
	if !ok {
		return nil
	}
	return &Ctx{caller: caller}
}

// RangeBot 遍历所有bot (Ctx)实例
//
// 单次操作返回 true 则继续遍历，否则退出
func RangeBot(iter func(id string, ctx *Ctx) bool) {
	clients.Range(func(key string, value *Bot) bool {
		return iter(key, &Ctx{caller: value})
	})
}

// GetFirstSuperUser 在 ids 中获得 SuperUsers 列表的首个 qq
//
// 找不到返回 nil
func (bot *Bot) GetFirstSuperUser(ids ...string) string {
	m := make(map[string]struct{}, len(ids)*4)
	for _, qq := range ids {
		m[qq] = struct{}{}
	}
	for _, qq := range bot.SuperUsers {
		if _, ok := m[qq]; ok {
			return qq
		}
	}
	return ""
}
