package nano

// User 用户对象
//
// https://bot.q.qq.com/wiki/develop/api/openapi/user/model.html
type User struct {
	ID               string `json:"id"`
	Username         string `json:"username"`
	Avatar           string `json:"avatar"`
	Bot              bool   `json:"bot"`
	UnionOpenid      string `json:"union_openid"`
	UnionUserAccount string `json:"union_user_account"`
}

// GetMyInfo 获取当前用户（机器人）详情
//
// https://bot.q.qq.com/wiki/develop/api/openapi/user/me.html
func (bot *Bot) GetMyInfo() (*User, error) {
	return bot.getOpenAPIofUser("/users/@me")
}

// GetMyGuilds 获取当前用户（机器人）频道列表，支持分页
//
// https://bot.q.qq.com/wiki/develop/api/openapi/user/guilds.html
func (bot *Bot) GetMyGuilds(before, after string, limit int) (guilds []Guild, err error) {
	err = bot.GetOpenAPI(WriteHTTPQueryIfNotNil("/users/@me/guilds",
		"before", before,
		"after", after,
		"limit", limit,
	), &guilds)
	return
}
