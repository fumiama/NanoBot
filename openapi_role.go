package nano

import (
	"errors"
	"io"
)

const (
	RoleIDAll          = "1" // 全体成员
	RoleIDAdmin        = "2" // 管理员
	RoleIDCreater      = "4" // 群主/创建者
	RoleIDChannelAdmin = "5" // 子频道管理员
)

var (
	ErrMustGiveChannelID = errors.New("must give channel_id")
)

// Role 频道身份组对象
//
// https://bot.q.qq.com/wiki/develop/api/openapi/guild/role_model.html
type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Color       uint32 `json:"color"`
	Hoist       uint32 `json:"hoist"`
	Number      uint32 `json:"number"`
	MemberLimit uint32 `json:"member_limit"`
}

// GuildRoleList 频道身份组列表
//
// https://bot.q.qq.com/wiki/develop/api/openapi/guild/get_guild_roles.html#%E8%BF%94%E5%9B%9E
type GuildRoleList struct {
	GuildID      string `json:"guild_id"`
	Roles        []Role `json:"roles"`
	RoleNumLimit string `json:"role_num_limit"`
}

// GetGuildRoleListIn 获取 guild_id 指定的频道下的身份组列表
//
// https://bot.q.qq.com/wiki/develop/api/openapi/guild/get_guild_roles.html
func (bot *Bot) GetGuildRoleListIn(id string) (*GuildRoleList, error) {
	return bot.getOpenAPIofGuildRoleList("/guilds/" + id + "/roles")
}

// GuildRoleCreate 创建频道身份组响应
//
// https://bot.q.qq.com/wiki/develop/api/openapi/guild/post_guild_role.html#%E8%BF%94%E5%9B%9E
type GuildRoleCreate struct {
	RoleID string `json:"role_id"`
	Role   Role   `json:"role"`
}

// CreateGuildRoleOf 创建频道身份组
//
// https://bot.q.qq.com/wiki/develop/api/openapi/guild/post_guild_role.html
//
// 参数为非必填，但至少需要传其中之一，默认为空或 0
func (bot *Bot) CreateGuildRoleOf(id string, name string, color uint32, hoist int32) (*GuildRoleCreate, error) {
	return bot.postOpenAPIofGuildRoleCreate("/guilds/"+id+"/roles", "", WriteBodyFromJSON(&struct {
		N string `json:"name,omitempty"`
		C uint32 `json:"color,omitempty"`
		H int32  `json:"hoist,omitempty"`
	}{name, color, hoist}))
}

// GuildRolePatch 修改频道身份组
//
// https://bot.q.qq.com/wiki/develop/api/openapi/guild/patch_guild_role.html#%E8%BF%94%E5%9B%9E
type GuildRolePatch struct {
	GuildID string `json:"guild_id"`
	RoleID  string `json:"role_id"`
	Role    Role   `json:"role"`
}

// PatchGuildRoleOf 修改频道 guild_id 下 role_id 指定的身份组
//
// https://bot.q.qq.com/wiki/develop/api/openapi/guild/patch_guild_role.html
func (bot *Bot) PatchGuildRoleOf(guildid, roleid string, name string, color uint32, hoist int32) (*GuildRolePatch, error) {
	return bot.patchOpenAPIofGuildRolePatch("/guilds/"+guildid+"/roles/"+roleid, WriteBodyFromJSON(&struct {
		N string `json:"name,omitempty"`
		C uint32 `json:"color,omitempty"`
		H int32  `json:"hoist,omitempty"`
	}{name, color, hoist}))
}

// DeleteGuildRoleOf 删除频道 guild_id下 role_id 对应的身份组
//
// https://bot.q.qq.com/wiki/develop/api/openapi/guild/delete_guild_role.html
func (bot *Bot) DeleteGuildRoleOf(guildid, roleid string) error {
	return bot.DeleteOpenAPI("/guilds/"+guildid+"/roles/"+roleid, "", nil)
}

// GuildRoleChannelID 频道身份组成员返回 只填充了子频道 id 字段的对象
//
// https://bot.q.qq.com/wiki/develop/api/openapi/guild/put_guild_member_role.html#%E5%8F%82%E6%95%B0
type GuildRoleChannelID struct {
	Channel struct {
		ID string `json:"id"`
	} `json:"channel"`
}

// AddRoleToMemberOfGuild 将频道 guild_id 下的用户 user_id 添加到身份组 role_id
//
// https://bot.q.qq.com/wiki/develop/api/openapi/guild/put_guild_member_role.html
//
// 返回 channel_id
func (bot *Bot) AddRoleToMemberOfGuild(guildid, userid, roleid, channelid string) (string, error) {
	var body io.Reader
	if roleid == RoleIDChannelAdmin {
		if channelid == "" {
			return "", ErrMustGiveChannelID
		}
		body = WriteBodyFromJSON(&GuildRoleChannelID{
			Channel: struct {
				ID string `json:"id"`
			}{channelid},
		})
	}
	r, err := bot.putOpenAPIofGuildRoleChannelID("/guilds/"+guildid+"/members/"+userid+"/roles/"+roleid, body)
	if err != nil {
		return "", err
	}
	return r.Channel.ID, nil
}

// RemoveRoleFromMemberOfGuild 将用户 user_id 从 频道 guild_id 的 role_id 身份组中移除
//
// https://bot.q.qq.com/wiki/develop/api/openapi/guild/delete_guild_member_role.html
func (bot *Bot) RemoveRoleFromMemberOfGuild(guildid, userid, roleid, channelid string) error {
	var body io.Reader
	if roleid == RoleIDChannelAdmin {
		if channelid == "" {
			return ErrMustGiveChannelID
		}
		body = WriteBodyFromJSON(&GuildRoleChannelID{
			Channel: struct {
				ID string `json:"id"`
			}{channelid},
		})
	}
	return bot.DeleteOpenAPI("/guilds/"+guildid+"/members/"+userid+"/roles/"+roleid, "", body)
}
