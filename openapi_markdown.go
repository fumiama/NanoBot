package nano

// MessageMarkdown https://bot.q.qq.com/wiki/develop/api/openapi/message/model.html#messagemarkdown
type MessageMarkdown struct {
	TemplateID       int                     `json:"template_id,omitempty"`
	CustomTemplateID string                  `json:"custom_template_id,omitempty"`
	Params           []MessageMarkdownParams `json:"params,omitempty"`
	Content          string                  `json:"content,omitempty"` // 原生 markdown 内容,与上面三个参数互斥,参数都传值将报错
}

// MessageMarkdownParams https://bot.q.qq.com/wiki/develop/api/openapi/message/model.html#messagemarkdownparams
type MessageMarkdownParams struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

// MessageKeyboard https://bot.q.qq.com/wiki/develop/api/openapi/message/model.html#messagekeyboard
type MessageKeyboard struct {
	ID      string          `json:"id,omitempty"`
	Content *InlineKeyboard `json:"content,omitempty"` // 自定义 keyboard 内容,与 id 参数互斥,参数都传值将报错
}

// InlineKeyboard 消息按钮对象
//
// https://bot.q.qq.com/wiki/develop/api/openapi/message/message_keyboard.html
type InlineKeyboard struct {
	Rows     []InlineKeyboardRow `json:"rows"`
	BotAppID int                 `json:"bot_appid"`
}

// InlineKeyboardRow https://bot.q.qq.com/wiki/develop/api/openapi/message/message_keyboard.html#inlinekeyboardrow
type InlineKeyboardRow struct {
	Buttons []InlineKeyboardButton `json:"buttons"`
}

// InlineKeyboardButton https://bot.q.qq.com/wiki/develop/api/openapi/message/message_keyboard.html#button
type InlineKeyboardButton struct {
	ID         string                         `json:"id"`
	RenderData InlineKeyboardButtonRenderData `json:"render_data"`
	Action     InlineKeyboardButtonAction     `json:"action"`
}

// InlineKeyboardButtonRenderDataStyle https://bot.q.qq.com/wiki/develop/api/openapi/message/message_keyboard.html#renderstyle
type InlineKeyboardButtonRenderDataStyle int

const (
	InlineKeyboardButtonRenderDataStyleGray InlineKeyboardButtonRenderDataStyle = iota
	InlineKeyboardButtonRenderDataStyleBlue
)

// InlineKeyboardButtonRenderData https://bot.q.qq.com/wiki/develop/api/openapi/message/message_keyboard.html#renderdata
type InlineKeyboardButtonRenderData struct {
	Label        string                              `json:"label"`
	VisitedLabel string                              `json:"visited_label"`
	Style        InlineKeyboardButtonRenderDataStyle `json:"style"`
}

// InlineKeyboardButtonActionType https://bot.q.qq.com/wiki/develop/api/openapi/message/message_keyboard.html#actiontype
type InlineKeyboardButtonActionType int

const (
	InlineKeyboardButtonActionTypeHTTP InlineKeyboardButtonActionType = iota
	InlineKeyboardButtonActionTypeCallback
	InlineKeyboardButtonActionTypeAtBot
)

// InlineKeyboardButtonAction https://bot.q.qq.com/wiki/develop/api/openapi/message/message_keyboard.html#action
type InlineKeyboardButtonAction struct {
	Type                 InlineKeyboardButtonActionType       `json:"type"`
	Permission           InlineKeyboardButtonActionPermission `json:"permission"`
	ClickLimit           int                                  `json:"click_limit"`
	UnsupportTips        string                               `json:"unsupport_tips"`
	Data                 string                               `json:"data"`
	AtBotShowChannelList bool                                 `json:"at_bot_show_channel_list"`
}

// InlineKeyboardButtonActionPermissionType https://bot.q.qq.com/wiki/develop/api/openapi/message/message_keyboard.html#permissiontype
type InlineKeyboardButtonActionPermissionType int

const (
	InlineKeyboardButtonActionPermissionTypeShimeiUser InlineKeyboardButtonActionPermissionType = iota
	InlineKeyboardButtonActionPermissionTypeAdmin
	InlineKeyboardButtonActionPermissionTypeAll
	InlineKeyboardButtonActionPermissionTypeShimeiRole
)

// InlineKeyboardButtonActionPermission https://bot.q.qq.com/wiki/develop/api/openapi/message/message_keyboard.html#permission
type InlineKeyboardButtonActionPermission struct {
	Type           InlineKeyboardButtonActionPermissionType `json:"type"`
	SpecifyRoleIDs []string                                 `json:"specify_role_ids"`
	SpecifyUserIDs []string                                 `json:"specify_user_ids"`
}
