package monitor

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/monitor/agent"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Prefix("/api/monitor").
			GetPost("", new(IndexAction)).
			GetPost("/agent/:agentId/item/:itemId", new(agent.ItemAction)).
			EndAll()
	})
}
