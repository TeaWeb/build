package teaproxy

import (
	"github.com/TeaWeb/build/internal/teaevents"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/logs"
	"sync/atomic"
	"time"
)

// 当前QPS
var qps = int32(0)

// 对外可读取的QPS
var QPS = int32(0)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		// 计算QPS
		teautils.Every(1*time.Second, func(ticker *teautils.Ticker) {
			QPS = qps
			atomic.StoreInt32(&qps, 0)
		})
	})

	teaevents.On(teaevents.EventTypeReload, func(event teaevents.EventInterface) {
		// 重启服务
		err := SharedManager.Restart()
		if err != nil {
			logs.Error(err)
		}
	})
}
