package teaagents

import (
	"encoding/json"
	"time"
)

// 监控项事件
type ItemEvent struct {
	Event     string      `json:"event"`
	AgentId   string      `json:"agentId"`
	AppId     string      `json:"appId"`
	ItemId    string      `json:"itemId"`
	Value     interface{} `json:"value"`
	Error     string      `json:"error"`
	BeginAt   int64       `json:"beginAt"`
	Timestamp int64       `json:"timestamp"`
	CostMs    float64     `json:"costMs"`
}

// 获取新监控项事件
func NewItemEvent(agentId string, appId string, itemId string, value interface{}, err error, beginAt int64, costMs float64) *ItemEvent {
	errorString := ""
	if err != nil {
		errorString = err.Error()
	}
	return &ItemEvent{
		Event:     "ItemEvent",
		AgentId:   agentId,
		AppId:     appId,
		ItemId:    itemId,
		Value:     value,
		Error:     errorString,
		BeginAt:   beginAt,
		Timestamp: time.Now().Unix(),
		CostMs:    costMs,
	}
}

// 转换为JSON
func (this *ItemEvent) AsJSON() ([]byte, error) {
	return json.Marshal(this)
}
