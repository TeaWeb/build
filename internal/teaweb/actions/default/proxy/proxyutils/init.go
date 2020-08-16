package proxyutils

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaevents"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		// 启动统计
		go func() {
			for _, server := range teaconfigs.LoadServerConfigsFromDir(Tea.ConfigDir()) {
				if !server.On {
					continue
				}
				if server.StatBoard == nil && server.RealtimeBoard == nil {
					continue
				}
				ReloadServerStats(server.Id)
			}
		}()

		// 处理事件
		teaevents.On(teaconfigs.EventBackendDown, func(event teaevents.EventInterface) {
			realEvent, ok := event.(*teaconfigs.BackendDownEvent)
			if !ok {
				return
			}
			if realEvent.Server == nil || !realEvent.Server.On {
				return
			}
			err := NotifyProxyBackendDownMessage(realEvent)
			if err != nil {
				logs.Error(err)
			}
		})
		teaevents.On(teaconfigs.EventBackendUp, func(event teaevents.EventInterface) {
			realEvent, ok := event.(*teaconfigs.BackendUpEvent)
			if !ok {
				return
			}
			if realEvent.Server == nil || !realEvent.Server.On {
				return
			}
			err := NotifyProxyBackendUpMessage(realEvent)
			if err != nil {
				logs.Error(err)
			}
		})
	})
}
