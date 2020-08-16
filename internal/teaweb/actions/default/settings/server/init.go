package server

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/settings"
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantAll,
			}).
			Helper(new(settings.Helper)).
			Prefix("/settings/server").
			Get("/http", new(HttpAction)).
			Post("/httpUpdate", new(HttpUpdateAction)).
			Get("/https", new(HttpsAction)).
			Post("/httpsUpdate", new(HttpsUpdateAction)).
			GetPost("/security", new(SecurityAction)).
			EndAll()
	})
}
