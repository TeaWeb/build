package agentutils

import "sync"

var agentStateMap = map[string]*AgentState{} //agentId => state
var agentStateLocker = sync.Mutex{}

// Agent状态
type AgentState struct {
	Version  string  // 版本号
	OsName   string  // 操作系统
	Speed    float64 // 连接速度，ms
	IP       string  // IP地址
	IsActive bool    // 是否在线
}

// 查找Agent状态
func FindAgentState(agentId string) *AgentState {
	agentStateLocker.Lock()
	defer agentStateLocker.Unlock()

	state, ok := agentStateMap[agentId]
	if ok {
		return state
	}
	state = &AgentState{}
	agentStateMap[agentId] = state
	return state
}
