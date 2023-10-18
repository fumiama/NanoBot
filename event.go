package nano

import (
	"encoding/json"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Event ...
type Event struct {
	// Type is payload.T
	Type string
	// Seq 序列号
	Seq uint32
	// Value 是 D
	Value any
	// value is the reflect value of Value
	value reflect.Value
}

// processEvent 处理需要关注的业务事件
func (bot *Bot) processEvent(payload *WebsocketPayload) {
	tp := UnderlineToCamel(payload.T)
	if bot.Handler != nil {
		ev, ok := bot.handlers[tp]
		if !ok {
			return
		}
		log.Debugln(getLogHeader(), "使用 handlers 处理", tp, "事件")
		x := reflect.New(ev.t)
		err := json.Unmarshal(payload.D, x.Interface())
		if err != nil {
			log.Warnln(getLogHeader(), "解析", tp, "事件时出现错误:", err)
			return
		}
		go ev.h(payload.S, bot, x.UnsafePointer())
		return
	}
	ctx := &Ctx{
		Event: Event{
			Type: tp,
			Seq:  payload.S,
		},
		State:  State{},
		caller: bot,
	}
	switch tp {
	case "DirectMessageCreate":
		ctx.IsToMe = true
		fallthrough
	case "MessageCreate", "AtMessageCreate":
		tp = "Message"
	case "DirectMessageDelete":
		ctx.IsToMe = true
		fallthrough
	case "MessageDelete", "PublicMessageDelete":
		tp = "MessageDelete"
	}
	matcherLock.RLock()
	n := len(matcherMap[tp])
	if n == 0 {
		matcherLock.RUnlock()
		return
	}
	log.Debugln(getLogHeader(), "pass", tp, "event to plugins")
	matchers := make([]*Matcher, n)
	copy(matchers, matcherMap[tp])
	matcherLock.RUnlock()
	x := reflect.New(types[ctx.Type])
	err := json.Unmarshal(payload.D, x.Interface())
	if err != nil {
		log.Warnln(getLogHeader(), "解析", ctx.Type, "事件时出现错误:", err)
		return
	}
	ctx.Value = x.Interface()
	ctx.value = x
	switch tp {
	case "Message":
		ctx.Message = (*Message)(x.UnsafePointer())
		if ctx.Message.MentionEveryone {
			ctx.IsToMe = true
		}
		log.Infoln(getLogHeader(), "=>", ctx.Message)
	case "MessageDelete":
		mdl := (*MessageDelete)(x.UnsafePointer())
		opmember, err := ctx.GetGuildMemberOf(mdl.Message.GuildID, mdl.OpUser.ID)
		if err != nil {
			log.Warnln(getLogHeader(), "获取撤回消息者详情错误:", err)
			return
		}
		ctx.Message = (*MessageDelete)(x.UnsafePointer()).Message
		ctx.Message.Member = &Member{
			GuildID: mdl.Message.GuildID,
			User:    ctx.Message.Author,
		}
		ctx.Message.Author = opmember.User
	}
	go match(ctx, matchers)
}

func match(ctx *Ctx, matchers []*Matcher) {
	if ctx.Message != nil && ctx.Message.Content != "" { // 确保无空
		if !ctx.IsToMe {
			ctx.IsToMe = func(ctx *Ctx) bool {
				name := ctx.GetReady().User.Username
				if strings.HasPrefix(ctx.Message.Content, name) {
					log.Debugln(getLogHeader(), "message before process:", ctx.Message.Content)
					ctx.Message.Content = strings.TrimLeft(ctx.Message.Content[len(name):], " ")
					log.Debugln(getLogHeader(), "message after process:", ctx.Message.Content)
					return true
				}
				atme := ctx.AtMe()
				if strings.HasPrefix(ctx.Message.Content, atme) {
					log.Debugln(getLogHeader(), "message before process:", ctx.Message.Content)
					ctx.Message.Content = strings.TrimLeft(ctx.Message.Content[len(atme):], " ")
					log.Debugln(getLogHeader(), "message after process:", ctx.Message.Content)
					return true
				}
				return OnlyPrivate(ctx)
			}(ctx)
		}
	}
	log.Debugln(getLogHeader(), "message is to me:", ctx.IsToMe)
loop:
	for _, matcher := range matchers {
		for k := range ctx.State { // Clear State
			delete(ctx.State, k)
		}
		matcherLock.RLock()
		m := matcher.copy()
		matcherLock.RUnlock()
		ctx.ma = m

		// pre handler
		if m.Engine != nil {
			for _, handler := range m.Engine.preHandler {
				if !handler(ctx) { // 有 pre handler 未满足
					if m.Break { // 阻断后续
						break loop
					}
					continue loop
				}
			}
		}

		for _, rule := range m.Rules {
			if rule != nil && !rule(ctx) { // 有 Rule 的条件未满足
				if m.Break { // 阻断后续
					break loop
				}
				continue loop
			}
		}

		// mid handler
		if m.Engine != nil {
			for _, handler := range m.Engine.midHandler {
				if !handler(ctx) { // 有 mid handler 未满足
					if m.Break { // 阻断后续
						break loop
					}
					continue loop
				}
			}
		}

		if m.Process != nil {
			m.Process(ctx) // 处理事件
		}
		if matcher.Temp { // 临时 Matcher 删除
			matcher.Delete()
		}

		if m.Engine != nil {
			// post handler
			for _, handler := range m.Engine.postHandler {
				handler(ctx)
			}
		}

		if m.Block { // 阻断后续
			break loop
		}
	}
}
