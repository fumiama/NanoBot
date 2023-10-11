package nano

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/pkg/errors"
)

//go:generate go run codegen/getopenapiof/main.go ShardWSSGateway User Guild Channel Member RoleMembers GuildRoleList ChannelPermissions Message

// GetOpenAPI 从 ep 获取 json 结构化数据写到 ptr, ptr 除 Slice 外必须在开头继承 CodeMessageBase
func (bot *Bot) GetOpenAPI(ep, contenttype string, ptr any) error {
	req, err := NewHTTPEndpointGetRequestWithAuth(ep, contenttype, bot.Authorization())
	if err != nil {
		return errors.Wrap(err, getCallerFuncName())
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, getCallerFuncName())
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Wrap(errors.New("code: "+strconv.Itoa(resp.StatusCode)+", msg: "+resp.Status), getCallerFuncName())
	}
	if ptr == nil {
		return nil
	}
	err = json.NewDecoder(resp.Body).Decode(ptr)
	if err != nil {
		return errors.Wrap(err, getCallerFuncName())
	}
	if reflect.ValueOf(ptr).Elem().Kind() == reflect.Slice {
		return nil
	}
	respbase := (*CodeMessageBase)(*(*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(&ptr), unsafe.Sizeof(uintptr(0)))))
	if respbase.C != 0 {
		return errors.Wrap(errors.New("code: "+strconv.Itoa(respbase.C)+", msg: "+respbase.M), getCallerFuncName())
	}
	return nil
}

//go:generate go run codegen/putopenapiof/main.go GuildRoleChannelID

// PutOpenAPI 向 ep 发送 PUT 并获取 json 结构化数据返回写到 ptr, ptr 除 Slice 外必须在开头继承 CodeMessageBase
func (bot *Bot) PutOpenAPI(ep, contenttype string, ptr any, body io.Reader) error {
	req, err := NewHTTPEndpointPutRequestWithAuth(ep, contenttype, bot.Authorization(), body)
	if err != nil {
		return errors.Wrap(err, getCallerFuncName())
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, getCallerFuncName())
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Wrap(errors.New("code: "+strconv.Itoa(resp.StatusCode)+", msg: "+resp.Status), getCallerFuncName())
	}
	if ptr == nil {
		return nil
	}
	err = json.NewDecoder(resp.Body).Decode(ptr)
	if err != nil {
		return errors.Wrap(err, getCallerFuncName())
	}
	if reflect.ValueOf(ptr).Elem().Kind() == reflect.Slice {
		return nil
	}
	respbase := (*CodeMessageBase)(*(*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(&ptr), unsafe.Sizeof(uintptr(0)))))
	if respbase.C != 0 {
		return errors.Wrap(errors.New("code: "+strconv.Itoa(respbase.C)+", msg: "+respbase.M), getCallerFuncName())
	}
	return nil
}

// DeleteOpenAPI 向 ep 发送 DELETE 请求
func (bot *Bot) DeleteOpenAPI(ep, contenttype string, body io.Reader) error {
	req, err := NewHTTPEndpointDeleteRequestWithAuth(ep, contenttype, bot.Authorization(), body)
	if err != nil {
		return errors.Wrap(err, getCallerFuncName())
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, getCallerFuncName())
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Wrap(errors.New("code: "+strconv.Itoa(resp.StatusCode)+", msg: "+resp.Status), getCallerFuncName())
	}
	return nil
}

//go:generate go run codegen/postopenapiof/main.go Channel GuildRoleCreate Message

// PostOpenAPI 从 ep 得到 json 结构化数据返回值写到 ptr, ptr 除 Slice 外必须在开头继承 CodeMessageBase
func (bot *Bot) PostOpenAPI(ep, contenttype string, ptr any, body io.Reader) error {
	req, err := NewHTTPEndpointPostRequestWithAuth(ep, contenttype, bot.Authorization(), body)
	if err != nil {
		return errors.Wrap(err, getCallerFuncName())
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, getCallerFuncName())
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Wrap(errors.New("code: "+strconv.Itoa(resp.StatusCode)+", msg: "+resp.Status), getCallerFuncName())
	}
	if ptr == nil {
		return nil
	}
	err = json.NewDecoder(resp.Body).Decode(ptr)
	if err != nil {
		return errors.Wrap(err, getCallerFuncName())
	}
	if reflect.ValueOf(ptr).Elem().Kind() == reflect.Slice {
		return nil
	}
	respbase := (*CodeMessageBase)(*(*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(&ptr), unsafe.Sizeof(uintptr(0)))))
	if respbase.C != 0 {
		return errors.Wrap(errors.New("code: "+strconv.Itoa(respbase.C)+", msg: "+respbase.M), getCallerFuncName())
	}
	return nil
}

//go:generate go run codegen/patchopenapiof/main.go Channel GuildRolePatch

// PatchOpenAPI 从 ep 得到 json 结构化数据返回值写到 ptr, ptr 除 Slice 外必须在开头继承 CodeMessageBase
func (bot *Bot) PatchOpenAPI(ep, contenttype string, ptr any, body io.Reader) error {
	req, err := NewHTTPEndpointPatchRequestWithAuth(ep, contenttype, bot.Authorization(), body)
	if err != nil {
		return errors.Wrap(err, getCallerFuncName())
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, getCallerFuncName())
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Wrap(errors.New("code: "+strconv.Itoa(resp.StatusCode)+", msg: "+resp.Status), getCallerFuncName())
	}
	if ptr == nil {
		return nil
	}
	err = json.NewDecoder(resp.Body).Decode(ptr)
	if err != nil {
		return errors.Wrap(err, getCallerFuncName())
	}
	if reflect.ValueOf(ptr).Elem().Kind() == reflect.Slice {
		return nil
	}
	respbase := (*CodeMessageBase)(*(*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(&ptr), unsafe.Sizeof(uintptr(0)))))
	if respbase.C != 0 {
		return errors.Wrap(errors.New("code: "+strconv.Itoa(respbase.C)+", msg: "+respbase.M), getCallerFuncName())
	}
	return nil
}
