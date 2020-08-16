package login

import (
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		// 自定义登录URL
		prefix := "/login"
		security := configs.SharedAdminConfig().Security
		if security != nil {
			prefix = security.NewLoginURL()
		}

		server.
			Prefix(prefix).
			GetPost("", new(IndexAction)).
			EndAll()
	})
}
