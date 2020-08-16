package agentutils

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"sync"
)

var agentRuntimeMap = map[string]*agents.AgentConfig{} // 当前在运行的 agent id => agent
var agentRuntimeLocker = sync.Mutex{}

// 查找正在运行中的Agent，用来维护Agent的状态
func FindAgentRuntime(agentConfig *agents.AgentConfig) *agents.AgentConfig {
	if agentConfig == nil {
		return nil
	}
	agentRuntimeLocker.Lock()
	defer agentRuntimeLocker.Unlock()

	agent, found := agentRuntimeMap[agentConfig.Id]
	if found {
		return agent
	} else {
		agentRuntimeMap[agentConfig.Id] = agentConfig
	}
	return agentConfig
}
