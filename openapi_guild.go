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
}

// GetGuildByID 获取 guild_id 指定的频道的详情
//
// https://bot.q.qq.com/wiki/develop/api/openapi/guild/get_guild.html
func (bot *Bot) GetGuildByID(id string) (*Guild, error) {
	return bot.getOpenAPIofGuild("/guilds/" + id)
}
