package agentutils

import (
	"sync"
	"time"
)

var agentQueueMap = map[string]*AgentQueue{} // agentId => queue
var agentQueueLocker = &sync.Mutex{}

type AgentQueue struct {
	c chan *AgentEvent
}

func NewAgentQueue() *AgentQueue {
	return &AgentQueue{
		c: make(chan *AgentEvent, 32),
	}
}

// 通知事件
func PostAgentEvent(agentId string, event *AgentEvent) {
	state := FindAgentState(agentId)
	if !state.IsActive {
		return
	}

	agentQueueLocker.Lock()
	queue, ok := agentQueueMap[agentId]
	if !ok {
		agentQueueLocker.Unlock()
		return
	}
	agentQueueLocker.Unlock()
	select {
	case queue.c <- event:
	default:
	}
}

// 等待事件
func Wait(agentId string) *AgentEvent {
	agentQueueLocker.Lock()
	queue, ok := agentQueueMap[agentId]
	if !ok {
		queue = NewAgentQueue()
		agentQueueMap[agentId] = queue
	}
	agentQueueLocker.Unlock()

	timer := time.NewTimer(59 * time.Second)

	select {
	case event := <-queue.c:
		timer.Stop()
		return event
	case <-timer.C:
	}

	return nil
}
