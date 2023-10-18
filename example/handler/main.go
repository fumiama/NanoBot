package main

import (
	"strings"

	nano "github.com/fumiama/NanoBot"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	nano.OpenAPI = nano.SandboxAPI
	nano.Run(nil, &nano.Bot{
		AppID:   "你的AppID",
		Token:   "你的Token",
		Secret:  "你的Secret, 目前没用到, 可以不填",
		Intents: nano.IntentPublic,
		Handler: &nano.Handler{
			OnAtMessageCreate: func(s uint32, bot *nano.Bot, d *nano.Message) {
				u := ""
				if len(d.Attachments) > 0 {
					u = d.Attachments[0].URL
					if !strings.HasPrefix(u, "http") {
						u = "http://" + u
					}
				}
				_, err := bot.PostMessageToChannel(d.ChannelID, &nano.MessagePost{
					Content:        "您发送了: " + d.Content,
					Image:          u,
					ReplyMessageID: d.ID,
					MessageReference: &nano.MessageReference{
						MessageID: d.ID,
					},
				})
				if err != nil {
					bot.PostMessageToChannel(d.ChannelID, &nano.MessagePost{
						Content:        "[ERROR]: " + err.Error(),
						ReplyMessageID: d.ID,
					})
				}
			},
		},
	})
}
