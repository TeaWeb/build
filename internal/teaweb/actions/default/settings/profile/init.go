package profile

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/settings"
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.Static("/avatar", Tea.ConfigFile("avatars/"))

		server.
			Helper(new(helpers.UserMustAuth)).
			Helper(new(settings.Helper)).
			Prefix("/settings/profile").
			Get("", new(IndexAction)).
			GetPost("/update", new(UpdateAction)).
			GetPost("/updateAvatar", new(UpdateAvatarAction)).
			EndAll()
	})
}
