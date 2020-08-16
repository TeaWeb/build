package login

import (
	"github.com/iwind/TeaGo/actions"
	"sync"
)

var observers []func(action actions.ActionWrapper) bool
var observersLocker sync.Mutex

func Observe(f func(action actions.ActionWrapper) bool) {
	observersLocker.Lock()
	defer observersLocker.Unlock()

	observers = append(observers, f)
}

func Notify(action actions.ActionWrapper) bool {
	observersLocker.Lock()
	defer observersLocker.Unlock()

	for _, ob := range observers {
		b := ob(action)
		if !b {
			return false
		}
	}

	return true
}
