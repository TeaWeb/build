package groups

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents"
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
			Helper(new(agents.Helper)).
			Prefix("/agents/groups").
			Get("", new(IndexAction)).
			GetPost("/add", new(AddAction)).
			Post("/delete", new(DeleteAction)).
			Get("/detail", new(DetailAction)).
			GetPost("/update", new(UpdateAction)).
			Post("/move", new(MoveAction)).
			Get("/noticeReceivers", new(NoticeReceiversAction)).
			GetPost("/addNoticeReceivers", new(AddNoticeReceiversAction)).
			Post("/deleteNoticeReceivers", new(DeleteNoticeReceiversAction)).
			Post("/generateKey", new(GenerateKeyAction)).
			EndAll()
	})
}
