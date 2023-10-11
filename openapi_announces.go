package nano

// Announces 公告对象
//
// https://bot.q.qq.com/wiki/develop/api/openapi/announces/model.html#announces
type Announces struct {
	GuildID           string             `json:"guild_id,omitempty"`
	ChannelID         string             `json:"channel_id,omitempty"`
	MessageID         string             `json:"message_id,omitempty"`
	AnnouncesType     uint32             `json:"announces_type,omitempty"` // 公告类别 0:成员公告 1:欢迎公告，默认成员公告
	RecommendChannels []RecommendChannel `json:"recommend_channels,omitempty"`
}

// RecommendChannel 推荐子频道对象
//
// https://bot.q.qq.com/wiki/develop/api/openapi/announces/model.html#recommendchannel
type RecommendChannel struct {
	ChannelID string `json:"channel_id"`
	Introduce string `json:"introduce"`
}

// PostAnnounceInGuild 创建频道全局公告，公告类型分为 消息类型的频道公告 和 推荐子频道类型的频道公告
//
// https://bot.q.qq.com/wiki/develop/api/openapi/announces/post_guild_announces.html
//
// 会重写 content 为返回值
func (bot *Bot) PostAnnounceInGuild(id string, content *Announces) error {
	return bot.PostOpenAPI("/guilds/"+id+"/announces", "", content, WriteBodyFromJSON(content))
}

// DeleteAnnounceInGuild 删除频道 guild_id 下指定 message_id 的全局公告
//
// https://bot.q.qq.com/wiki/develop/api/openapi/announces/delete_guild_announces.html
//
// message_id 有值时，会校验 message_id 合法性，若不校验校验 message_id，请将 message_id 设置为 all
func (bot *Bot) DeleteAnnounceInGuild(guildid, messageid string) error {
	return bot.DeleteOpenAPI("/guilds/"+guildid+"/announces/"+messageid, "", nil)
}
