package nano

import (
	"encoding/json"
	"reflect"

	log "github.com/sirupsen/logrus"
)

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
}
