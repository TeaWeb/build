package agentutils

import (
	"sync"
)

// 监控项状态列表
var itemStateMap = map[string]*ItemState{} // itemId => ItemState
var itemStateLocker = sync.Mutex{}

// 添加监控项状态
func AddItemState(itemId string, itemState *ItemState) {
	itemStateLocker.Lock()
	itemStateMap[itemId] = itemState
	itemStateLocker.Unlock()
}

// 查找监控项状态
func FindItemState(itemId string) (state *ItemState, ok bool) {
	itemStateLocker.Lock()
	state, ok = itemStateMap[itemId]
	itemStateLocker.Unlock()
	return
}

// 删除监控项状态
func RemoveItemState(itemId string) {
	itemStateLocker.Lock()
	delete(itemStateMap, itemId)
	itemStateLocker.Unlock()
}
