package teastats

import (
	"github.com/TeaWeb/build/internal/tealogs"
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo"
)

func init() {
	// 注册筛选器
	RegisterFilter(
		new(RequestAllPeriodFilter),
		new(RequestPagePeriodFilter),
		new(RequestIPPeriodFilter),
		new(StatusAllPeriodFilter),
		new(StatusPagePeriodFilter),
		new(TrafficAllPeriodFilter),
		new(TrafficPagePeriodFilter),
		new(PVAllPeriodFilter),
		new(PVPagePeriodFilter),
		new(UVAllPeriodFilter),
		new(UVPagePeriodFilter),
		new(IPAllPeriodFilter),
		new(IPPagePeriodFilter),
		new(MethodAllPeriodFilter),
		new(MethodPagePeriodFilter),
		new(CostAllPeriodFilter),
		new(CostPagePeriodFilter),
		new(RefererDomainPeriodFilter),
		new(RefererURLPeriodFilter),
		new(LandingPagePeriodFilter),
		new(BackendAllPeriodFilter),
		new(LocationAllPeriodFilter),
		new(RewriteAllPeriodFilter),
		new(FastcgiAllPeriodFilter),
		new(DeviceAllPeriodFilter),
		new(OSAllPeriodFilter),
		new(BrowserAllPeriodFilter),
		new(RegionAllPeriodFilter),
		new(ProvinceAllPeriodFilter),
		new(CityAllPeriodFilter),
		new(WAFBlockAllPeriodFilter),
	)

	// 注册AccessLogHook
	tealogs.AddAccessLogHook(&tealogs.AccessLogHook{
		Process: func(accessLog *accesslogs.AccessLog) (goNext bool) {
			if !accessLog.ShouldStat() {
				return true
			}
			serverQueue := FindServerQueue(accessLog.ServerId)
			if serverQueue == nil {
				return true
			}
			serverQueue.Filter(accessLog)
			return true
		},
	})

	// kv storage
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		sharedKV = NewKVStorage("stat.leveldb")
	})

	// 停止
	TeaGo.BeforeStop(func(server *TeaGo.Server) {
		go func() {
			for _, serverQueue := range AllStartedServers {
				serverQueue.(*ServerQueue).Stop()
			}
		}()
	})
}
