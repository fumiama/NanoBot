package main

import (
	_ "github.com/fumiama/NanoBot/example/echo"

	nano "github.com/fumiama/NanoBot"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	nano.OpenAPI = nano.SandboxAPI
	nano.OnMessageFullMatch("help").SetBlock(true).
		Handle(func(ctx *nano.Ctx) {
			_, _ = ctx.SendPlainMessage(false, "echo string")
		})
	nano.Run(nil, &nano.Bot{
		AppID:      "你的AppID",
		Token:      "你的Token",
		Secret:     "你的Secret, 目前没用到, 可以不填",
		Intents:    nano.IntentGuildPublic,
		SuperUsers: []string{"用户ID1", "用户ID2"},
	})
}
