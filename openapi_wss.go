package nano

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

var (
	ErrEmptyToken    = errors.New("empty token")
	ErrInvalidExpire = errors.New("invalid expire")
)

// GetGeneralWSSGateway 获取通用 WSS 接入点
//
// https://bot.q.qq.com/wiki/develop/api/openapi/wss/url_get.html
func (bot *Bot) GetGeneralWSSGateway() (string, error) {
	resp := struct {
		CodeMessageBase
		U string `json:"url"`
	}{}
	err := bot.GetOpenAPI("/gateway", "", &resp)
	return resp.U, err
}

// ShardWSSGateway 带分片 WSS 接入点响应数据
//
// https://bot.q.qq.com/wiki/develop/api/openapi/wss/shard_url_get.html#%E8%BF%94%E5%9B%9E
type ShardWSSGateway struct {
	URL               string `json:"url"`
	Shards            int    `json:"shards"`
	SessionStartLimit struct {
		Total          int `json:"total"`
		Remaining      int `json:"remaining"`
		ResetAfter     int `json:"reset_after"`
		MaxConcurrency int `json:"max_concurrency"`
	} `json:"session_start_limit"`
}

// GetShardWSSGateway 获取带分片 WSS 接入点
//
// https://bot.q.qq.com/wiki/develop/api/openapi/wss/shard_url_get.html
func (bot *Bot) GetShardWSSGateway() (*ShardWSSGateway, error) {
	return bot.getOpenAPIofShardWSSGateway("/gateway/bot")
}

// GetAppAccessToken 获取接口凭证并保存到 bot.Token
//
// https://bot.q.qq.com/wiki/develop/api-231017/dev-prepare/interface-framework/api-use.html#%E8%8E%B7%E5%8F%96%E6%8E%A5%E5%8F%A3%E5%87%AD%E8%AF%81
func (bot *Bot) GetAppAccessToken() error {
	req, err := newHTTPEndpointRequestWithAuth("POST", "", AccessTokenAPI, "", "", WriteBodyFromJSON(&struct {
		A string `json:"appId"`
		S string `json:"clientSecret"`
	}{bot.AppID, bot.Secret}))
	if err != nil {
		return err
	}
	resp, err := bot.client.Do(req)
	if err != nil {
		return errors.Wrap(err, getThisFuncName())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.Wrap(errors.New("http status code: "+strconv.Itoa(resp.StatusCode)), getThisFuncName())
	}
	body := struct {
		C int    `json:"code"`
		M string `json:"message"`
		T string `json:"access_token"`
		E string `json:"expires_in"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return errors.Wrap(err, getThisFuncName())
	}
	if body.C != 0 {
		return errors.Wrap(errors.New("code: "+strconv.Itoa(body.C)+", msg: "+body.M), getThisFuncName())
	}
	if body.T == "" {
		return errors.Wrap(ErrEmptyToken, getThisFuncName())
	}
	if body.E == "" {
		return errors.Wrap(ErrInvalidExpire, getThisFuncName())
	}
	bot.token = body.T
	bot.expiresec, err = strconv.ParseInt(body.E, 10, 64)
	if err != nil {
		return errors.Wrap(err, getThisFuncName())
	}
	if bot.expiresec <= 0 {
		return errors.Wrap(ErrInvalidExpire, getThisFuncName())
	}
	return nil
}
