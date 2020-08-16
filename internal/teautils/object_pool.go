package teautils

import (
	"time"
)

// 对象池
type ObjectPool struct {
	c      chan interface{}
	ticker *Ticker

	lastSize int
	newFunc  func() interface{}
}

// 创建新对象
func NewObjectPool(maxSize int, newFunc func() interface{}) *ObjectPool {
	if maxSize <= 0 {
		maxSize = 1024
	}
	pool := &ObjectPool{
		c:       make(chan interface{}, maxSize),
		newFunc: newFunc,
	}
	pool.start()
	return pool
}

func (this *ObjectPool) start() {
	// 清除Timer
	this.ticker = NewTicker(1 * time.Minute)
	go func() {
		for this.ticker.Next() {
			currentSize := len(this.c)
			if currentSize <= 32 || this.lastSize == 0 || this.lastSize != currentSize {
				this.lastSize = currentSize
				continue
			}

			i := 0

		For:
			for {
				select {
				case _ = <-this.c:
					i++
					if i >= currentSize/2 {
						break For
					}
				default:
					break For
				}
			}
		}
	}()
}

// 获取一个新的对象
func (this *ObjectPool) Get() (obj interface{}) {
	select {
	case obj = <-this.c:
	default:
		return this.newFunc()
	}
	return
}

// 放回一个使用过的对象
func (this *ObjectPool) Put(obj interface{}) {
	select {
	case this.c <- obj:
	default:
		// 已达最大容量，则抛弃
	}
}

// 当前的数量
func (this *ObjectPool) Size() int {
	return len(this.c)
}

// 销毁
func (this *ObjectPool) Destroy() {
	this.ticker.Stop()
}
