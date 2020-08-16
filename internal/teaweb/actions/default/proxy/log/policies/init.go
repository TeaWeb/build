package policies

import (
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Prefix("/proxy/log/policies").
			Helper(new(helpers.UserMustAuth)).
			Helper(new(Helper)).
			Get("", new(IndexAction)).
			GetPost("/add", new(AddAction)).
			Post("/delete", new(DeleteAction)).
			Get("/policy", new(PolicyAction)).
			GetPost("/update", new(UpdateAction)).
			EndAll()
	})
}
