package teamemory

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestCell_List(t *testing.T) {
	cell := NewCell()
	cell.Write(1, &Item{
		ValueInt64: 1,
	})
	cell.Write(2, &Item{
		ValueInt64: 2,
	})
	cell.Write(3, &Item{
		ValueInt64: 3,
	})

	{
		t.Log("====")
		l := cell.list
		for e := l.head; e != nil; e = e.Next {
			t.Log("element:", e.ValueInt64)
		}
	}

	cell.Write(1, &Item{
		ValueInt64: 1,
	})
	cell.Write(3, &Item{
		ValueInt64: 3,
	})
	cell.Write(3, &Item{
		ValueInt64: 3,
	})
	cell.Write(2, &Item{
		ValueInt64: 2,
	})
	cell.Delete(2)

	{
		t.Log("====")
		l := cell.list
		for e := l.head; e != nil; e = e.Next {
			t.Log("element:", e.ValueInt64)
		}
	}

	for _, m := range cell.mapping {
		t.Log(m.ValueInt64)
	}
}

func TestCell_LimitSize(t *testing.T) {
	cell := NewCell()
	cell.LimitSize = 1024

	for i := int64(0); i < 100; i ++ {
		key := []byte(fmt.Sprintf("%d", i))
		cell.Write(HashKey(key), &Item{
			Key:        key,
			ValueInt64: i,
			Type:       ItemInt64,
		})
	}

	t.Log("totalBytes:", cell.totalBytes)

	{
		t.Log("====")
		l := cell.list
		s := []string{}
		for e := l.head; e != nil; e = e.Next {
			s = append(s, fmt.Sprintf("%d", e.ValueInt64))
		}
		t.Log("{ " + strings.Join(s, ", ") + " }")
	}

	t.Log("mapping:", len(cell.mapping))
	s := []string{}
	for _, item := range cell.mapping {
		s = append(s, fmt.Sprintf("%d", item.ValueInt64))
	}
	t.Log("{ " + strings.Join(s, ", ") + " }")
}

func TestCell_MemoryUsage(t *testing.T) {
	//runtime.GOMAXPROCS(4)

	cell := NewCell()
	cell.LimitSize = 1024 * 1024 * 1024 * 1

	before := time.Now()

	wg := sync.WaitGroup{}
	wg.Add(4)

	for j := 0; j < 4; j ++ {
		go func(j int) {
			start := j * 50 * 10000
			for i := start; i < start+50*10000; i ++ {
				key := []byte(strconv.Itoa(i) + "VERY_LONG_STRING")
				cell.Write(HashKey(key), &Item{
					Key:        key,
					ValueInt64: int64(i),
					Type:       ItemInt64,
				})
			}
			wg.Done()
		}(j)
	}

	wg.Wait()
	t.Log("items:", len(cell.mapping))
	t.Log(time.Since(before).Seconds(), "s", "totalBytes:", cell.totalBytes/1024/1024, "M")
	//time.Sleep(10 * time.Second)
}

func BenchmarkCell_Write(b *testing.B) {
	runtime.GOMAXPROCS(1)

	cell := NewCell()

	for i := 0; i < b.N; i ++ {
		key := []byte(strconv.Itoa(i) + "_LONG_KEY_LONG_KEY_LONG_KEY_LONG_KEY")
		cell.Write(HashKey(key), &Item{
			Key:        key,
			ValueInt64: int64(i),
			Type:       ItemInt64,
		})
	}

	b.Log("items:", len(cell.mapping))
}

func TestCell_Read(t *testing.T) {
	cell := NewCell()

	cell.Write(1, &Item{
		ValueInt64: 1,
		ExpireAt:   time.Now().Unix() + 3600,
	})
	cell.Write(2, &Item{
		ValueInt64: 2,
		ExpireAt:   time.Now().Unix() + 3600,
	})
	cell.Write(3, &Item{
		ValueInt64: 3,
		ExpireAt:   time.Now().Unix() + 3600,
	})

	{
		s := []string{}
		cell.list.Range(func(item *Item) (goNext bool) {
			s = append(s, fmt.Sprintf("%d", item.ValueInt64))
			return true
		})
		t.Log("before:", s)
	}

	t.Log(cell.Read(1).ValueInt64)

	{
		s := []string{}
		cell.list.Range(func(item *Item) (goNext bool) {
			s = append(s, fmt.Sprintf("%d", item.ValueInt64))
			return true
		})
		t.Log("after:", s)
	}

	t.Log(cell.Read(2).ValueInt64)

	{
		s := []string{}
		cell.list.Range(func(item *Item) (goNext bool) {
			s = append(s, fmt.Sprintf("%d", item.ValueInt64))
			return true
		})
		t.Log("after:", s)
	}
}

func TestCell_Recycle(t *testing.T) {
	cell := NewCell()
	cell.Write(1, &Item{
		ValueInt64: 1,
		ExpireAt:   time.Now().Unix() - 1,
	})

	cell.Write(2, &Item{
		ValueInt64: 2,
		ExpireAt:   time.Now().Unix() + 1,
	})

	cell.Recycle()

	{
		s := []string{}
		cell.list.Range(func(item *Item) (goNext bool) {
			s = append(s, fmt.Sprintf("%d", item.ValueInt64))
			return true
		})
		t.Log("after:", s)
	}

	t.Log(cell.list.Len(), cell.totalBytes)
}
