package teaagents

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/iwind/TeaGo/logs"
	"time"
)

// 加载Local配置
func loadLocalConfig() error {
	agent := agents.NewAgentConfigFromFile("agent.local.conf")
	if agent == nil {
		err := agents.LocalAgentConfig().Save()
		if err != nil {
			logs.Println("[agent]" + err.Error())
		} else {
			return loadLocalConfig()
		}

		time.Sleep(30 * time.Second)
		return loadLocalConfig()
	}
	err := agent.Validate()
	if err != nil {
		logs.Println("[agent]" + err.Error())
		time.Sleep(30 * time.Second)
		return loadLocalConfig()
	}
	runningAgent = agent
	connectConfig.Key = agent.Key

	if !isBooting {
		// 定时任务
		scheduleTasks()

		// 监控项数据
		scheduleItems()
	}
	return nil
}
