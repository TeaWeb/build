package teautils

import (
	"github.com/iwind/TeaGo/logs"
	"sync"
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	ticker := NewTicker(3 * time.Second)
	go func() {
		time.Sleep(10 * time.Second)
		ticker.Stop()
	}()
	for ticker.Next() {
		logs.Println("tick")
	}
	t.Log("finished")
}

func TestTicker2(t *testing.T) {
	ticker := NewTicker(1 * time.Second)
	go func() {
		time.Sleep(5 * time.Second)
		ticker.Stop()
	}()
	for {
		logs.Println("loop")
		select {
		case <-ticker.C:
			logs.Println("tick")
		case <-ticker.S:
			return
		}
	}
}

func TestTickerEvery(t *testing.T) {
	i := 0
	wg := &sync.WaitGroup{}
	wg.Add(1)
	Every(2*time.Second, func(ticker *Ticker) {
		i++
		logs.Println("TestTickerEvery i:", i)
		if i >= 4 {
			ticker.Stop()
			wg.Done()
		}
	})
	wg.Wait()
}
