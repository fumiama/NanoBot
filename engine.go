package nano

//go:generate go run codegen/engine/main.go

// 生成空引擎
func newEngine() *Engine {
	return &Engine{
		preHandler:  []Rule{},
		midHandler:  []Rule{},
		postHandler: []Process{},
	}
}

var defaultEngine = newEngine()

// Engine is the pre_handler, mid_handler, post_handler manager
type Engine struct {
	preHandler  []Rule
	midHandler  []Rule
	postHandler []Process
	matchers    []*Matcher
	prio        int
	service     string
	datafolder  string
}

// Delete 移除该 Engine 注册的所有 Matchers
func (e *Engine) Delete() {
	for _, m := range e.matchers {
		m.Delete()
	}
}

// UsePreHandler 向该 Engine 添加新 PreHandler(Rule),
// 会在 Rule 判断前触发，如果 preHandler
// 没有通过，则 Rule, Matcher 不会触发
//
// 可用于分群组管理插件等
func (e *Engine) UsePreHandler(rules ...Rule) {
	e.preHandler = append(e.preHandler, rules...)
}

// UseMidHandler 向该 Engine 添加新 MidHandler(Rule),
// 会在 Rule 判断后， Matcher 触发前触发，如果 midHandler
// 没有通过，则 Matcher 不会触发
//
// 可用于速率限制等
func (e *Engine) UseMidHandler(rules ...Rule) {
	e.midHandler = append(e.midHandler, rules...)
}

// UsePostHandler 向该 Engine 添加新 PostHandler(Rule),
// 会在 Matcher 触发后触发，如果 PostHandler 返回 false,
// 则后续的 post handler 不会触发
//
// 可用于速率限制等
func (e *Engine) UsePostHandler(handler ...Process) {
	e.postHandler = append(e.postHandler, handler...)
}

// ApplySingle 应用反并发
func (e *Engine) ApplySingle(s *Single[int64]) *Engine {
	s.Apply(e)
	return e
}

// DataFolder 本插件数据目录, 默认 data/rbp/
func (e *Engine) DataFolder() string {
	return e.datafolder
}

// On 添加新的指定消息类型的匹配器(默认Engine)
func On(typ string, rules ...Rule) *Matcher { return defaultEngine.On(typ, rules...) }

// On 添加新的指定消息类型的匹配器
func (e *Engine) On(typ string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Type:   typ,
		Rules:  rules,
		Engine: e,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}
