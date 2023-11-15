package nano

import (
	"reflect"
	"regexp"
	"strings"
	"time"
)

// PrefixRule check if the text message has the prefix and trim the prefix
//
// 检查消息前缀
func PrefixRule(prefix string) Rule {
	return PrefixGroupRule(prefix)
}

// PrefixGroupRule check if the text message has the prefix and trim the prefix
//
// 检查消息前缀
func PrefixGroupRule(prefixes ...string) Rule {
	return func(ctx *Ctx) bool {
		switch msg := ctx.Value.(type) {
		case *Message:
			if msg.Content == "" { // 确保无空
				return false
			}
			for _, prefix := range prefixes {
				if strings.HasPrefix(msg.Content, MessageEscape(prefix)) {
					ctx.State["prefix"] = prefix
					arg := strings.TrimLeft(msg.Content[len(prefix):], " ")
					ctx.State["args"] = arg
					return true
				}
			}
			return false
		default:
			return false
		}
	}
}

// SuffixRule check if the text message has the suffix and trim the suffix
//
// 检查消息后缀
func SuffixRule(suffix string) Rule {
	return SuffixGroupRule(suffix)
}

// SuffixGroupRule check if the text message has the suffix and trim the suffix
//
// 检查消息后缀
func SuffixGroupRule(suffixes ...string) Rule {
	return func(ctx *Ctx) bool {
		switch msg := ctx.Value.(type) {
		case *Message:
			if msg.Content == "" { // 确保无空
				return false
			}
			for _, suffix := range suffixes {
				if strings.HasSuffix(msg.Content, MessageEscape(suffix)) {
					ctx.State["suffix"] = suffix
					arg := strings.TrimRight(msg.Content[:len(msg.Content)-len(suffix)], " ")
					ctx.State["args"] = arg
					return true
				}
			}
			return false
		default:
			return false
		}
	}
}

// CommandRule check if the message is a command and trim the command name
//
//	this rule only supports Message
func CommandRule(command string) Rule {
	return CommandGroupRule(command)
}

// CommandGroupRule check if the message is a command and trim the command name
//
//	this rule only supports Message
func CommandGroupRule(commands ...string) Rule {
	return func(ctx *Ctx) bool {
		msg, ok := ctx.Value.(*Message)
		if !ok || msg.Content == "" { // 确保无空
			return false
		}
		msg.Content = strings.TrimSpace(msg.Content)
		if msg.Content == "" { // 确保无空
			return false
		}
		cmdMessage := ""
		args := ""
		switch {
		case strings.HasPrefix(msg.Content, "/"):
			cmdMessage, args, _ = strings.Cut(msg.Content, " ")
			cmdMessage, _, _ = strings.Cut(cmdMessage, "@")
			cmdMessage = cmdMessage[1:]
		default:
			return false
		}
		for _, command := range commands {
			if strings.HasPrefix(cmdMessage, MessageEscape(command)) {
				ctx.State["command"] = command
				ctx.State["args"] = args
				return true
			}
		}
		return false
	}
}

// RegexRule check if the message can be matched by the regex pattern
func RegexRule(regexPattern string) Rule {
	regex := regexp.MustCompile(regexPattern)
	return func(ctx *Ctx) bool {
		switch msg := ctx.Value.(type) {
		case *Message:
			if msg.Content == "" { // 确保无空
				return false
			}
			if matched := regex.FindStringSubmatch(msg.Content); matched != nil {
				ctx.State["regex_matched"] = matched
				return true
			}
			return false
		default:
			return false
		}
	}
}

// ReplyRule check if the message is replying some message
//
//	this rule only supports Message
func ReplyRule(messageID string) Rule {
	return func(ctx *Ctx) bool {
		msg, ok := ctx.Value.(*Message)
		if !ok || msg.MessageReference == nil { // 确保无空
			return false
		}
		return messageID == msg.MessageReference.MessageID
	}
}

func KeywordRule(src string) Rule {
	return KeywordGroupRule(src)
}

// KeywordGroupRule check if the message has a keyword or keywords
func KeywordGroupRule(src ...string) Rule {
	return func(ctx *Ctx) bool {
		switch msg := ctx.Value.(type) {
		case *Message:
			if msg.Content == "" { // 确保无空
				return false
			}
			for _, str := range src {
				if strings.Contains(msg.Content, MessageEscape(str)) {
					ctx.State["keyword"] = str
					return true
				}
			}
			return false
		default:
			return false
		}
	}
}

// FullMatchRule check if src has the same copy of the message
func FullMatchRule(src string) Rule {
	return FullMatchGroupRule(src)
}

// FullMatchGroupRule check if src has the same copy of the message
func FullMatchGroupRule(src ...string) Rule {
	return func(ctx *Ctx) bool {
		switch msg := ctx.Value.(type) {
		case *Message:
			if msg.Content == "" { // 确保无空
				return false
			}
			for _, str := range src {
				if MessageEscape(str) == msg.Content {
					ctx.State["matched"] = msg.Content
					return true
				}
			}
			return false
		default:
			return false
		}
	}
}

// ShellRule 定义shell-like规则
//
//	this rule only supports Message
func ShellRule(cmd string, model interface{}) Rule {
	cmdRule := CommandRule(cmd)
	t := reflect.TypeOf(model)
	return func(ctx *Ctx) bool {
		if !cmdRule(ctx) {
			return false
		}
		// bind flag to struct
		args := ParseShell(ctx.State["args"].(string))
		val := reflect.New(t)
		fs := registerFlag(t, val)
		err := fs.Parse(args)
		if err != nil {
			return false
		}
		ctx.State["args"] = fs.Args()
		ctx.State["flag"] = val.Interface()
		return true
	}
}

// OnlyToMe only triggered in conditions of @bot or begin with the nicknames
//
//	this rule only supports Message
func OnlyToMe(ctx *Ctx) bool {
	return ctx.IsToMe
}

// CheckUser only triggered by specific person
func CheckUser(userID ...string) Rule {
	return func(ctx *Ctx) bool {
		switch msg := ctx.Value.(type) {
		case *Message:
			if msg.Author == nil { // 确保无空
				return false
			}
			for _, uid := range userID {
				if msg.Author.ID == uid {
					return true
				}
			}
			return false
		default:
			return false
		}
	}
}

// CheckChannel only triggered in specific channel
func CheckChannel(channelID ...string) Rule {
	return func(ctx *Ctx) bool {
		switch msg := ctx.Value.(type) {
		case *Message:
			if msg.ChannelID == "" { // 确保无空
				return false
			}
			for _, cid := range channelID {
				if msg.ChannelID == cid {
					return true
				}
			}
			return false
		default:
			return false
		}
	}
}

// CheckGuild only triggered in specific guild
func CheckGuild(guildID ...string) Rule {
	return func(ctx *Ctx) bool {
		switch msg := ctx.Value.(type) {
		case *Message:
			if msg.GuildID == "" { // 确保无空
				return false
			}
			for _, gid := range guildID {
				if msg.GuildID == gid {
					return true
				}
			}
			return false
		default:
			return false
		}
	}
}

// OnlyQQ 必须是 QQ 消息
func OnlyQQ(ctx *Ctx) bool {
	return ctx.IsQQ
}

// OnlyGuild 必须是频道消息
func OnlyGuild(ctx *Ctx) bool {
	return !ctx.IsQQ
}

// OnlyDirect 必须是频道私聊
func OnlyDirect(ctx *Ctx) bool {
	if ctx.Type != "" {
		return strings.HasPrefix(ctx.Type, "Direct")
	}
	return false
}

// OnlyChannel 必须是频道 Channel
func OnlyChannel(ctx *Ctx) bool {
	return !OnlyDirect(ctx) && !OnlyQQ(ctx)
}

// OnlyPublic 消息类型包含 At 或 Public (包括QQ群)
func OnlyPublic(ctx *Ctx) bool {
	if ctx.Type != "" {
		return strings.HasPrefix(ctx.Type, "At") || strings.HasPrefix(ctx.Type, "Public")
	}
	return false
}

// OnlyPrivate is !OnlyPublic (包括QQ私聊)
func OnlyPrivate(ctx *Ctx) bool {
	return !OnlyPublic(ctx)
}

// OnlyQQGroup 只在 QQ 群
func OnlyQQGroup(ctx *Ctx) bool {
	return ctx.Type == "GroupAtMessageCreate"
}

// OnlyQQPrivate 只在 QQ 私聊
func OnlyQQPrivate(ctx *Ctx) bool {
	return ctx.Type == "C2cMessageCreate"
}

// SuperUserPermission only triggered by the bot's owner
func SuperUserPermission(ctx *Ctx) bool {
	switch msg := ctx.Value.(type) {
	case *Message:
		if msg.Author == nil { // 确保无空
			return false
		}
		for _, su := range ctx.caller.SuperUsers {
			if su == msg.Author.ID {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// CreaterPermission only triggered by the creater or higher permission
func CreaterPermission(ctx *Ctx) bool {
	if SuperUserPermission(ctx) {
		return true
	}
	switch msg := ctx.Value.(type) {
	case *Message:
		if msg.Author == nil || msg.Member == nil { // 确保无空
			return false
		}
		for _, role := range msg.Member.Roles {
			if role == RoleIDCreater {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// AdminPermission only triggered by the admins or higher permission
func AdminPermission(ctx *Ctx) bool {
	if SuperUserPermission(ctx) {
		return true
	}
	switch msg := ctx.Value.(type) {
	case *Message:
		if msg.Author == nil || msg.Member == nil { // 确保无空
			return false
		}
		for _, role := range msg.Member.Roles {
			if role == RoleIDCreater || role == RoleIDAdmin {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// ChannelAdminPermission only triggered by the channel admins or higher permission
func ChannelAdminPermission(ctx *Ctx) bool {
	if SuperUserPermission(ctx) {
		return true
	}
	switch msg := ctx.Value.(type) {
	case *Message:
		if msg.Author == nil || msg.Member == nil { // 确保无空
			return false
		}
		for _, role := range msg.Member.Roles {
			if role == RoleIDCreater || role == RoleIDAdmin || role == RoleIDChannelAdmin {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// UserOrGrpAdmin 允许用户单独使用或群管使用
func UserOrGrpAdmin(ctx *Ctx) bool {
	if OnlyPublic(ctx) {
		return AdminPermission(ctx)
	}
	return OnlyToMe(ctx)
}

// UserOrChannelAdmin 允许用户单独使用或频道管理使用
func UserOrChannelAdmin(ctx *Ctx) bool {
	if OnlyPublic(ctx) {
		return ChannelAdminPermission(ctx)
	}
	return OnlyToMe(ctx)
}

// HasAttachments 消息包含 Attachments (典型: 图片) 返回 true
func HasAttachments(ctx *Ctx) bool {
	msg, ok := ctx.Value.(*Message)
	if !ok || len(msg.Attachments) == 0 { // 确保无空
		return false
	}
	ctx.State["attachments"] = msg.Attachments
	return true
}

// MustProvidePhoto 消息不存在图片阻塞120秒至有图片，超时返回 false
func MustProvidePhoto(onmessage string, needphohint, failhint string) Rule {
	return func(ctx *Ctx) bool {
		msg, ok := ctx.Value.(*Message)
		if ok && len(msg.Attachments) > 0 { // 确保无空
			ctx.State["attachments"] = msg.Attachments
			return true
		}
		// 没有图片就索取
		if needphohint != "" {
			_, err := ctx.PostMessageToChannel(msg.ChannelID, &MessagePost{
				Content:          needphohint,
				MessageReference: &MessageReference{MessageID: msg.ID},
				ReplyMessageID:   msg.ID,
			})
			if err != nil {
				return false
			}
		}
		next := NewFutureEvent(onmessage, 999, false, ctx.CheckSession(), HasAttachments).Next()
		select {
		case <-time.After(time.Second * 120):
			if failhint != "" {
				_, _ = ctx.SendPlainMessage(true, failhint)
			}
			return false
		case newCtx := <-next:
			ctx.State["attachments"] = newCtx.State["attachments"]
			ctx.Event = newCtx.Event
			return true
		}
	}
}
