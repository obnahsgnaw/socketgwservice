package event

import (
	"github.com/obnahsgnaw/api/service/autheduser"
	"github.com/obnahsgnaw/application/service/event"
	"github.com/obnahsgnaw/pbhttp/core/project"
	"github.com/obnahsgnaw/pbhttp/core/provider/backend/rpc/opelog"
	"github.com/obnahsgnaw/socketgwservice/application/register"
	"github.com/obnahsgnaw/socketgwservice/config"
)

func init() {
	// 注册事件
	register.Register(func(p *register.Provider) {
		p.Backend().Event().Register(operate, func(e *event.Event) {
			// 断言事件数据
			d := e.Data[0].(*OperateData)
			if p.CommonConfig().ServiceEnabled(project.OpeLogService) {
				_ = opelog.AddOpeLog(p.BackendRps(), d.Operator, config.Project, d.Target, d.TargetId, d.Action, d.Content, d.DataBefore, d.DataAfter)
			}
		})
	})
}

type OperateData struct {
	Operator   autheduser.User // 操作人
	Target     string          // 操作对象 如表名
	TargetId   uint32          // 操作对象标识 如 主键ID
	Action     string          // 操作动作    如 create
	Content    string          // 操作描述内容 如 管理员xxx 创建了用户 xxx
	DataBefore interface{}     // 操作前数据
	DataAfter  interface{}     // 操作后数据
}

// NewOperateEvent new事件
func NewOperateEvent(data *OperateData) *event.Event {
	return register.Provide.Backend().Event().Build(operate, data)
}
