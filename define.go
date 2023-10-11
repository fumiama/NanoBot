package nano

const (
	// StandardAPI 正式环境接口域名
	StandardAPI = `https://api.sgroup.qq.com`
	// SandboxAPI 沙箱环境接口域名
	SandboxAPI = `https://sandbox.api.sgroup.qq.com`
)

var (
	OpenAPI = StandardAPI // OpenAPI 实际使用的 API, 可自行配置
)

// CodeMessageBase 各种消息都有的 code + message 基类
type CodeMessageBase struct {
	C int    `json:"code"`
	M string `json:"message"`
}
