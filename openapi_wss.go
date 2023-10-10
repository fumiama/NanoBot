package nano

// GetGeneralWSSGateway 获取通用 WSS 接入点
//
// https://bot.q.qq.com/wiki/develop/api/openapi/wss/url_get.html
func (bot *Bot) GetGeneralWSSGateway() (string, error) {
	resp := struct {
		CodeMessageBase
		U string `json:"url"`
	}{}
	err := bot.GetOpenAPI("/gateway", &resp)
	if err != nil {
		return "", err
	}
	return resp.U, nil
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
