package teastats

import (
	"testing"
	"time"
)

func TestKVStorage_Set(t *testing.T) {
	if sharedKV == nil {
		sharedKV = NewKVStorage("stat.leveldb")
		if sharedKV == nil {
			return
		}
	}
	err := sharedKV.Set("hello", "value", time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	value, err := sharedKV.Get("hello")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(value)
	_ = sharedKV.Close()
}

func TestKVStorage_Has(t *testing.T) {
	if sharedKV == nil {
		sharedKV = NewKVStorage("stat.leveldb")
		if sharedKV == nil {
			return
		}
	}
	t.Log(sharedKV.Has("hello"))
}
