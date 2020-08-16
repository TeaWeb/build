package teamemory

import (
	"bytes"
	"compress/gzip"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/timers"
	"math"
	"time"
)

// Memory Cache Grid
//
// |           Grid                 |
// | cell1, cell2, ..., cell1024    |
// | item1, item2, ..., item1000000 |
type Grid struct {
	cells      []*Cell
	countCells uint64

	recycleIndex    int
	recycleLooper   *timers.Looper
	recycleInterval int

	gzipLevel int

	limitSize  int64
	limitCount int
}

func NewGrid(countCells int, opt ...interface{}) *Grid {
	grid := &Grid{
		recycleIndex: -1,
	}

	for _, o := range opt {
		switch x := o.(type) {
		case *CompressOpt:
			grid.gzipLevel = x.Level
		case *LimitSizeOpt:
			grid.limitSize = x.Size
		case *LimitCountOpt:
			grid.limitCount = x.Count
		case *RecycleIntervalOpt:
			grid.recycleInterval = x.Interval
		}
	}

	cells := []*Cell{}
	if countCells <= 0 {
		countCells = 1
	} else if countCells > 100*10000 {
		countCells = 100 * 10000
	}
	for i := 0; i < countCells; i++ {
		cell := NewCell()
		cell.LimitSize = int64(math.Floor(float64(grid.limitSize) / float64(countCells)))
		cell.LimitCount = int(math.Floor(float64(grid.limitCount)) / float64(countCells))

		cells = append(cells, cell)
	}
	grid.cells = cells
	grid.countCells = uint64(len(cells))

	grid.recycleTimer()
	return grid
}

// get all cells in the grid
func (this *Grid) Cells() []*Cell {
	return this.cells
}

func (this *Grid) WriteItem(item *Item) {
	if this.countCells <= 0 {
		return
	}
	hashKey := item.HashKey()
	this.cellForHashKey(hashKey).Write(hashKey, item)
}

func (this *Grid) WriteInt64(key []byte, value int64, lifeSeconds int64) {
	this.WriteItem(&Item{
		Key:        key,
		Type:       ItemInt64,
		ValueInt64: value,
		ExpireAt:   time.Now().Unix() + lifeSeconds,
	})
}

func (this *Grid) IncreaseInt64(key []byte, delta int64, lifeSeconds int64) (result int64) {
	hashKey := HashKey(key)
	return this.cellForHashKey(hashKey).Increase64(key, time.Now().Unix()+lifeSeconds, hashKey, delta)
}

func (this *Grid) WriteString(key []byte, value string, lifeSeconds int64) {
	this.WriteBytes(key, []byte(value), lifeSeconds)
}

func (this *Grid) WriteBytes(key []byte, value []byte, lifeSeconds int64) {
	isCompressed := false
	if this.gzipLevel != gzip.NoCompression {
		buf := bytes.NewBuffer([]byte{})
		writer, err := gzip.NewWriterLevel(buf, this.gzipLevel)
		if err != nil {
			logs.Error(err)
			this.WriteItem(&Item{
				Key:        key,
				Type:       ItemBytes,
				ValueBytes: value,
				ExpireAt:   time.Now().Unix() + lifeSeconds,
			})
			return
		}

		_, err = writer.Write([]byte(value))
		if err != nil {
			logs.Error(err)
			this.WriteItem(&Item{
				Key:        key,
				Type:       ItemBytes,
				ValueBytes: value,
				ExpireAt:   time.Now().Unix() + lifeSeconds,
			})
			return
		}

		err = writer.Close()
		if err != nil {
			logs.Error(err)
			this.WriteItem(&Item{
				Key:        key,
				Type:       ItemBytes,
				ValueBytes: value,
				ExpireAt:   time.Now().Unix() + lifeSeconds,
			})
			return
		}
		value = buf.Bytes()
		isCompressed = true
	}

	this.WriteItem(&Item{
		Key:          key,
		Type:         ItemBytes,
		ValueBytes:   value,
		ExpireAt:     time.Now().Unix() + lifeSeconds,
		IsCompressed: isCompressed,
	})
}

func (this *Grid) WriteInterface(key []byte, value interface{}, lifeSeconds int64) {
	this.WriteItem(&Item{
		Key:            key,
		Type:           ItemInterface,
		ValueInterface: value,
		ExpireAt:       time.Now().Unix() + lifeSeconds,
		IsCompressed:   false,
	})
}

func (this *Grid) Read(key []byte) *Item {
	if this.countCells <= 0 {
		return nil
	}
	hashKey := HashKey(key)
	return this.cellForHashKey(hashKey).Read(hashKey)
}

func (this *Grid) Stat() *Stat {
	stat := &Stat{}
	for _, cell := range this.cells {
		cellStat := cell.Stat()
		stat.CountItems += cellStat.CountItems
		stat.TotalBytes += cellStat.TotalBytes
	}
	return stat
}

func (this *Grid) Delete(key []byte) {
	if this.countCells <= 0 {
		return
	}
	hashKey := HashKey(key)
	this.cellForHashKey(hashKey).Delete(hashKey)
}

func (this *Grid) Reset() {
	for _, cell := range this.cells {
		cell.Reset()
	}
}

func (this *Grid) Destroy() {
	if this.recycleLooper != nil {
		this.recycleLooper.Stop()
		this.recycleLooper = nil
	}
	this.cells = nil
}

func (this *Grid) cellForHashKey(hashKey uint64) *Cell {
	if hashKey < 0 {
		return this.cells[-hashKey%this.countCells]
	} else {
		return this.cells[hashKey%this.countCells]
	}
}

func (this *Grid) recycleTimer() {
	duration := 1 * time.Minute
	if this.recycleInterval > 0 {
		duration = time.Duration(this.recycleInterval) * time.Second
	}
	this.recycleLooper = timers.Loop(duration, func(looper *timers.Looper) {
		if this.countCells == 0 {
			return
		}
		this.recycleIndex++
		if this.recycleIndex > int(this.countCells-1) {
			this.recycleIndex = 0
		}
		this.cells[this.recycleIndex].Recycle()
	})
}
