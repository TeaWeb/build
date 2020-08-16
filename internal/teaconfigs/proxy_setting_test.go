package teaconfigs

import (
	"github.com/iwind/TeaGo/logs"
	"sync"
	"testing"
)

func TestSharedProxySetting(t *testing.T) {
	logs.PrintAsJSON(SharedProxySetting(), t)
}

func TestSharedProxySettingConcurrent(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			setting := SharedProxySetting()
			if setting == nil {
				t.Log("nil")
			}
		}()
	}
	wg.Wait()
}
