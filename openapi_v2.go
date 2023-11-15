package nano

type IDTimestampMessageResult struct {
	ID        string `json:"id"`
	Timestamp int    `json:"timestamp"`
}

// QQRobotStatus https://bot.q.qq.com/wiki/develop/api-231017/server-inter/group.html#%E4%BA%8B%E4%BB%B6
type QQRobotStatus struct {
	OpenID         string `json:"openid"`
	GroupOpenID    string `json:"group_openid"`
	OpMemberOpenID string `json:"op_member_openid"`
	Timestamp      int    `json:"timestamp"`
}
