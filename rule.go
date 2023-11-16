package nano

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/floatbox/process"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
)

func newctrl(service string, o *ctrl.Options[*Ctx]) Rule {
	c := m.NewControl(service, o)
	return func(ctx *Ctx) bool {
		ctx.State["manager"] = c
		if ctx.Message != nil {
			gid, _ := strconv.ParseUint(ctx.Message.ChannelID, 10, 64)
			uid, _ := strconv.ParseUint(ctx.Message.Author.ID, 10, 64)
			return c.Handler(int64(gid), int64(uid))
		}
		return false
	}
}

func Lookup(service string) (*ctrl.Control[*Ctx], bool) {
	return m.Lookup(service)
}

// respLimiterManager 请求响应限速器管理
//
//	每 1d 4次触发
var respLimiterManager = rate.NewManager[string](time.Hour*24, 4)

func init() {
	process.NewCustomOnce(&m).Do(func() {
		OnMessageCommandGroup([]string{
			"响应", "response", "沉默", "silence",
		}, UserOrGrpAdmin).SetBlock(true).Limit(func(ctx *Ctx) *rate.Limiter {
			return respLimiterManager.Load(ctx.Message.ChannelID)
		}).secondPriority().Handle(func(ctx *Ctx) {
			grp := ctx.SenderID()
			if grp == 0 {
				return
			}

			msg := ""
			switch ctx.State["command"] {
			case "响应", "response":
				if m.CanResponse(int64(grp)) {
					msg = ctx.GetReady().User.Username + "已经在工作了哦~"
					break
				}
				err := m.Response(int64(grp))
				if err == nil {
					msg = ctx.GetReady().User.Username + "将开始在此工作啦~"
				} else {
					msg = "ERROR: " + err.Error()
				}
			case "沉默", "silence":
				if !m.CanResponse(int64(grp)) {
					msg = ctx.GetReady().User.Username + "已经在休息了哦~"
					break
				}
				err := m.Silence(int64(grp))
				if err == nil {
					msg = ctx.GetReady().User.Username + "将开始休息啦~"
				} else {
					msg = "ERROR: " + err.Error()
				}
			default:
				msg = "ERROR: bad command\"" + fmt.Sprint(ctx.State["command"]) + "\""
			}
			_, _ = ctx.SendPlainMessage(false, msg)
		})

		OnMessageCommandGroup([]string{
			"全局响应", "allresponse", "全局沉默", "allsilence",
		}, SuperUserPermission).SetBlock(true).secondPriority().Handle(func(ctx *Ctx) {
			msg := ""
			cmd := ctx.State["command"].(string)
			switch {
			case strings.Contains(cmd, "响应") || strings.Contains(cmd, "response"):
				err := m.Response(0)
				if err == nil {
					msg = ctx.GetReady().User.Username + "将开始在此工作啦~"
				} else {
					msg = "ERROR: " + err.Error()
				}
			case strings.Contains(cmd, "沉默") || strings.Contains(cmd, "silence"):
				err := m.Silence(0)
				if err == nil {
					msg = ctx.GetReady().User.Username + "将开始休息啦~"
				} else {
					msg = "ERROR: " + err.Error()
				}
			default:
				msg = "ERROR: bad command\"" + cmd + "\""
			}
			_, _ = ctx.SendPlainMessage(false, msg)
		})

		OnMessageCommandGroup([]string{
			"启用", "enable", "禁用", "disable",
		}, UserOrGrpAdmin).SetBlock(true).secondPriority().Handle(func(ctx *Ctx) {
			grp := ctx.SenderID()
			if !m.CanResponse(int64(grp)) {
				return
			}
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			service, ok := Lookup(model.Args)
			if !ok {
				_, _ = ctx.SendPlainMessage(false, "没有找到指定服务!")
				return
			}
			if strings.Contains(model.Command, "启用") || strings.Contains(model.Command, "enable") {
				service.Enable(int64(grp))
				if service.Options.OnEnable != nil {
					service.Options.OnEnable(ctx)
				} else {
					_, _ = ctx.SendPlainMessage(false, "已启用服务: ", model.Args)
				}
			} else {
				service.Disable(int64(grp))
				if service.Options.OnDisable != nil {
					service.Options.OnDisable(ctx)
				} else {
					_, _ = ctx.SendPlainMessage(false, "已禁用服务: ", model.Args)
				}
			}
		})

		OnMessageCommandGroup([]string{
			"全局启用", "allenable", "全局禁用", "alldisable",
		}, OnlyToMe, SuperUserPermission).SetBlock(true).secondPriority().Handle(func(ctx *Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			service, ok := Lookup(model.Args)
			if !ok {
				_, _ = ctx.SendPlainMessage(false, "没有找到指定服务!")
				return
			}
			if strings.Contains(model.Command, "启用") || strings.Contains(model.Command, "enable") {
				service.Enable(0)
				_, _ = ctx.SendPlainMessage(false, "已全局启用服务: ", model.Args)
			} else {
				service.Disable(0)
				_, _ = ctx.SendPlainMessage(false, "已全局禁用服务: ", model.Args)
			}
		})

		OnMessageCommandGroup([]string{"还原", "reset"}, UserOrGrpAdmin).SetBlock(true).secondPriority().Handle(func(ctx *Ctx) {
			grp := ctx.SenderID()
			if !m.CanResponse(int64(grp)) {
				return
			}
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			service, ok := Lookup(model.Args)
			if !ok {
				_, _ = ctx.SendPlainMessage(false, "没有找到指定服务!")
				return
			}
			service.Reset(int64(grp))
			_, _ = ctx.SendPlainMessage(false, "已还原服务的默认启用状态: ", model.Args)
		})

		OnMessageCommandGroup([]string{
			"禁止", "ban", "允许", "permit",
		}, AdminPermission).SetBlock(true).secondPriority().Handle(func(ctx *Ctx) {
			grp := ctx.SenderID()
			if !m.CanResponse(int64(grp)) {
				return
			}
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			args := strings.Split(model.Args, " ")
			if len(args) >= 2 {
				service, ok := Lookup(args[0])
				if !ok {
					_, _ = ctx.SendPlainMessage(false, "没有找到指定服务!")
					return
				}
				grp := ctx.SenderID()
				msg := "*" + args[0] + "报告*"
				issu := SuperUserPermission(ctx)
				if strings.Contains(model.Command, "允许") || strings.Contains(model.Command, "permit") {
					for _, usr := range args[1:] {
						uid, err := strconv.ParseInt(usr, 10, 64)
						if err == nil {
							if issu {
								service.Permit(uid, int64(grp))
								msg += "\n+ 已允许" + usr
							} else {
								member, err := ctx.GetGuildMemberOf(ctx.Message.GuildID, usr)
								if err == nil && !member.Pending {
									service.Permit(uid, int64(grp))
									msg += "\n+ 已允许" + usr
								} else {
									msg += "\nx " + usr + " 不在本群"
								}
							}
						}
					}
				} else {
					for _, usr := range args[1:] {
						uid, err := strconv.ParseInt(usr, 10, 64)
						if err == nil {
							if issu {
								service.Ban(uid, int64(grp))
								msg += "\n- 已禁止" + usr
							} else {
								member, err := ctx.GetGuildMemberOf(ctx.Message.GuildID, usr)
								if err == nil && !member.Pending {
									service.Ban(uid, int64(grp))
									msg += "\n- 已禁止" + usr
								} else {
									msg += "\nx " + usr + " 不在本群"
								}
							}
						}
					}
				}
				_, _ = ctx.SendPlainMessage(false, msg)
				return
			}
			_, _ = ctx.SendPlainMessage(false, "参数错误!")
		})

		OnMessageCommandGroup([]string{
			"全局禁止", "allban", "全局允许", "allpermit",
		}, SuperUserPermission).SetBlock(true).secondPriority().Handle(func(ctx *Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			args := strings.Split(model.Args, " ")
			if len(args) >= 2 {
				service, ok := Lookup(args[0])
				if !ok {
					_, _ = ctx.SendPlainMessage(false, "没有找到指定服务!")
					return
				}
				msg := "*" + args[0] + "全局报告*"
				if strings.Contains(model.Command, "允许") || strings.Contains(model.Command, "permit") {
					for _, usr := range args[1:] {
						uid, err := strconv.ParseInt(usr, 10, 64)
						if err == nil {
							service.Permit(uid, 0)
							msg += "\n+ 已允许" + usr
						}
					}
				} else {
					for _, usr := range args[1:] {
						uid, err := strconv.ParseInt(usr, 10, 64)
						if err == nil {
							service.Ban(uid, 0)
							msg += "\n- 已禁止" + usr
						}
					}
				}
				_, _ = ctx.SendPlainMessage(false, msg)
				return
			}
			_, _ = ctx.SendPlainMessage(false, "参数错误!")
		})

		OnMessageCommandGroup([]string{
			"封禁", "block", "解封", "unblock",
		}, SuperUserPermission).SetBlock(true).secondPriority().Handle(func(ctx *Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			args := strings.Split(model.Args, " ")
			if len(args) >= 1 {
				msg := "*报告*"
				if strings.Contains(model.Command, "解") || strings.Contains(model.Command, "un") {
					for _, usr := range args {
						uid, err := strconv.ParseInt(usr, 10, 64)
						if err == nil {
							if m.DoUnblock(uid) == nil {
								msg += "\n- 已解封" + usr
							}
						}
					}
				} else {
					for _, usr := range args {
						uid, err := strconv.ParseInt(usr, 10, 64)
						if err == nil {
							if m.DoBlock(uid) == nil {
								msg += "\n+ 已封禁" + usr
							}
						}
					}
				}
				_, _ = ctx.SendPlainMessage(false, msg)
				return
			}
			_, _ = ctx.SendPlainMessage(false, "参数错误!")
		})

		OnMessageCommandGroup([]string{
			"改变默认启用状态", "allflip",
		}, SuperUserPermission).SetBlock(true).secondPriority().Handle(func(ctx *Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			service, ok := Lookup(model.Args)
			if !ok {
				_, _ = ctx.SendPlainMessage(false, "没有找到指定服务!")
				return
			}
			err := service.Flip()
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return
			}
			_, _ = ctx.SendPlainMessage(false, "已改变全局默认启用状态: ", model.Args)
		})

		OnMessageCommandGroup([]string{"用法", "usage"}, UserOrGrpAdmin).SetBlock(true).secondPriority().
			Handle(func(ctx *Ctx) {
				model := extension.CommandModel{}
				_ = ctx.Parse(&model)
				service, ok := Lookup(model.Args)
				if !ok {
					_, _ = ctx.SendPlainMessage(false, "没有找到指定服务!")
					return
				}
				if service.Options.Help != "" {
					grp := ctx.SenderID()
					_, _ = ctx.SendPlainMessage(false, service.EnableMarkIn(int64(grp)), " ", service)
				} else {
					_, _ = ctx.SendPlainMessage(false, "该服务无帮助!")
				}
			})

		OnMessageCommandGroup([]string{"服务列表", "service_list"}, UserOrGrpAdmin).SetBlock(true).secondPriority().
			Handle(func(ctx *Ctx) {
				grp := ctx.SenderID()
				m.RLock()
				msg := make([]any, 1, len(m.M)*4+1)
				m.RUnlock()
				msg[0] = "--------服务列表--------\n发送\"/用法 name\"查看详情\n发送\"/响应\"启用会话"
				ForEachByPrio(func(i int, service *ctrl.Control[*Ctx]) bool {
					msg = append(msg, "\n", i+1, ": ", service.EnableMarkIn(int64(grp)), service.Service)
					return true
				})
				_, _ = ctx.SendPlainMessage(false, msg...)
			})

		OnMessageCommandGroup([]string{"服务详情", "service_detail"}, UserOrGrpAdmin).SetBlock(true).secondPriority().
			Handle(func(ctx *Ctx) {
				grp := ctx.SenderID()
				m.RLock()
				msgs := make([]any, 1, len(m.M)*7+1)
				m.RUnlock()
				msgs[0] = "---服务详情---\n"
				ForEachByPrio(func(i int, service *ctrl.Control[*Ctx]) bool {
					msgs = append(msgs, i+1, ": ", service.EnableMarkIn(int64(grp)), service.Service, "\n", service, "\n\n")
					return true
				})
				_, _ = ctx.SendPlainMessage(false, msgs...)
			})
	})
}
