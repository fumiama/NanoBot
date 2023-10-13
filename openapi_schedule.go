package nano

// Schedule 日程对象
//
// https://bot.q.qq.com/wiki/develop/api/openapi/schedule/model.html
type Schedule struct {
	ID             string  `json:"id,omitempty"`
	Name           string  `json:"name,omitempty"`
	Description    string  `json:"description,omitempty"`
	StartTimestamp string  `json:"start_timestamp,omitempty"`
	EndTimestamp   string  `json:"end_timestamp,omitempty"`
	Creator        *Member `json:"creator,omitempty"`
	JumpChannelID  string  `json:"jump_channel_id,omitempty"`
	RemindType     string  `json:"remind_type,omitempty"` // https://bot.q.qq.com/wiki/develop/api/openapi/schedule/model.html#remindtype
}

// GetChannelSchedules 获取channel_id指定的子频道中当天的日程列表
//
// https://bot.q.qq.com/wiki/develop/api/openapi/schedule/get_schedules.html
func (bot *Bot) GetChannelSchedules(id string, since uint64) (schedules []Schedule, err error) {
	if since == 0 {
		err = bot.GetOpenAPI("/channels/"+id+"/schedules", "", &schedules)
	} else {
		err = bot.GetOpenAPIWithBody("/channels/"+id+"/schedules", "", &schedules, WriteBodyFromJSON(&struct {
			S uint64 `json:"since"`
		}{since}))
	}
	return
}

// GetScheduleInChannel 获取日程子频道 channel_id 下 schedule_id 指定的的日程的详情
//
// https://bot.q.qq.com/wiki/develop/api/openapi/schedule/get_schedule.html
func (bot *Bot) GetScheduleInChannel(channelid string, scheduleid string) (*Schedule, error) {
	return bot.getOpenAPIofSchedule("/channels/" + channelid + "/schedules/" + scheduleid)
}

// CreateScheduleInChannel 在 channel_id 指定的日程子频道下创建一个日程
//
// https://bot.q.qq.com/wiki/develop/api/openapi/schedule/post_schedule.html
//
// schedule 会被写入返回的对象
func (bot *Bot) CreateScheduleInChannel(id string, schedule *Schedule) error {
	return bot.PostOpenAPI("/channels/"+id+"/schedules", "", schedule, WriteBodyFromJSON(&struct {
		S *Schedule `json:"schedule"`
	}{schedule}))
}

// PatchScheduleInChannel 修改日程子频道 channel_id 下 schedule_id 指定的日程的详情
//
// https://bot.q.qq.com/wiki/develop/api/openapi/schedule/patch_schedule.html
//
// schedule 会被写入返回的对象
func (bot *Bot) PatchScheduleInChannel(channelid string, scheduleid string, schedule *Schedule) error {
	return bot.PatchOpenAPI("/channels/"+channelid+"/schedules/"+scheduleid, "", schedule, WriteBodyFromJSON(&struct {
		S *Schedule `json:"schedule"`
	}{schedule}))
}

// DeleteScheduleInChannel 删除日程子频道 channel_id 下 schedule_id 指定的日程
//
// https://bot.q.qq.com/wiki/develop/api/openapi/schedule/delete_schedule.html
func (bot *Bot) DeleteScheduleInChannel(channelid string, scheduleid string) error {
	return bot.DeleteOpenAPI("/channels/"+channelid+"/schedules/"+scheduleid, "", nil)
}
