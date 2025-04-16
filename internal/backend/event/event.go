package event

import (
	"github.com/obnahsgnaw/application/service/event"
	"github.com/obnahsgnaw/socketgwservice/application/register"
)

func init() {
	register.Register(func(p *register.Provider) {
		// 注册事件
		p.Backend().Event().Register(demo, func(e *event.Event) {
			// 断言事件数据
			id := e.Data[0].(string)
			name := e.Data[1].(string)
			println(id, name)
		})
	})
}

// NewDemoEvent new事件
func NewDemoEvent(id, name string) *event.Event {
	return register.Provide.Backend().Event().Build(demo, id, name)
}

// 触发事件
// *event.Event.Fire()
