package teacache

import (
	"fmt"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/assert"
	"github.com/iwind/TeaGo/logs"
	"runtime"
	"sync"
	"testing"
)

func TestFileManager(t *testing.T) {
	a := assert.NewAssertion(t)

	m := NewFileManager()
	m.dir = Tea.Root + "/cache"

	err := m.Write("123456", []byte("abc"))
	t.Log(err)
	a.IsNotNil(err)

	m.dir = Tea.Root + "/cache"
	logs.Println(m.dir)

	err = m.Write("123456", []byte("abcd"))
	a.IsNil(err)
	if err != nil {
		t.Log(err)
	}

	data, err := m.Read("123456")
	a.IsNil(err, func() string {
		return err.Error()
	})
	a.IsTrue(string(data) == "abcd", func() string {
		return "data:" + string(data)
	})
	t.Log(string(data))
}

func TestFileManagerConcurrent(t *testing.T) {
	m := NewFileManager()
	m.dir = Tea.TmpDir() + "/cache"

	wg := sync.WaitGroup{}
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func(i int) {
			_, err := m.Read("123456")
			if err != nil {
				//logs.Println(err)
			} else {
				logs.Println("read success")
			}

			err = m.Write("123456", []byte("abc"))
			if err != nil {
				//logs.Println(i, err)
			} else {
				logs.Println("write success")
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestFileManager_Stat(t *testing.T) {
	m := NewFileManager()
	m.dir = Tea.Root + "/./cache"
	t.Log(m.Stat())
}

func TestFileManager_Clean(t *testing.T) {
	m := NewFileManager()
	m.dir = Tea.Root + "/./cache"
	t.Log(m.Clean())
}

func TestFileManager_DeletePrefixes(t *testing.T) {
	m := NewFileManager()
	m.dir = Tea.Root + "/cache/"

	canWrite := false
	if canWrite {
		for i := 0; i < 3; i++ {
			err := m.Write("abc"+fmt.Sprintf("%03d", i), []byte("I AM DATA"))
			if err != nil {
				t.Fatal(err)
			}
		}
		_ = m.Write("bcd000", []byte("I AM BCD"))

		data, err := m.Read("abc000")
		if err != nil {
			t.Fatal(err)
		}
		t.Log("read:["+string(data)+"]", len(data), "bytes")
	}

	count, err := m.DeletePrefixes([]string{"abc"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("OK", count, "keys")
}

func BenchmarkFileManager_Read(b *testing.B) {
	runtime.GOMAXPROCS(1)

	m := NewFileManager()
	m.dir = Tea.Root + "/cache/"

	for i := 0; i < b.N; i++ {
		data, _ := m.Read("abc000")
		if len(data) == 0 {
			b.Fatal("invalid data:", string(data))
		}
	}
}
