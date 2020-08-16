package teamemory

import (
	"compress/gzip"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestMemoryGrid_Write(t *testing.T) {
	grid := NewGrid(5, NewRecycleIntervalOpt(2), NewLimitSizeOpt(10240))
	t.Log("123456:", grid.Read([]byte("123456")))

	grid.WriteInt64([]byte("abc"), 1, 5)
	t.Log(grid.Read([]byte("abc")).ValueInt64)

	grid.WriteString([]byte("abc"), "123", 5)
	t.Log(string(grid.Read([]byte("abc")).Bytes()))

	grid.WriteBytes([]byte("abc"), []byte("123"), 5)
	t.Log(grid.Read([]byte("abc")).Bytes())

	grid.Delete([]byte("abc"))
	t.Log(grid.Read([]byte("abc")))

	for i := 0; i < 100; i++ {
		grid.WriteInt64([]byte(fmt.Sprintf("%d", i)), 123, 1)
	}

	t.Log("before recycle:")
	for index, cell := range grid.cells {
		t.Log("cell:", index, len(cell.mapping), "items")
	}

	time.Sleep(3 * time.Second)
	t.Log("after recycle:")
	for index, cell := range grid.cells {
		t.Log("cell:", index, len(cell.mapping), "items")
	}

	grid.Destroy()
}

func TestMemoryGrid_Write_LimitCount(t *testing.T) {
	grid := NewGrid(2, NewLimitCountOpt(10))
	for i := 0; i < 100; i++ {
		grid.WriteInt64([]byte(strconv.Itoa(i)), int64(i), 30)
	}
	t.Log(grid.Stat().CountItems, "items")
}

func TestMemoryGrid_Compress(t *testing.T) {
	grid := NewGrid(5, NewCompressOpt(1))
	grid.WriteString([]byte("hello"), strings.Repeat("abcd", 10240), 30)
	t.Log(len(string(grid.Read([]byte("hello")).String())))
	t.Log(len(grid.Read([]byte("hello")).ValueBytes))
}

func BenchmarkMemoryGrid_Performance(b *testing.B) {
	grid := NewGrid(1024)
	for i := 0; i < b.N; i++ {
		grid.WriteInt64([]byte("key:"+strconv.Itoa(i)), int64(i), 3600)
	}
}

func TestMemoryGrid_Performance(t *testing.T) {
	runtime.GOMAXPROCS(1)

	grid := NewGrid(1024)

	now := time.Now()

	s := []byte(strings.Repeat("abcd", 10*1024))

	for i := 0; i < 100000; i++ {
		grid.WriteBytes([]byte(fmt.Sprintf("key:%d_%d", i, 1)), s, 3600)
		item := grid.Read([]byte(fmt.Sprintf("key:%d_%d", i, 1)))
		if item != nil {
			_ = item.String()
		}
	}

	countItems := 0
	for _, cell := range grid.cells {
		countItems += len(cell.mapping)
	}
	t.Log(countItems, "items")

	t.Log(time.Since(now).Seconds()*1000, "ms")
}

func TestMemoryGrid_Performance_Concurrent(t *testing.T) {
	//runtime.GOMAXPROCS(1)

	grid := NewGrid(1024)

	now := time.Now()

	s := []byte(strings.Repeat("abcd", 10*1024))

	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())

	for c := 0; c < runtime.NumCPU(); c++ {
		go func(c int) {
			defer wg.Done()
			for i := 0; i < 50000; i++ {
				grid.WriteBytes([]byte(fmt.Sprintf("key:%d_%d", i, c)), s, 3600)
				item := grid.Read([]byte(fmt.Sprintf("key:%d_%d", i, c)))
				if item != nil {
					_ = item.String()
				}
			}
		}(c)
	}

	wg.Wait()
	countItems := 0
	for _, cell := range grid.cells {
		countItems += len(cell.mapping)
	}
	t.Log(countItems, "items")

	t.Log(time.Since(now).Seconds()*1000, "ms")
}

func TestMemoryGrid_CompressPerformance(t *testing.T) {
	runtime.GOMAXPROCS(1)

	grid := NewGrid(1024, NewCompressOpt(gzip.BestCompression))

	now := time.Now()
	data := []byte(strings.Repeat("abcd", 1024))

	for i := 0; i < 100000; i++ {
		grid.WriteBytes([]byte(fmt.Sprintf("key:%d", i)), data, 3600)
		item := grid.Read([]byte(fmt.Sprintf("key:%d", i+100)))
		if item != nil {
			_ = item.String()
		}
	}

	countItems := 0
	for _, cell := range grid.cells {
		countItems += len(cell.mapping)
	}
	t.Log(countItems, "items")

	t.Log(time.Since(now).Seconds()*1000, "ms")
}

func TestMemoryGrid_IncreaseInt64(t *testing.T) {
	grid := NewGrid(1024)
	grid.WriteInt64([]byte("abc"), 123, 10)
	grid.IncreaseInt64([]byte("abc"), 123, 10)
	grid.IncreaseInt64([]byte("abc"), 123, 10)
	item := grid.Read([]byte("abc"))
	if item == nil {
		t.Fatal("item == nil")
	}

	if item.ValueInt64 != 369 {
		t.Fatal("not 369")
	}
}

func TestMemoryGrid_Destroy(t *testing.T) {
	grid := NewGrid(1024)
	grid.WriteInt64([]byte("abc"), 123, 10)
	t.Log(grid.recycleLooper, grid.cells)
	grid.Destroy()
	t.Log(grid.recycleLooper, grid.cells)

	if grid.recycleLooper != nil {
		t.Fatal("looper != nil")
	}
}

func TestMemoryGrid_Recycle(t *testing.T) {
	cell := NewCell()
	timestamp := time.Now().Unix()
	for i := 0; i < 300_0000; i++ {
		cell.Write(uint64(i), &Item{
			ExpireAt: timestamp - 30,
		})
	}
	before := time.Now()
	cell.Recycle()
	t.Log(time.Since(before).Seconds()*1000, "ms")
	t.Log(len(cell.mapping))

	runtime.GC()
	printMem(t)
}

func printMem(t *testing.T) {
	mem := &runtime.MemStats{}
	runtime.ReadMemStats(mem)
	t.Log(mem.Alloc/1024/1024, "M")
}
