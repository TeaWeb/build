package login

import (
	"github.com/iwind/TeaGo"
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/settings"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(helpers.UserMustAuth)).
			Helper(new(settings.Helper)).
			Prefix("/settings/login").
			Get("", new(IndexAction)).
			GetPost("/update", new(UpdateAction)).
			Post("/generateKey", new(GenerateKeyAction)).
			EndAll()
	})
}
