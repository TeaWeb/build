package settings

import (
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantAgent,
			}).
			Prefix("/agents/settings").
			Helper(new(Helper)).
			Get("", new(IndexAction)).
			GetPost("/update", new(UpdateAction)).
			Get("/install", new(InstallAction)).
			Get("/noticeReceivers", new(NoticeReceiversAction)).
			GetPost("/addNoticeReceivers", new(AddNoticeReceiversAction)).
			Post("/deleteNoticeReceivers", new(DeleteNoticeReceiversAction)).
			Post("/on", new(OnAction)).
			Post("/off", new(OffAction)).
			EndAll()
	})
}
