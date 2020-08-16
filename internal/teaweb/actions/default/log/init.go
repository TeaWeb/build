package log

import (
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

// 初始化
func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		// 路由设置
		server.
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantLog,
			}).
			Helper(new(Helper)).
			Prefix("/log").
			GetPost("/runtime", new(RuntimeAction)).
			GetPost("/audit", new(AuditAction)).
			EndAll()
	})
}
