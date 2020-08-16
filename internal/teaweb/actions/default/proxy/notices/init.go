package notices

import (
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		// 注册路由
		server.
			Prefix("/proxy/notices").
			Helper(new(helpers.UserMustAuth)).
			Helper(new(Helper)).
			Get("", new(IndexAction)).
			GetPost("/addNoticeReceiver", new(AddNoticeReceiverAction)).
			Post("/deleteNoticeReceiver", new(DeleteNoticeReceiverAction)).
			EndAll()
	})
}
