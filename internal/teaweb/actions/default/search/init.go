package search

import (
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Prefix("/search").
			Helper(new(helpers.UserMustAuth)).
			GetPost("", new(IndexAction)).
			EndAll()
	})
}
