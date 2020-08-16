package rewrite

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy"
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Prefix("/proxy/rewrite").
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantProxy,
			}).
			Helper(new(Helper)).
			Helper(new(proxy.Helper)).
			Get("/data", new(DataAction)).
			GetPost("/add", new(AddAction)).
			Post("/delete", new(DeleteAction)).
			GetPost("/update", new(UpdateAction)).
			Post("/test", new(TestAction)).
			EndAll()
	})
}
