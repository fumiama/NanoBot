package nano

// DMS 私信会话对象
//
// https://bot.q.qq.com/wiki/develop/api/openapi/dms/model.html
type DMS struct {
	GuildID    string `json:"guild_id"`
	ChannelID  string `json:"channel_id"`
	CreateTime string `json:"create_time"` // 创建私信会话时间戳
}

// CreatePrivateChat 机器人和在同一个频道内的成员创建私信会话
//
// https://bot.q.qq.com/wiki/develop/api/openapi/dms/post_dms.html
func (bot *Bot) CreatePrivateChat(guildid, userid string) (*DMS, error) {
	return bot.postOpenAPIofDMS("/users/@me/dms", "", WriteBodyFromJSON(&struct {
		R string `json:"recipient_id"`
		S string `json:"source_guild_id"`
	}{userid, guildid}))
}

// PostMessageToUser 发送私信消息，前提是已经创建了私信会话
//
// https://bot.q.qq.com/wiki/develop/api/openapi/dms/post_dms_messages.html
//
// - 私信的 guild_id 在创建私信会话时以及私信消息事件中获取
func (bot *Bot) PostMessageToUser(id string, content *MessagePost) (*Message, error) {
	return bot.postMessageTo("/dms/"+id+"/messages", content)
}

// DeleteMessageOfUser 撤回私信频道 guild_id 中 message_id 指定的私信消息, 只能用于撤回机器人自己发送的私信
//
// https://bot.q.qq.com/wiki/develop/api/openapi/dms/delete_dms.html
func (bot *Bot) DeleteMessageOfUser(guildid, messageid string, hidetip bool) error {
	return bot.DeleteOpenAPI(WriteHTTPQueryIfNotNil("/dms/"+guildid+"/messages/"+messageid,
		"hidetip", hidetip,
	), "", nil)
}
