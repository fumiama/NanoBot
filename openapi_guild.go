package nano

import "time"

// Guild 频道对象
//
// https://bot.q.qq.com/wiki/develop/api/openapi/guild/model.html
type Guild struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Icon        string    `json:"icon"`
	OwnerID     string    `json:"owner_id"`
	Owner       bool      `json:"owner"`
	JoinedAt    time.Time `json:"joined_at"`
	MemberCount int       `json:"member_count"`
	MaxMembers  int       `json:"max_members"`
	Description string    `json:"description"`
	OpUserID    string    `json:"op_user_id"` // https://bot.q.qq.com/wiki/develop/api/gateway/guild.html#%E4%BA%8B%E4%BB%B6%E5%86%85%E5%AE%B9
}

// GetGuildByID 获取 guild_id 指定的频道的详情
//
// https://bot.q.qq.com/wiki/develop/api/openapi/guild/get_guild.html
func (bot *Bot) GetGuildByID(id string) (*Guild, error) {
	return bot.getOpenAPIofGuild("/guilds/" + id)
}

// SetAllMuteInGuild 禁言全员 / 解除全员禁言
//
// https://bot.q.qq.com/wiki/develop/api/openapi/guild/patch_guild_mute.html
func (bot *Bot) SetAllMuteInGuild(id string, endtimestamp string, seconds string) error {
	return bot.PatchOpenAPI("/guilds/"+id+"/mute", "", nil, WriteBodyFromJSON(&struct {
		T string `json:"mute_end_timestamp"`
		S string `json:"mute_seconds"`
	}{endtimestamp, seconds}))
}

// SetUserMuteInGuild 禁言 / 解除禁言频道 guild_id 下的成员 user_id
//
// https://bot.q.qq.com/wiki/develop/api/openapi/guild/patch_guild_mute.html
func (bot *Bot) SetUserMuteInGuild(guildid, userid string, endtimestamp string, seconds string) error {
	return bot.PatchOpenAPI("/guilds/"+guildid+"/members/"+userid+"/mute", "", nil, WriteBodyFromJSON(&struct {
		T string `json:"mute_end_timestamp"`
		S string `json:"mute_seconds"`
	}{endtimestamp, seconds}))
}

// SetUsersMuteInGuild 批量禁言 / 解除禁言频道 guild_id 下的成员 user_id
//
// https://bot.q.qq.com/wiki/develop/api/openapi/guild/patch_guild_mute.html
func (bot *Bot) SetUsersMuteInGuild(guildid string, endtimestamp string, seconds string, userids ...string) ([]string, error) {
	resp := &struct {
		CodeMessageBase
		U []string `json:"user_ids"`
	}{}
	err := bot.PatchOpenAPI("/guilds/"+guildid+"/mute", "", resp, WriteBodyFromJSON(&struct {
		T string   `json:"mute_end_timestamp"`
		S string   `json:"mute_seconds"`
		U []string `json:"user_ids"`
	}{endtimestamp, seconds, userids}))
	return resp.U, err
}
