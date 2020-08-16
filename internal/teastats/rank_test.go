package teastats

import (
	"fmt"
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestRank_Top(t *testing.T) {
	rank := NewRank(10, 1000000)

	j := 0
	for i := 0; i < 500*10000; i++ {
		if i%12 == 0 {
			rank.Add("192.168.1.a")
		} else if i%11 == 0 {
			rank.Add("192.168.1.b")
		} else if i%10 == 0 {
			rank.Add("192.168.1.c")
		} else if i%9 == 0 {
			rank.Add("192.168.1.d")
		} else if i%8 == 0 {
			rank.Add("192.168.1.e")
		} else if i%7 == 0 {
			rank.Add("192.168.1.f")
		} else if i%6 == 0 {
			rank.Add("192.168.1.g")
		} else if i%5 == 0 {
			rank.Add("192.168.1.h")
		} else if i%4 == 0 {
			rank.Add("192.168.1.i")
		} else if i%3 == 0 {
			rank.Add("192.168.1.j")
		} else if i%2 == 0 {
			rank.Add("192.168.1.k")
		} else {
			rank.Add(fmt.Sprintf("192.168.1.%d", j))
			j++
		}
	}

	logs.PrintAsJSON(rank.Top(), t)
	t.Log(len(rank.top))
	t.Log(len(rank.buffer))
	t.Log("dirty keys:", len(rank.dirtyKeys))

	db := NewKVStorage("test.leveldb").db
	rank.Save(db, "hello")
}

func TestRank_Load(t *testing.T) {
	kv := NewKVStorage("test.leveldb")
	if kv == nil {
		return
	}
	defer func() {
		_ = kv.Close()
	}()
	db := kv.db

	rank := NewRank(10, 1000000)
	rank.Load(db, "hello")
	logs.PrintAsJSON(rank.Top(), t)
	t.Log(len(rank.top))
	t.Log(len(rank.buffer))
}
