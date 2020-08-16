package teamemory

import (
	"crypto/md5"
	"github.com/dchest/siphash"
	"strconv"
	"testing"
)

func TestItem_Size(t *testing.T) {
	item := &Item{
		ValueInt64: 1024,
		Key:        []byte("123"),
		ValueBytes: []byte("Hello, World"),
	}
	t.Log(item.Size())
}

func BenchmarkItem_Size(b *testing.B) {
	item := &Item{
		ValueInt64: 1024,
		Key:        []byte("123"),
		ValueBytes: []byte("Hello, World"),
	}
	for i := 0; i < b.N; i ++ {
		_ = item.Size()
	}
}

func TestItem_HashKey(t *testing.T) {
	t.Log(HashKey([]byte("2")))
}

func TestItem_siphash(t *testing.T) {
	result := siphash.Hash(0, 0, []byte("123456"))
	t.Log(result)
}

func TestItem_unique(t *testing.T) {
	m := map[uint64]bool{}
	for i := 0; i < 1000*10000; i ++ {
		s := "Hello,World,LONG KEY,LONG KEY,LONG KEY,LONG KEY" + strconv.Itoa(i)
		result := siphash.Hash(0, 0, []byte(s))
		_, ok := m[result]
		if ok {
			t.Log("found same", i)
			break
		} else {
			m[result] = true
		}
	}

	t.Log(siphash.Hash(0, 0, []byte("01")))
	t.Log(siphash.Hash(0, 0, []byte("10")))
}

func BenchmarkItem_HashKeyMd5(b *testing.B) {
	for i := 0; i < b.N; i ++ {
		h := md5.New()
		h.Write([]byte("HELLO_KEY_" + strconv.Itoa(i)))
		_ = h.Sum(nil)
	}
}

func BenchmarkItem_siphash(b *testing.B) {
	for i := 0; i < b.N; i ++ {
		_ = siphash.Hash(0, 0, []byte("HELLO_KEY_"+strconv.Itoa(i)))
	}
}
