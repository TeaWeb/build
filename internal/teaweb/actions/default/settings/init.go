package settings

import (
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

var serverChanged = false

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantAll,
			}).
			Helper(new(Helper)).
			Prefix("/settings").
			Get("", new(IndexAction)).
			EndAll()
	})
}

func NotifyServerChange() {
	serverChanged = true
}

func ServerIsChanged() bool {
	return serverChanged
}
