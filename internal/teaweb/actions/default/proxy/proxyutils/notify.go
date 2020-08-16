package proxyutils

import "sync"

var proxyChanged = false

// 监控者
var observers []func()
var observerLocker = sync.Mutex{}

func NotifyChange() {
	proxyChanged = true

	// 执行监控者
	observerLocker.Lock()
	defer observerLocker.Unlock()
	for _, observer := range observers {
		go observer()
	}
}

func FinishChange() {
	proxyChanged = false
}

func ProxyIsChanged() bool {
	return proxyChanged
}

func Observe(f func()) {
	observerLocker.Lock()
	defer observerLocker.Unlock()
	observers = append(observers, f)
}
