package api

import (
	"github.com/iwind/TeaGo/timers"
	"sync"
	"time"
)

var SharedApiWatching = NewApiWatching()

// API监控管理
type ApiWatching struct {
	apis   map[string]time.Time // path => time
	locker sync.Mutex
}

func NewApiWatching() *ApiWatching {
	w := &ApiWatching{
		apis: map[string]time.Time{},
	}

	timers.Loop(5*time.Second, func(looper *timers.Looper) {
		w.gc()
	})

	return w
}

func (this *ApiWatching) gc() {
	this.locker.Lock()
	defer this.locker.Unlock()

	for path, t := range this.apis {
		if time.Since(t).Seconds() >= 5 {
			delete(this.apis, path)
		}
	}
}

func (this *ApiWatching) Add(path string) {
	this.locker.Lock()
	defer this.locker.Unlock()

	this.apis[path] = time.Now()
}

func (this *ApiWatching) Remove(path string) {
	this.locker.Lock()
	defer this.locker.Unlock()

	delete(this.apis, path)
}

func (this *ApiWatching) Contains(path string) bool {
	this.locker.Lock()
	defer this.locker.Unlock()

	_, ok := this.apis[path]
	return ok
}
