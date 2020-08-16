package teamemory

import (
	"math"
	"sync"
	"time"
)

type Cell struct {
	LimitSize  int64
	LimitCount int

	mapping    map[uint64]*Item // key => item
	list       *List            // { item1, item2, ... }
	totalBytes int64
	locker     sync.RWMutex
}

func NewCell() *Cell {
	return &Cell{
		mapping: map[uint64]*Item{},
		list:    NewList(),
	}
}

func (this *Cell) Write(hashKey uint64, item *Item) {
	if item == nil {
		return
	}
	this.locker.Lock()

	oldItem, ok := this.mapping[hashKey]
	if ok {
		this.list.Remove(oldItem)

		if this.LimitSize > 0 {
			this.totalBytes -= oldItem.Size()
		}
	}

	// limit count
	if this.LimitCount > 0 && len(this.mapping) >= this.LimitCount {
		this.locker.Unlock()
		return
	}

	// trim memory
	size := item.Size()
	shouldTrim := false
	if this.LimitSize > 0 && this.LimitSize < this.totalBytes+size {
		this.Trim()
		shouldTrim = true
	}

	// compare again
	if shouldTrim {
		if this.LimitSize > 0 && this.LimitSize < this.totalBytes+size {
			this.locker.Unlock()
			return
		}
	}

	this.totalBytes += size

	this.list.Add(item)
	this.mapping[hashKey] = item

	this.locker.Unlock()
}

func (this *Cell) Increase64(key []byte, expireAt int64, hashKey uint64, delta int64) (result int64) {
	this.locker.Lock()
	item, ok := this.mapping[hashKey]
	if ok {
		// reset to zero if expired
		if item.ExpireAt < time.Now().Unix() {
			item.ValueInt64 = 0
			item.ExpireAt = expireAt
		}
		item.IncreaseInt64(delta)
		result = item.ValueInt64
	} else {
		item := NewItem(key, ItemInt64)
		item.ValueInt64 = delta
		item.ExpireAt = expireAt
		this.mapping[hashKey] = item
		result = delta
	}
	this.locker.Unlock()
	return
}

func (this *Cell) Read(hashKey uint64) *Item {
	this.locker.Lock()

	item, ok := this.mapping[hashKey]
	if ok {
		this.list.Remove(item)
		this.list.Add(item)

		this.locker.Unlock()

		if item.ExpireAt < time.Now().Unix() {
			return nil
		}
		return item
	}

	this.locker.Unlock()
	return nil
}

func (this *Cell) Stat() *CellStat {
	this.locker.RLock()
	defer this.locker.RUnlock()

	return &CellStat{
		TotalBytes: this.totalBytes,
		CountItems: len(this.mapping),
	}
}

// trim NOT ACTIVE items
// should called in locker context
func (this *Cell) Trim() {
	l := len(this.mapping)
	if l == 0 {
		return
	}

	inactiveSize := int(math.Ceil(float64(l) / 10)) // trim 10% items
	this.list.Range(func(item *Item) (goNext bool) {
		inactiveSize--
		delete(this.mapping, item.HashKey())
		this.list.Remove(item)
		this.totalBytes -= item.Size()
		return inactiveSize > 0
	})
}

func (this *Cell) Delete(hashKey uint64) {
	this.locker.Lock()
	item, ok := this.mapping[hashKey]
	if ok {
		delete(this.mapping, hashKey)
		this.list.Remove(item)
		this.totalBytes -= item.Size()
	}
	this.locker.Unlock()
}

// range all items in the cell
func (this *Cell) Range(f func(item *Item)) {
	this.locker.Lock()
	for _, item := range this.mapping {
		f(item)
	}
	this.locker.Unlock()
}

func (this *Cell) Recycle() {
	this.locker.Lock()
	if len(this.mapping) == 0 {
		this.locker.Unlock()
		return
	}

	timestamp := time.Now().Unix()
	for key, item := range this.mapping {
		if item.ExpireAt < timestamp {
			delete(this.mapping, key)
			this.list.Remove(item)
			this.totalBytes -= item.Size()
		}
	}

	this.locker.Unlock()
}

func (this *Cell) Reset() {
	this.locker.Lock()
	this.list.Reset()
	this.mapping = map[uint64]*Item{}
	this.totalBytes = 0
	this.locker.Unlock()
}
