package teaproxy

import (
	"sync"
	"testing"
)

func BenchmarkClientLock(b *testing.B) {
	locker := sync.RWMutex{}
	for i := 0; i < b.N; i++ {
		locker.RLock()
		locker.RUnlock()
	}
}
