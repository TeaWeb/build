package servers

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy"
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Prefix("/proxy/servers").
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantProxy,
			}).
			Helper(new(proxy.Helper)).
			Get("/waf", new(WafAction)).
			Post("/waf/update", new(WafUpdateAction)).
			EndAll()
	})
}
