package shared

import (
	"github.com/TeaWeb/build/internal/teaevents"
	"sync"
)

var Locker = new(FileLocker)

// global file modify locker
type FileLocker struct {
	locker sync.Mutex
}

// lock
func (this *FileLocker) Lock() {
	this.locker.Lock()
}

// unlock for read
func (this *FileLocker) ReadUnlock() {
	this.locker.Unlock()
}

// unlock for write and notify
func (this *FileLocker) WriteUnlockNotify() {
	this.locker.Unlock()
	teaevents.Post(teaevents.NewConfigChangedEvent())
}

// unlock for write
func (this *FileLocker) WriteUnlock() {
	this.locker.Unlock()
}
