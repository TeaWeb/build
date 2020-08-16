package teadb

import (
	"github.com/TeaWeb/build/internal/teaevents"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
)

func init() {
	// 在测试环境下直接建立数据库，在二进制环境下需要等服务启动的时候才启动
	if Tea.IsTesting() {
		SetupDB()
	} else {
		TeaGo.BeforeStart(func(server *TeaGo.Server) {
			SetupDB()
		})
	}

	// 重启服务
	teaevents.On(teaevents.EventTypeReload, func(event teaevents.EventInterface) {
		db := SharedDB()
		if db != nil {
			err := db.Shutdown()
			if err != nil {
				logs.Println("[db]restart error:", err.Error())
			}
			ChangeDB()
		}
	})
}
