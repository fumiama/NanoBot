package nano

import "time"

// Bot 一个机器人实例的配置
type Bot struct {
	AppID      string        // AppID is BotAppID（开发者ID）
	Token      string        // Token is 机器人令牌
	Secret     string        // Secret is 机器人密钥
	SuperUsers []string      // SuperUsers 超级用户
	Timeout    time.Duration // Timeout is API 调用超时
	Handler    *Handler      // Handler 注册对各种事件的处理

	handlers map[string]GeneralHandleType // handlers 方便调用的 handler
}

// Init 初始化, 只需执行一次
func (b *Bot) Init() {

}

// Authorization 返回 Authorization Header value
func (bot *Bot) Authorization() string {
	return "Bot " + bot.AppID + "." + bot.Token
}
