package nano

import (
	"runtime"
	"strings"
)

// getCurrentFuncName 获取当前函数名
func getCurrentFuncName() string {
	pc, _, _, ok := runtime.Caller(1)
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

// getCallerFuncName 获取调用者函数名
func getCallerFuncName() string {
	pc, _, _, ok := runtime.Caller(2)
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
