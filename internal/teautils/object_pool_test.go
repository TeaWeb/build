package teautils

import (
	"github.com/iwind/TeaGo/maps"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewObjectPool(t *testing.T) {
	pool := NewObjectPool(5, func() interface{} {
		return new(maps.Map)
	})
	obj := pool.Get()
	if obj == nil {
		t.Fatal("obj should not be nil")
	}
	pool.Put(obj)
	obj2 := pool.Get()
	if obj != obj2 {
		t.Fatal("obj != obj2")
	}
}

func TestObjectPool_Get(t *testing.T) {
	var count = int64(0)
	pool := NewObjectPool(10240, func() interface{} {
		atomic.AddInt64(&count, 1)
		return new(maps.Map)
	})

	for j := 0; j < 10; j++ {
		concurrent := 10240
		wg := sync.WaitGroup{}
		wg.Add(concurrent)
		for i := 0; i < concurrent; i++ {
			go func(i int) {
				obj := pool.Get()
				_ = obj

				if i%2 == 0 {
					time.Sleep(100 * time.Millisecond)
				} else if i%3 == 0 {
					time.Sleep(200 * time.Millisecond)
				}

				pool.Put(obj)
				wg.Done()
			}(i)
		}
		wg.Wait()
	}

	t.Log("new objects:", count)
}

func BenchmarkObjectPool_Get(b *testing.B) {
	pool := NewObjectPool(10240, func() interface{} {
		return new(maps.Map)
	})

	for i := 0; i < b.N; i++ {
		obj := pool.Get()
		_ = obj
		pool.Put(obj)
	}
}
