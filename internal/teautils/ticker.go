package teautils

import (
	"time"
)

// 类似于time.Ticker，但能够真正地停止
type Ticker struct {
	raw *time.Ticker

	S chan bool
	C <-chan time.Time

	isStopped bool
}

// 创建新Ticker
func NewTicker(duration time.Duration) *Ticker {
	raw := time.NewTicker(duration)
	return &Ticker{
		raw: raw,
		C:   raw.C,
		S:   make(chan bool, 1),
	}
}

// 查找下一个Tick
func (this *Ticker) Next() bool {
	select {
	case <-this.raw.C:
		return true
	case <-this.S:
		return false
	}
}

// 停止
func (this *Ticker) Stop() {
	if this.isStopped {
		return
	}

	this.isStopped = true

	this.raw.Stop()
	this.S <- true
}
