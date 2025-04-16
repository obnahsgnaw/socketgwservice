package frontend

import (
	"github.com/obnahsgnaw/socketgwservice/application/register"
	_ "github.com/obnahsgnaw/socketgwservice/internal/frontend/controller"
	_ "github.com/obnahsgnaw/socketgwservice/internal/frontend/controller/upload"
	_ "github.com/obnahsgnaw/socketgwservice/internal/frontend/event"
	_ "github.com/obnahsgnaw/socketgwservice/internal/frontend/sockethandler/tcp"
	_ "github.com/obnahsgnaw/socketgwservice/internal/frontend/sockethandler/wss"
)

// 注册数据初始处理器 RegisterHandle()
func init() {
	register.Register(func(p *register.Provider) {
		p.Database().RegisterDataInitializer(
		//
		)
	})
}
