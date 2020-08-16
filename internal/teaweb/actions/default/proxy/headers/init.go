package headers

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy"
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantProxy,
			}).
			Helper(new(proxy.Helper)).
			Prefix("/proxy/headers").
			Get("", new(IndexAction)).
			Get("/data", new(DataAction)).
			GetPost("/add", new(AddAction)).
			Post("/delete", new(DeleteAction)).
			GetPost("/update", new(UpdateAction)).
			GetPost("/addIgnore", new(AddIgnoreAction)).
			GetPost("/updateIgnore", new(UpdateIgnoreAction)).
			Post("/deleteIgnore", new(DeleteIgnoreAction)).
			GetPost("/addRequestHeader", new(AddRequestHeaderAction)).
			GetPost("/updateRequestHeader", new(UpdateRequestHeaderAction)).
			Post("/deleteRequestHeader", new(DeleteRequestHeaderAction)).
			EndAll()
	})
}
