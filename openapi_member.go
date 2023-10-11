package nano

import "time"

// Member 成员对象 Member and MemberWithGuildID
//
// https://bot.q.qq.com/wiki/develop/api/openapi/member/model.html
type Member struct {
	GuildID  string    `json:"guild_id"` // MemberWithGuildID only
	User     User      `json:"user"`
	Nick     string    `json:"nick"`
	Roles    []string  `json:"roles"`
	JoinedAt time.Time `json:"joined_at"`
	Deaf     bool      `json:"deaf"`
	Mute     bool      `json:"mute"`
	Pending  bool      `json:"pending"`
}

// GetGuildMembersIn 获取 guild_id 指定的频道中所有成员的详情列表，支持分页
//
// https://bot.q.qq.com/wiki/develop/api/openapi/member/get_members.html
func (bot *Bot) GetGuildMembersIn(id, after string, limit uint32) (members []Member, err error) {
	err = bot.GetOpenAPI(WriteHTTPQueryIfNotNil("/guilds/"+id+"/members",
		"after", after,
		"limit", limit,
	), "", &members)
	return
}

// RoleMembers 频道身份组成员列表
//
// https://bot.q.qq.com/wiki/develop/api/openapi/member/get_role_members.html#%E8%BF%94%E5%9B%9E
type RoleMembers struct {
	Data []Member `json:"data"`
	Next string   `json:"next"`
}

// GetRoleMembersOf 获取 guild_id 频道中指定role_id身份组下所有成员的详情列表，支持分页
//
// https://bot.q.qq.com/wiki/develop/api/openapi/member/get_role_members.html
func (bot *Bot) GetRoleMembersOf(guildid, roleid, startindex string, limit uint32) (*RoleMembers, error) {
	return bot.getOpenAPIofRoleMembers(WriteHTTPQueryIfNotNil("/guilds/"+guildid+"/roles/"+roleid+"/members",
		"start_index", startindex,
		"limit", limit,
	))
}

// GetGuildMemberOf 获取 guild_id 指定的频道中 user_id 对应成员的详细信息
//
// https://bot.q.qq.com/wiki/develop/api/openapi/member/get_member.html
func (bot *Bot) GetGuildMemberOf(guildid, userid string) (*Member, error) {
	return bot.getOpenAPIofMember("/guilds/" + guildid + "/members/" + userid)
}

// DeleteGuildMemberOf 删除 guild_id 指定的频道下的成员 user_id
//
// https://bot.q.qq.com/wiki/develop/api/openapi/member/delete_member.html
//
// - delhistmsgdays: 消息撤回时间范围仅支持固定的天数：3，7，15，30。 特殊的时间范围：-1: 撤回全部消息。默认值为0不撤回任何消息。
func (bot *Bot) DeleteGuildMemberOf(guildid, userid string, addblklst bool, delhistmsgdays int) error {
	return bot.DeleteOpenAPI("/guilds/"+guildid+"/members/"+userid, "", WriteBodyFromJSON(&struct {
		A bool `json:"add_blacklist"`
		D int  `json:"delete_history_msg_days"`
	}{addblklst, delhistmsgdays}))
}
