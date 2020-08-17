package main

import (
	"github.com/TeaWeb/build/internal/teaagent/agentconst"
	teaagents "github.com/TeaWeb/build/internal/teaagent/agents"
	"github.com/TeaWeb/plugin/pkg/loader"
	"github.com/TeaWeb/plugin/pkg/plugins"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/timers"
	"time"
)

func main() {
	logs.Println("load plugin")
	p := plugins.NewPlugin()
	p.Name = "Agent"
	p.Code = "agent.teaweb"
	p.Developer = "TeaWeb"
	p.Version = "v" + agentconst.AgentVersion
	p.Date = "2019-11-01"
	p.Site = "https://github.com/TeaWeb/agent"
	p.Description = "本地Agent插件"
	p.OnStart(func() {
		timers.Delay(2*time.Second, func(timer *time.Timer) {
			teaagents.Start()
		})
	})
	loader.Start(p)
}
