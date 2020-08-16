package cluster

import (
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Prefix("/settings/cluster").
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantAll,
			}).
			Helper(new(Helper)).
			Get("", new(IndexAction)).
			GetPost("/update", new(UpdateAction)).
			Post("/connect", new(ConnectAction)).
			Post("/sync", new(SyncAction)).
			EndAll()

	})
}
