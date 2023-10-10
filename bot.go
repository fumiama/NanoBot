package nano

import "time"

// Bot 一个机器人实例的配置
type Bot struct {
	AppID   string        // AppID is BotAppID（开发者ID）
	Token   string        // Token is 机器人令牌
	Secret  string        // Secret is 机器人密钥
	Timeout time.Duration // Timeout is API 调用超时
}

// Authorization 返回 Authorization Header value
func (bot *Bot) Authorization() string {
	return "Bot " + bot.AppID + "." + bot.Token
}
