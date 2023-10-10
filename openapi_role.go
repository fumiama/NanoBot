package nano

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

// GetGuildRoleList 获取 guild_id指定的频道下的身份组列表
//
// https://bot.q.qq.com/wiki/develop/api/openapi/guild/get_guild_roles.html
func (bot *Bot) GetGuildRoleList(id string) (*GuildRoleList, error) {
	return bot.getOpenAPIofGuildRoleList("/guilds/" + id + "/roles")
}

// GuildRoleCreate 创建频道身份组响应
//
// https://bot.q.qq.com/wiki/develop/api/openapi/guild/post_guild_role.html#%E8%BF%94%E5%9B%9E
type GuildRoleCreate struct {
	RoleID string `json:"role_id"`
	Role   Role   `json:"role"`
}

// CreateGuildRole 创建频道身份组
//
// https://bot.q.qq.com/wiki/develop/api/openapi/guild/post_guild_role.html
//
// 参数为非必填，但至少需要传其中之一，默认为空或 0
func (bot *Bot) CreateGuildRole(id string, name string, color uint32, hoist int32) (*GuildRoleCreate, error) {
	return bot.postOpenAPIofGuildRoleCreate("/guilds/"+id+"/roles", WriteBodyFromJSON(&struct {
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

// PatchGuildRole 修改频道 guild_id 下 role_id 指定的身份组
//
// https://bot.q.qq.com/wiki/develop/api/openapi/guild/patch_guild_role.html
func (bot *Bot) PatchGuildRole(guildid, roleid string, name string, color uint32, hoist int32) (*GuildRolePatch, error) {
	return bot.patchOpenAPIofGuildRolePatch("/guilds/"+guildid+"/roles/"+roleid, WriteBodyFromJSON(&struct {
		N string `json:"name,omitempty"`
		C uint32 `json:"color,omitempty"`
		H int32  `json:"hoist,omitempty"`
	}{name, color, hoist}))
}

func (bot *Bot) DeleteGuildRole(guildid, roleid string) error {
	return bot.DeleteOpenAPI("/guilds/"+guildid+"/roles/"+roleid, nil)
}
