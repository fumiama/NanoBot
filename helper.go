package nano

import (
	"runtime"
	"strings"
)

func getFuncNameWithSkip(n int) string {
	pc, _, _, ok := runtime.Caller(n)
	if !ok {
		return ""
	}
	fullname := runtime.FuncForPC(pc).Name()
	i := strings.LastIndex(fullname, ".") + 1
	if i <= 0 || i >= len(fullname) {
		return fullname
	}
	return fullname[i:]
}

// getThisFuncName 获取正在执行的函数名
func getThisFuncName() string {
	return getFuncNameWithSkip(1)
}

// getCallerFuncName 获取调用者函数名
func getCallerFuncName() string {
	return getFuncNameWithSkip(2)
}

// MessageEscape 消息转义
//
// https://bot.q.qq.com/wiki/develop/api/openapi/message/message_format.html
func MessageEscape(text string) string {
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	return text
}

// MessageUnescape 消息解转义
//
// https://bot.q.qq.com/wiki/develop/api/openapi/message/message_format.html
func MessageUnescape(text string) string {
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	return text
}
