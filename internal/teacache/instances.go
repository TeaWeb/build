package teacache

import (
	"sync"
)

var cachePolicyMap = map[string]ManagerInterface{}
var cachePolicyMapLocker = sync.RWMutex{}

// 重置管理器
func ResetCachePolicyManager(filename string) {
	cachePolicyMapLocker.Lock()
	defer cachePolicyMapLocker.Unlock()

	manager, ok := cachePolicyMap[filename]
	if ok {
		manager.Close()
		delete(cachePolicyMap, filename)
	}
}

// 获取管理器
func FindCachePolicyManager(filename string) ManagerInterface {
	cachePolicyMapLocker.Lock()
	defer cachePolicyMapLocker.Unlock()

	manager, ok := cachePolicyMap[filename]
	if ok {
		return manager
	}
	return nil
}
