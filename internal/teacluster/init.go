package teacluster

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaevents"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/timers"
	"time"
)

func init() {
	if !ClusterEnabled {
		return
	}

	// register actions
	RegisterActionType(
		new(SuccessAction),
		new(FailAction),
		new(RegisterAction),
		new(PushAction),
		new(PullAction),
		new(NotifyAction),
		new(SumAction),
		new(SyncAction),
		new(PingAction),
		new(RunAction),
	)

	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		// build
		SharedManager.BuildSum()

		// start manager
		go func() {
			ticker := teautils.NewTicker(60 * time.Second)
			for {
				err := SharedManager.Start()
				if err != nil {
					logs.Println("[cluster]" + err.Error())
				}

				// retry N seconds later
				select {
				case <-ticker.C:
					// every N seconds
				case <-SharedManager.RestartChan:
					// retry immediately
				}
			}
		}()

		// start ping
		timers.Loop(60*time.Second, func(looper *timers.Looper) {
			node := teaconfigs.SharedNodeConfig()
			if node != nil && node.On && SharedManager.IsActive() {
				err := SharedManager.Write(&PingAction{})
				if err != nil {
					logs.Println("[cluster]" + err.Error())
				}
			}
		})
	})

	TeaGo.BeforeStop(func(server *TeaGo.Server) {
		if SharedManager != nil {
			err := SharedManager.Stop()
			if err != nil {
				logs.Error(err)
			}
		}
	})

	teaevents.On(teaevents.EventTypeConfigChanged, func(event teaevents.EventInterface) {
		node := teaconfigs.SharedNodeConfig()
		if node != nil && node.On {
			SharedManager.SetIsChanged(true)
		}
	})
}
