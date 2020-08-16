package database

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
			Prefix("/settings/database").
			Get("", new(IndexAction)).
			Post("/tables", new(TablesAction)).
			Post("/tableStat", new(TableStatAction)).
			Post("/deleteTable", new(DeleteTableAction)).
			EndAll()
	})
}
