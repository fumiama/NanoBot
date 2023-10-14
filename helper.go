package nano

import (
	"encoding/base64"
	"net/url"
	"runtime"
	"strings"
	"unsafe"
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

// UnderlineToCamel convert abc_def to AbcDef
func UnderlineToCamel(s string) string {
	sb := strings.Builder{}
	isnextupper := true
	for _, c := range []byte(strings.ToLower(s)) {
		if c == '_' {
			isnextupper = true
			continue
		}
		if isnextupper {
			sb.WriteString(strings.ToUpper(string(c)))
			isnextupper = false
			continue
		}
		sb.WriteByte(c)
	}
	return sb.String()
}

// resolveURI github.com/wdvxdr1123/ZeroBot/driver/uri.go
func resolveURI(addr string) (network, address string) {
	network, address = "tcp", addr
	uri, err := url.Parse(addr)
	if err == nil && uri.Scheme != "" {
		scheme, ext, _ := strings.Cut(uri.Scheme, "+")
		if ext != "" {
			network = ext
			uri.Scheme = scheme // remove `+unix`/`+tcp4`
			if ext == "unix" {
				uri.Host, uri.Path, _ = strings.Cut(uri.Path, ":")
				uri.Host = base64.StdEncoding.EncodeToString(StringToBytes(uri.Host)) // special handle for unix
			}
			address = uri.String()
		}
	}
	return
}

// slice is the runtime representation of a slice.
// It cannot be used safely or portably and its representation may
// change in a later release.
//
// Unlike reflect.SliceHeader, its Data field is sufficient to guarantee the
// data it references will not be garbage collected.
type slice struct {
	data unsafe.Pointer
	len  int
	cap  int
}

// BytesToString 没有内存开销的转换
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes 没有内存开销的转换
func StringToBytes(s string) (b []byte) {
	bh := (*slice)(unsafe.Pointer(&b))
	sh := (*slice)(unsafe.Pointer(&s))
	bh.data = sh.data
	bh.len = sh.len
	bh.cap = sh.len
	return b
}
