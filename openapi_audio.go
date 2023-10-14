package nano

// AudioAction 音频事件
//
// https://bot.q.qq.com/wiki/develop/api/openapi/audio/model.html
type AudioAction struct {
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
	AudioURL  string `json:"audio_url"`
	Text      string `json:"text"`
}

// AudioControlStatus https://bot.q.qq.com/wiki/develop/api/openapi/audio/model.html#status
type AudioControlStatus int

const (
	AudioControlStatusStart AudioControlStatus = iota
	AudioControlStatusPause
	AudioControlStatusResume
	AudioControlStatusStop
)

// AudioControl 控制子频道 channel_id 下的音频
//
// https://bot.q.qq.com/wiki/develop/api/openapi/audio/audio_control.html
type AudioControl struct {
	AudioURL string             `json:"audio_url"`
	Text     string             `json:"text"`
	Status   AudioControlStatus `json:"status"`
}

// ControlAudioInChannel 控制子频道 channel_id 下的音频
//
// https://bot.q.qq.com/wiki/develop/api/openapi/audio/audio_control.html
func (bot *Bot) ControlAudioInChannel(id string, control *AudioControl) error {
	return bot.PostOpenAPI("/channels/"+id+"/audio", "", &CodeMessageBase{}, WriteBodyFromJSON(control))
}

// OpenMic 机器人在 channel_id 对应的语音子频道上麦
//
// https://bot.q.qq.com/wiki/develop/api/openapi/audio/put_mic.html
func (bot *Bot) OpenMicInChannel(id string) error {
	return bot.PutOpenAPI("/channels/"+id+"/mic", "", &CodeMessageBase{}, nil)
}

// CloseMicInChannel 机器人在 channel_id 对应的语音子频道下麦
//
// https://bot.q.qq.com/wiki/develop/api/openapi/audio/delete_mic.html
func (bot *Bot) CloseMicInChannel(id string) error {
	return bot.DeleteOpenAPI("/channels/"+id+"/mic", "", nil)
}

// AudioLiveChannelUsersChange 音视频/直播子频道成员进出事件
//
// https://bot.q.qq.com/wiki/develop/api/gateway/audio_or_live_channel_member.html
type AudioLiveChannelUsersChange struct {
	GuildID     string `json:"guild_id"`
	ChannelID   string `json:"channel_id"`
	ChannelType int    `json:"channel_type"`
	UserID      string `json:"user_id"`
}
