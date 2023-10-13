package nano

import "strconv"

// Emoji 表情对象
//
// https://bot.q.qq.com/wiki/develop/api/openapi/emoji/model.html
type Emoji struct {
	ID   string `json:"id"`
	Type uint32 `json:"type"`
}

// MessageReaction https://bot.q.qq.com/wiki/develop/api/openapi/reaction/model.html#messagereaction
type MessageReaction struct {
	UserID    string          `json:"user_id"`
	GuildID   string          `json:"guild_id"`
	ChannelID string          `json:"channel_id"`
	Target    *ReactionTarget `json:"target"`
	Emoji     *Emoji          `json:"emoji"`
}

// ReactionTargetType https://bot.q.qq.com/wiki/develop/api/openapi/reaction/model.html#reactiontargettype
type ReactionTargetType int

// ReactionTarget https://bot.q.qq.com/wiki/develop/api/openapi/reaction/model.html#reactiontarget
type ReactionTarget struct {
	ID   string             `json:"id"`
	Type ReactionTargetType `json:"type"`
}

// GiveMessageReaction 对消息 message_id 进行表情表态
//
// https://bot.q.qq.com/wiki/develop/api/openapi/reaction/put_message_reaction.html
func (bot *Bot) GiveMessageReaction(channelid, messageid string, emoji Emoji) error {
	return bot.PutOpenAPI("/channels/"+channelid+"/messages/"+messageid+"/reactions/"+strconv.FormatUint(uint64(emoji.Type), 10)+"/"+emoji.ID, "", nil, nil)
}

// DeleteMessageReaction 删除自己对消息 message_id 的表情表态
//
// https://bot.q.qq.com/wiki/develop/api/openapi/reaction/delete_own_message_reaction.html
func (bot *Bot) DeleteMessageReaction(channelid, messageid string, emoji Emoji) error {
	return bot.DeleteOpenAPI("/channels/"+channelid+"/messages/"+messageid+"/reactions/"+strconv.FormatUint(uint64(emoji.Type), 10)+"/"+emoji.ID, "", nil)
}

// MessageReactionUsers https://bot.q.qq.com/wiki/develop/api/openapi/reaction/get_reaction_users.html#%E8%BF%94%E5%9B%9E
type MessageReactionUsers struct {
	Users  []User `json:"users"`
	Cookie string `json:"cookie"`
	IsEnd  bool   `json:"is_end"`
}

// GetMessageReactionUsers 拉取对消息 message_id 指定表情表态的用户列表
//
// https://bot.q.qq.com/wiki/develop/api/openapi/reaction/get_reaction_users.html
func (bot *Bot) GetMessageReactionUsers(channelid, messageid string, emoji Emoji, cookie string, limit int) (*MessageReactionUsers, error) {
	return bot.getOpenAPIofMessageReactionUsers(WriteHTTPQueryIfNotNil(
		"/channels/"+channelid+"/messages/"+messageid+"/reactions/"+strconv.FormatUint(uint64(emoji.Type), 10)+"/"+emoji.ID,
		"cookie", cookie,
		"limit", limit,
	))
}
