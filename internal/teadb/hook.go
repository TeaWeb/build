package teadb

import "sync"

var hooks = []*Hook{}
var hooksLocker = sync.Mutex{}

func AddHook(hook *Hook) {
	hooksLocker.Lock()
	defer hooksLocker.Unlock()

	hooks = append(hooks, hook)
}

type Hook struct {
	Setup func()
}

func callHookSetup() {
	for _, hook := range hooks {
		go hook.Setup()
	}
}
