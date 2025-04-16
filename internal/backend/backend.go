package backend

import (
	"github.com/obnahsgnaw/socketgwservice/application/register"
	_ "github.com/obnahsgnaw/socketgwservice/internal/backend/controller"
	_ "github.com/obnahsgnaw/socketgwservice/internal/backend/controller/upload"
	_ "github.com/obnahsgnaw/socketgwservice/internal/backend/event"
)

// 注册数据初始处理器 RegisterHandle()
func init() {
	register.Register(func(p *register.Provider) {
		p.Database().RegisterDataInitializer(
		//
		)
	})
}
