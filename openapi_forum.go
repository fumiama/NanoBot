package nano

import "time"

// Thread 话题频道内发表的主帖称为主题
//
// https://bot.q.qq.com/wiki/develop/api/openapi/forum/model.html#thread
type Thread struct {
	GuildID    string      `json:"guild_id"`
	ChannelID  string      `json:"channel_id"`
	AuthorID   string      `json:"author_id"`
	ThreadInfo *ThreadInfo `json:"thread_info"`
}

// ThreadInfo 帖子事件包含的主帖内容相关信息
//
// https://bot.q.qq.com/wiki/develop/api/openapi/forum/model.html#threadinfo
type ThreadInfo struct {
	ThreadID string    `json:"thread_id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	DateTime time.Time `json:"date_time"`
}

// Post 话题频道内对主题的评论称为帖子
//
// https://bot.q.qq.com/wiki/develop/api/openapi/forum/model.html#post
type Post struct {
	GuildID   string    `json:"guild_id"`
	ChannelID string    `json:"channel_id"`
	AuthorID  string    `json:"author_id"`
	PostInfo  *PostInfo `json:"post_info"`
}

// PostInfo 帖子事件包含的帖子内容信息
//
// https://bot.q.qq.com/wiki/develop/api/openapi/forum/model.html#postinfo
type PostInfo struct {
	ThreadID string    `json:"thread_id"`
	PostID   string    `json:"post_id"`
	Content  string    `json:"content"`
	DateTime time.Time `json:"date_time"`
}

// Reply 话题频道对帖子回复或删除时生产该事件中包含该对象
//
// https://bot.q.qq.com/wiki/develop/api/openapi/forum/model.html#reply
type Reply struct {
	GuildID   string     `json:"guild_id"`
	ChannelID string     `json:"channel_id"`
	AuthorID  string     `json:"author_id"`
	ReplyInfo *ReplyInfo `json:"reply_info"`
}

// ReplyInfo 回复事件包含的回复内容信息
//
// https://bot.q.qq.com/wiki/develop/api/openapi/forum/model.html#replyinfo
type ReplyInfo struct {
	ThreadID string    `json:"thread_id"`
	PostID   string    `json:"post_id"`
	ReplyID  string    `json:"reply_id"`
	Content  string    `json:"content"`
	DateTime time.Time `json:"date_time"`
}

// AuditResult 论坛帖子审核结果事件
//
// https://bot.q.qq.com/wiki/develop/api/openapi/forum/model.html#auditresult
type AuditResult struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
	AuthorID  string `json:"author_id"`
	ThreadID  string `json:"thread_id"`
	PostID    string `json:"post_id"`
	ReplyID   string `json:"reply_id"`
	Type      uint32 `json:"type"`
	Result    uint32 `json:"result"`
	ErrMsg    string `json:"err_msg"`
}

// GetChannelThreads 获取子频道下的帖子列表
//
// https://bot.q.qq.com/wiki/develop/api/openapi/forum/get_threads_list.html
func (bot *Bot) GetChannelThreads(id string) (threads []Thread, isfinish bool, err error) {
	resp := &struct {
		CodeMessageBase
		T []Thread `json:"threads"`
		I uint32   `json:"is_finish"`
	}{}
	err = bot.GetOpenAPI("/channels/"+id+"/threads", "", resp)
	threads = resp.T
	isfinish = resp.I > 0
	return
}

// GetThreadInfo 获取子频道下的帖子详情
//
// https://bot.q.qq.com/wiki/develop/api/openapi/forum/get_thread.html
func (bot *Bot) GetThreadInfo(channelid, threadid string) (*ThreadInfo, error) {
	resp := &struct {
		CodeMessageBase
		T ThreadInfo `json:"thread"`
	}{}
	err := bot.GetOpenAPI("/channels/"+channelid+"/threads/"+threadid, "", resp)
	return &resp.T, err
}

// PostThread 发表帖子
//
// https://bot.q.qq.com/wiki/develop/api/openapi/forum/put_thread.html
func (bot *Bot) PostThreadInChannel(id string, title string, content string, format uint32) (taskid string, createtime string, err error) {
	resp := &struct {
		CodeMessageBase
		T string `json:"task_id"`
		C string `json:"create_time"`
	}{}
	err = bot.PutOpenAPI("/channels/"+id+"/threads", "", resp, WriteBodyFromJSON(&struct {
		T string `json:"title"`
		C string `json:"content"`
		F uint32 `json:"format"`
	}{title, content, format}))
	taskid = resp.T
	createtime = resp.C
	return
}

// DeleteThreadInChannel 删除指定子频道下的某个帖子
//
// https://bot.q.qq.com/wiki/develop/api/openapi/forum/delete_thread.html
func (bot *Bot) DeleteThreadInChannel(channelid, threadid string) error {
	return bot.DeleteOpenAPI("/channels/"+channelid+"/threads/"+threadid, "", nil)
}
