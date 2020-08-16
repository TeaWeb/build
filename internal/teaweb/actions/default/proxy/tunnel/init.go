package tunnel

import (
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Prefix("/proxy/tunnel").
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantProxy,
			}).
			Get("", new(IndexAction)).
			GetPost("/update", new(UpdateAction)).
			Post("/generateSecret", new(GenerateSecretAction)).
			EndAll()
	})
}
