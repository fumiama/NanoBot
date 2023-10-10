package nano

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"unsafe"

	"github.com/pkg/errors"
)

//go:generate go run codegen/getopenapiof/main.go ShardWSSGateway User Guild GuildArray Channel ChannelArray

// GetOpenAPI 从 ep 获取 json 结构化数据写到 ptr, ptr 必须在开头继承 CodeMessageBase
func (bot *Bot) GetOpenAPI(ep string, ptr any) error {
	req, err := NewHTTPEndpointGetRequestWithAuth(ep, bot.Authorization())
	if err != nil {
		return errors.Wrap(err, getCallerFuncName())
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, getCallerFuncName())
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(ptr)
	if err != nil {
		return errors.Wrap(err, getCallerFuncName())
	}
	respbbase := (*CodeMessageBase)(*(*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(&ptr), unsafe.Sizeof(uintptr(0)))))
	if respbbase.C != 0 {
		return errors.Wrap(errors.New("code: "+strconv.Itoa(respbbase.C)+", msg: "+respbbase.M), getCallerFuncName())
	}
	return nil
}

//go:generate go run codegen/postopenapiof/main.go Channel

// PostOpenAPI 从 ep 得到 json 结构化数据返回值写到 ptr, ptr 必须在开头继承 CodeMessageBase
func (bot *Bot) PostOpenAPI(ep string, ptr any, body io.Reader) error {
	req, err := NewHTTPEndpointPostRequestWithAuth(ep, bot.Authorization(), body)
	if err != nil {
		return errors.Wrap(err, getCallerFuncName())
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, getCallerFuncName())
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(ptr)
	if err != nil {
		return errors.Wrap(err, getCallerFuncName())
	}
	respbbase := (*CodeMessageBase)(*(*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(&ptr), unsafe.Sizeof(uintptr(0)))))
	if respbbase.C != 0 {
		return errors.Wrap(errors.New("code: "+strconv.Itoa(respbbase.C)+", msg: "+respbbase.M), getCallerFuncName())
	}
	return nil
}
