package nano

// PinsMessage 精华消息对象
//
// https://bot.q.qq.com/wiki/develop/api/openapi/pins/model.html#pinsmessage
type PinsMessage struct {
	GuildID    string   `json:"guild_id"`
	ChannelID  string   `json:"channel_id"`
	MessageIDs []string `json:"message_ids"`
}

// PinMessageInChannel 添加子频道 channel_id 内的精华消息
//
// https://bot.q.qq.com/wiki/develop/api/openapi/pins/put_pins_message.html
func (bot *Bot) PinMessageInChannel(channelid, messageid string) (*PinsMessage, error) {
	return bot.putOpenAPIofPinsMessage("/channels/"+channelid+"/pins/"+messageid, nil)
}

// UnpinMessageInChannel 子频道 channel_id 下指定 message_id 的精华消息
//
// https://bot.q.qq.com/wiki/develop/api/openapi/pins/delete_pins_message.html
//
// 删除子频道内全部精华消息，请将 message_id 设置为 all
func (bot *Bot) UnpinMessageInChannel(channelid, messageid string) error {
	return bot.DeleteOpenAPI("/channels/"+channelid+"/pins/"+messageid, "", nil)
}

// GetPinMessagesOfChannel 获取子频道 channel_id 内的精华消息
//
// https://bot.q.qq.com/wiki/develop/api/openapi/pins/get_pins_message.html
func (bot *Bot) GetPinMessagesOfChannel(id string) (*PinsMessage, error) {
	return bot.getOpenAPIofPinsMessage("/channels/" + id + "/pins")
}
