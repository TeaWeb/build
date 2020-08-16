package dashboard

import (
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(helpers.UserMustAuth)).
			Helper(new(Helper)).
			Prefix("/dashboard").
			GetPost("", new(IndexAction)).
			Get("/logs", new(LogsAction)).
			EndAll()
	})
}
