package api

import (
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/iwind/TeaGo/utils/time"
	"sync"
	"testing"
	"time"
)

func TestAPILimit(t *testing.T) {
	if teatesting.IsGlobal() {
		return
	}

	limit := NewAPILimit()
	limit.Concurrent = 3
	limit.Validate()

	wg := sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(i int) {
			limit.Begin()
			t.Log("process", i, timeutil.Format("H:i:s"))
			time.Sleep(1 * time.Second)
			limit.Done()

			wg.Done()
		}(i)
	}

	wg.Wait()
	t.Log("done")
}
