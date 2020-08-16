package agent

import "github.com/iwind/TeaGo"

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Prefix("/api/agent").
			Helper(new(Helper)).
			Get("", new(IndexAction)).
			Get("/pull", new(PullAction)).
			Put("/push", new(PushAction)). // 兼容老版本Agent
			Post("/push", new(PushAction)).
			Get("/upgrade", new(UpgradeAction)).
			EndAll().
			// 不需要认证的API
			Get("/api/agent/ip", new(IpAction)).
			EndAll()
	})
}
