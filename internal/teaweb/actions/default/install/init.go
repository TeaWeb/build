package install

import (
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Prefix("/install").
			Helper(new(helpers.UserMustAuth)).
			GetPost("", new(IndexAction)).
			Post("/skip", new(SkipAction)).
			Post("/test", new(TestAction)).
			Post("/save", new(SaveAction)).
			EndAll()
	})
}
