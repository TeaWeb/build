package notices

import (
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(&helpers.UserMustAuth{}).
			Helper(new(Helper)).
			Prefix("/notices").
			Get("", new(IndexAction)).
			Get("/badge", new(BadgeAction)).
			Post("/setRead", new(SetReadAction)).
			Get("/setting", new(SettingAction)).
			Get("/level", new(LevelAction)).
			GetPost("/addReceiver", new(AddReceiverAction)).
			Get("/medias", new(MediasAction)).
			GetPost("/addMedia", new(AddMediaAction)).
			Post("/deleteMedia", new(DeleteMediaAction)).
			Get("/mediaDetail", new(MediaDetailAction)).
			GetPost("/updateMedia", new(UpdateMediaAction)).
			GetPost("/testMedia", new(TestMediaAction)).
			Post("/deleteReceiver", new(DeleteReceiverAction)).
			Get("/receiver", new(ReceiverAction)).
			GetPost("/updateReceiver", new(UpdateReceiverAction)).
			Post("/sound", new(SoundAction)).
			EndAll()
	})
}
