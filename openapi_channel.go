package nano

// ChannelType https://bot.q.qq.com/wiki/develop/api/openapi/channel/model.html#channeltype
type ChannelType int

const (
	ChannelTypeText        ChannelType    = iota // 文字子频道
	ChannelTypeReserved1                         // 保留，不可用
	ChannelTypeAudio                             // 语音子频道
	ChannelTypeReserved2                         // 保留，不可用
	ChannelTypeSubchannel                        // 子频道分组
	ChannelTypeLive        = 10000 + iota        // 直播子频道
	ChannelTypeApplication                       // 应用子频道
	ChannelTypeForum                             // 论坛子频道
)

// ChannelSubType https://bot.q.qq.com/wiki/develop/api/openapi/channel/model.html#channelsubtype
type ChannelSubType int

const (
	ChannelSubTypeChat     ChannelSubType = iota // 闲聊
	ChannelSubTypeAnnounce                       // 公告
	ChannelSubTypeKouryaku                       // 攻略
	ChannelSubTypeGame                           // 开黑
)

// PrivateType https://bot.q.qq.com/wiki/develop/api/openapi/channel/model.html#privatetype
type PrivateType int

const (
	PrivateTypePublic         PrivateType = iota // 公开频道
	PrivateTypeOnlyAdmin                         // 群主管理员可见
	PrivateTypeAdminAndShimei                    // 群主管理员+指定成员，可使用 修改子频道权限接口 指定成员
)

// SpeakPermission https://bot.q.qq.com/wiki/develop/api/openapi/channel/model.html#speakpermission
type SpeakPermission int

const (
	SpeakPermissionInvalid        = iota // 无效类型
	SpeakPermissionAll                   // 所有人
	SpeakPermissionAdminAndShimei        // 群主管理员+指定成员，可使用 修改子频道权限接口 指定成员
)

// Channel 子频道对象
//
// https://bot.q.qq.com/wiki/develop/api/openapi/channel/model.html
type Channel struct {
	ID              string          `json:"id"`
	GuildID         string          `json:"guild_id"`
	Name            string          `json:"name"`
	Type            ChannelType     `json:"type"`
	SubType         ChannelSubType  `json:"sub_type"`
	Position        int             `json:"position"`
	ParentID        string          `json:"parent_id"`
	OwnerID         string          `json:"owner_id"`
	PrivateType     PrivateType     `json:"private_type"`
	SpeakPermission SpeakPermission `json:"speak_permission"`
	ApplicationID   string          `json:"application_id"` // ApplicationID see https://bot.q.qq.com/wiki/develop/api/openapi/channel/model.html#%E5%BA%94%E7%94%A8%E5%AD%90%E9%A2%91%E9%81%93%E7%9A%84%E5%BA%94%E7%94%A8%E7%B1%BB%E5%9E%8B
	Permissions     string          `json:"permissions"`
}

// ChannelArray []Channel 的别名
type ChannelArray []Channel

// GetChannelsOfGuild 获取 guild_id 指定的频道下的子频道列表
//
// https://bot.q.qq.com/wiki/develop/api/openapi/channel/get_channels.html
func (bot *Bot) GetChannelsOfGuild(id string) (*ChannelArray, error) {
	return bot.getOpenAPIofChannelArray("/guilds/" + id + "/channels")
}

// GetChannelByID 用于获取 channel_id 指定的子频道的详情
//
// https://bot.q.qq.com/wiki/develop/api/openapi/channel/get_channel.html
func (bot *Bot) GetChannelByID(id string) (*Channel, error) {
	return bot.getOpenAPIofChannel("/channels/" + id)
}

// ChannelPost 子频道 post 操作所用对象
//
// https://bot.q.qq.com/wiki/develop/api/openapi/channel/post_channels.html
type ChannelPost struct {
	Name            string          `json:"name"`
	Type            ChannelType     `json:"type"`
	SubType         ChannelSubType  `json:"sub_type"`
	Position        int             `json:"position"`
	ParentID        string          `json:"parent_id"`
	OwnerID         string          `json:"owner_id,omitempty"`
	PrivateType     PrivateType     `json:"private_type"`
	PrivateUserIds  []string        `json:"private_user_ids,omitempty"`
	SpeakPermission SpeakPermission `json:"speak_permission,omitempty"`
	ApplicationID   string          `json:"application_id,omitempty"`
}

// CreateChannel 用于在 guild_id 指定的频道下创建一个子频道
//
// https://bot.q.qq.com/wiki/develop/api/openapi/channel/post_channels.html
func (bot *Bot) CreateChannel(id string, config *ChannelPost) (*Channel, error) {
	return bot.postOpenAPIofChannel("/guilds/"+id+"/channels", WriteBodyFromJSON(config))
}

// ChannelPatch 子频道 patch 操作所用对象
//
// https://bot.q.qq.com/wiki/develop/api/openapi/channel/patch_channel.html
type ChannelPatch struct {
	Name            string           `json:"name,omitempty"`
	Position        int              `json:"position,omitempty"`
	ParentID        *string          `json:"parent_id,omitempty"`
	PrivateType     *PrivateType     `json:"private_type,omitempty"`
	SpeakPermission *SpeakPermission `json:"speak_permission,omitempty"`
}

// PatchChannel 修改 channel_id 指定的子频道的信息
//
// https://bot.q.qq.com/wiki/develop/api/openapi/channel/patch_channel.html
func (bot *Bot) PatchChannel(id string, config *ChannelPatch) (*Channel, error) {
	return bot.patchOpenAPIofChannel("/channels/"+id, WriteBodyFromJSON(config))
}

// DeleteChannel 删除 channel_id 指定的子频道
//
// https://bot.q.qq.com/wiki/develop/api/openapi/channel/delete_channel.html
func (bot *Bot) DeleteChannel(id string) error {
	return bot.DeleteOpenAPI("/channels/" + id)
}

// GetOnlineNumsInChannel 查询音视频/直播子频道 channel_id 的在线成员数
//
// https://bot.q.qq.com/wiki/develop/api/openapi/channel/get_online_nums.html
func (bot *Bot) GetOnlineNumsInChannel(id string) (int, error) {
	resp := struct {
		CodeMessageBase
		N int `json:"online_nums"`
	}{}
	err := bot.GetOpenAPI("/channels/"+id+"/online_nums", &resp)
	return resp.N, err
}
