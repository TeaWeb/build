package teacache

import (
	"github.com/iwind/TeaGo/Tea"
	"testing"
	"time"
)

func TestLevelDBManager_Write(t *testing.T) {
	m := NewLevelDBManager()
	m.SetOptions(map[string]interface{}{
		"dir": Tea.Root + "/cache",
	})
	m.Life = 30 * time.Second
	t.Log(m.Write("hello", []byte("world123")))
	_ = m.Close()
}

func TestLevelDBManager_Read(t *testing.T) {
	m := NewLevelDBManager()
	m.SetOptions(map[string]interface{}{
		"dir": Tea.Root + "/cache",
	})
	data, err := m.Read("hello")
	if err != nil && err != ErrNotFound {
		t.Fatal(err)
	}

	t.Log(string(data))
	_ = m.Close()
}

func TestLevelDBManager_Stat(t *testing.T) {
	m := NewLevelDBManager()
	m.SetOptions(map[string]interface{}{
		"dir": Tea.Root + "/cache",
	})
	t.Log(m.Stat())
	_ = m.Close()
}

func TestLevelDBManager_CleanExpired(t *testing.T) {
	m := NewLevelDBManager()
	m.SetOptions(map[string]interface{}{
		"dir": Tea.Root + "/cache",
	})
	m.Life = 30 * time.Second
	err := m.CleanExpired()
	if err != nil {
		t.Fatal(err)
	}
	_ = m.Close()
}

func TestLevelDBManager_Clean(t *testing.T) {
	m := NewLevelDBManager()
	m.SetOptions(map[string]interface{}{
		"dir": Tea.Root + "/cache",
	})
	m.Life = 30 * time.Second
	err := m.Clean()
	if err != nil {
		t.Fatal(err)
	}
	_ = m.Close()
}

func TestLevelDBManager_DeletePrefixes(t *testing.T) {
	m := NewLevelDBManager()
	m.SetOptions(map[string]interface{}{
		"dir": Tea.Root + "/cache",
	})
	canWrite := false
	if canWrite {
		err := m.Write("abc000", []byte("Hello"))
		if err != nil {
			t.Fatal(err)
		}
	}
	count, err := m.DeletePrefixes([]string{"abc"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(count, "keys deleted")
}
