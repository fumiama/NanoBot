package nano

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/pkg/errors"
)

const (
	// StandardAPI 正式环境接口域名
	StandardAPI = `https://api.sgroup.qq.com`
	// SandboxAPI 沙箱环境接口域名
	SandboxAPI = `https://sandbox.api.sgroup.qq.com`
	// AccessTokenAPI 获取接口凭证的 API
	AccessTokenAPI = "https://bots.qq.com/app/getAppAccessToken"
)

var (
	OpenAPI = StandardAPI // OpenAPI 实际使用的 API, 默认 StandardAPI, 可自行赋值配置
)

// CodeMessageBase 各种消息都有的 code + message 基类
type CodeMessageBase struct {
	C int    `json:"code"`
	M string `json:"message"`
}

func (bot *Bot) dohttprequest(constructer HTTPRequsetConstructer, ep, contenttype string, ptr any, body io.Reader) error {
	appid := ""
	if bot.IsV2() {
		appid = bot.AppID
	}
	req, err := constructer(ep, contenttype, bot.Authorization(), appid, body)
	if err != nil {
		return errors.Wrap(err, getCallerFuncName())
	}
	resp, err := bot.client.Do(req)
	if err != nil {
		return errors.Wrap(err, getCallerFuncName())
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}
	errsb := strings.Builder{}
	var respbase *CodeMessageBase
	if resp.StatusCode >= http.StatusBadRequest {
		errsb.WriteString("code: ")
		errsb.WriteString(resp.Status)
	}
	if ptr == nil {
		goto RET
	}
	err = json.NewDecoder(resp.Body).Decode(ptr)
	if err != nil {
		if errsb.Len() > 0 {
			errsb.WriteString(", ")
		}
		errsb.WriteString("json: ")
		errsb.WriteString(err.Error())
		goto RET
	}
	if reflect.ValueOf(ptr).Elem().Kind() == reflect.Slice {
		return nil
	}
	respbase = (*CodeMessageBase)(*(*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(&ptr), unsafe.Sizeof(uintptr(0)))))
	if respbase.C != 0 {
		if errsb.Len() > 0 {
			errsb.WriteString(", ")
		}
		errsb.WriteString("err: [")
		errsb.WriteString(strconv.Itoa(respbase.C))
		errsb.WriteString("] ")
		if len([]rune(respbase.M)) > 256 {
			errsb.WriteString(string([]rune(respbase.M)[:256]))
			errsb.WriteString("...")
		} else {
			errsb.WriteString(respbase.M)
		}
	}
RET:
	if errsb.Len() > 0 {
		return errors.Wrap(errors.New(errsb.String()), getCallerFuncName())
	}
	return nil
}

//go:generate go run codegen/getopenapiof/main.go ShardWSSGateway User Guild Channel Member RoleMembers GuildRoleList ChannelPermissions Message MessageSetting PinsMessage Schedule MessageReactionUsers

// GetOpenAPI 从 ep 获取 json 结构化数据写到 ptr, ptr 除 Slice 外必须在开头继承 CodeMessageBase
func (bot *Bot) GetOpenAPI(ep, contenttype string, ptr any) error {
	return bot.dohttprequest(NewHTTPEndpointGetRequestWithAuth, ep, contenttype, ptr, nil)
}

// GetOpenAPIWithBody 不规范地从 ep 获取 json 结构化数据写到 ptr, ptr 除 Slice 外必须在开头继承 CodeMessageBase
func (bot *Bot) GetOpenAPIWithBody(ep, contenttype string, ptr any, body io.Reader) error {
	return bot.dohttprequest(NewHTTPEndpointGetRequestWithAuth, ep, contenttype, ptr, body)
}

//go:generate go run codegen/putopenapiof/main.go GuildRoleChannelID PinsMessage

// PutOpenAPI 向 ep 发送 PUT 并获取 json 结构化数据返回写到 ptr, ptr 除 Slice 外必须在开头继承 CodeMessageBase
func (bot *Bot) PutOpenAPI(ep, contenttype string, ptr any, body io.Reader) error {
	return bot.dohttprequest(NewHTTPEndpointPutRequestWithAuth, ep, contenttype, ptr, body)
}

// DeleteOpenAPI 向 ep 发送 DELETE 请求
func (bot *Bot) DeleteOpenAPI(ep, contenttype string, body io.Reader) error {
	return bot.dohttprequest(NewHTTPEndpointDeleteRequestWithAuth, ep, contenttype, nil, body)
}

// DeleteOpenAPIWithPtr 带返回值地向 ep 发送 DELETE 请求
func (bot *Bot) DeleteOpenAPIWithPtr(ep, contenttype string, ptr any, body io.Reader) error {
	return bot.dohttprequest(NewHTTPEndpointDeleteRequestWithAuth, ep, contenttype, ptr, body)
}

//go:generate go run codegen/postopenapiof/main.go Channel GuildRoleCreate Message DMS IDTimestampMessageResult

// PostOpenAPI 从 ep 得到 json 结构化数据返回值写到 ptr, ptr 除 Slice 外必须在开头继承 CodeMessageBase
func (bot *Bot) PostOpenAPI(ep, contenttype string, ptr any, body io.Reader) error {
	return bot.dohttprequest(NewHTTPEndpointPostRequestWithAuth, ep, contenttype, ptr, body)
}

//go:generate go run codegen/patchopenapiof/main.go Channel GuildRolePatch

// PatchOpenAPI 从 ep 得到 json 结构化数据返回值写到 ptr, ptr 除 Slice 外必须在开头继承 CodeMessageBase
func (bot *Bot) PatchOpenAPI(ep, contenttype string, ptr any, body io.Reader) error {
	return bot.dohttprequest(NewHTTPEndpointPatchRequestWithAuth, ep, contenttype, ptr, body)
}
