package agentutils

// Agent事件
type AgentEvent struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

// 获取新对象
func NewAgentEvent(name string, data interface{}) *AgentEvent {
	return &AgentEvent{
		Name: name,
		Data: data,
	}
}
