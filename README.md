<div align="center">
  <a href="https://crypko.ai/crypko/GtWYDpVMx5GYm/">
  <img src=".github/nano.jpeg" alt="东云名乃" width = "256">
  </a><br>

  <h1>NanoBot</h1>
  类 ZeroBot 的官方 QQ 频道适配器<br><br>

  <img src="https://counter.seku.su/cmoe?name=NanoBot&theme=r34" /><br>

</div>

## Instructions

> Note: This framework is built mainly for Chinese users thus may display hard-coded Chinese prompts during the interaction.

参见 QQ 官方[文档](https://bot.q.qq.com/wiki/)。

## 快速开始(基于插件)
> 查看`example`文件夹以获取更多信息

<table>
	<tr>
		<td align="center"><img src="https://github.com/fumiama/NanoBot/assets/41315874/6ef9fd95-ae99-449e-85e1-25797271e088"></td>
		<td align="center"><img src="https://github.com/fumiama/NanoBot/assets/41315874/edd374e4-b8a5-4cff-a463-8c3b30e537c4"></td>
        <td align="center"><img src="https://github.com/fumiama/NanoBot/assets/41315874/ed1b063f-44b0-4950-ac35-1e72745cf3f4"></td>
	</tr>
    <tr>
		<td align="center">开始响应</td>
		<td align="center">服务列表</td>
        <td align="center">查看用法</td>
	</tr>
</table>

![启用禁用](https://github.com/fumiama/NanoBot/assets/41315874/fc7f4774-f64b-44c5-9575-b9483bf3a455)


```go
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
	nano.Run(&nano.Bot{
		AppID:      "你的AppID",
		Token:      "你的Token",
		Secret:     "你的Secret, 目前没用到, 可以不填",
		Intents:    nano.IntentPublic,
		SuperUsers: []string{"用户ID1", "用户ID2"},
	})
}
```

## 更多选择(传统的事件驱动)

> 如果声明了 Handler, 所有插件将被禁用

![event-based example](https://github.com/fumiama/NanoBot/assets/41315874/414ef9a6-1da2-49ff-b28e-9e3009cdb41c)

```go
package main

import (
	"strings"

	nano "github.com/fumiama/NanoBot"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	nano.OpenAPI = nano.SandboxAPI
	nano.Run(&nano.Bot{
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
```

## Thanks

- [ZeroBot](https://github.com/wdvxdr1123/ZeroBot)
