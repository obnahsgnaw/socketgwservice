package internal

import (
	"github.com/obnahsgnaw/socketgwservice/application/register"
	_ "github.com/obnahsgnaw/socketgwservice/internal/backend"
	"github.com/obnahsgnaw/socketgwservice/internal/dal/field"
	_ "github.com/obnahsgnaw/socketgwservice/internal/frontend"
	"gorm.io/gorm"
)

func init() {
	// 注册数据库连接后初始化
	register.Register(func(p *register.Provider) {
		p.Database().RegisterMigrateModel(
			field.Models()...,
		)
		p.Database().RegisterConnectInitializer(func(db *gorm.DB) {
			//
		})
	})
}
