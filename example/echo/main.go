package echo

import (
	ctrl "github.com/FloatTech/zbpctrl"
	nano "github.com/fumiama/NanoBot"
)

func init() {
	nano.Register("echo", &ctrl.Options[*nano.Ctx]{
		DisableOnDefault: false,
		Help:             "- echo xxx",
	}).OnMessagePrefix("echo").SetBlock(true).
		Handle(func(ctx *nano.Ctx) {
			args := ctx.State["args"].(string)
			if args == "" {
				return
			}
			_, _ = ctx.SendPlainMessage(false, args)
		})
}
