package nano

// Channel 子频道对象
//
// https://bot.q.qq.com/wiki/develop/api/openapi/channel/model.html
type Channel struct {
	ID              string `json:"id"`
	GuildID         string `json:"guild_id"`
	Name            string `json:"name"`
	Type            int    `json:"type"`
	SubType         int    `json:"sub_type"`
	Position        int    `json:"position"`
	ParentID        string `json:"parent_id"`
	OwnerID         string `json:"owner_id"`
	PrivateType     int    `json:"private_type"`
	SpeakPermission int    `json:"speak_permission"`
	ApplicationID   string `json:"application_id"`
	Permissions     string `json:"permissions"`
}

// ChannelArray []Channel 的别名
type ChannelArray []Channel

// GetChannelsOfGuild 获取 id 指定的频道下的子频道列表
//
// https://bot.q.qq.com/wiki/develop/api/openapi/channel/get_channels.html
func (bot *Bot) GetChannelsOfGuild(id string) (*ChannelArray, error) {
	return bot.getOpenAPIofChannelArray("/guilds/" + id + "/channels")
}

// GetChannelByID 用于获取 id 指定的子频道的详情
//
// https://bot.q.qq.com/wiki/develop/api/openapi/channel/get_channel.html
func (bot *Bot) GetChannelByID(id string) (*Channel, error) {
	return bot.getOpenAPIofChannel("/channels/" + id)
}

// ChannelPost 子频道 post 操作所用对象
type ChannelPost struct {
	Name            string   `json:"name"`
	Type            int      `json:"type"`
	SubType         int      `json:"sub_type"`
	Position        int      `json:"position"`
	ParentID        string   `json:"parent_id"`
	OwnerID         string   `json:"owner_id,omitempty"`
	PrivateType     int      `json:"private_type"`
	PrivateUserIds  []string `json:"private_user_ids,omitempty"`
	SpeakPermission int      `json:"speak_permission,omitempty"`
	ApplicationID   string   `json:"application_id,omitempty"`
}
