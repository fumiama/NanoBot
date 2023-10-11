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
