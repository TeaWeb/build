package teawaf

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
	"time"
)

func TestIPTable_MatchIP(t *testing.T) {
	a := assert.NewAssertion(t)

	{
		table := NewIPTable()
		err := table.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsFalse(table.Match("192.168.1.100", 8080))
	}

	{
		table := NewIPTable()
		table.IP = "*"
		table.Port = "8080"
		err := table.Init()
		if err != nil {
			t.Fatal(err)
		}
		t.Log("port:", table.minPort, table.maxPort)
		a.IsTrue(table.Match("192.168.1.100", 8080))
		a.IsFalse(table.Match("192.168.1.100", 8081))
	}

	{
		table := NewIPTable()
		table.IP = "*"
		table.Port = "8080-8082"
		err := table.Init()
		if err != nil {
			t.Fatal(err)
		}
		t.Log("port:", table.minPort, table.maxPort)
		a.IsTrue(table.Match("192.168.1.100", 8080))
		a.IsTrue(table.Match("192.168.1.100", 8081))
		a.IsFalse(table.Match("192.168.1.100", 8083))
	}

	{
		table := NewIPTable()
		table.IP = "*"
		table.Port = "*-8082"
		err := table.Init()
		if err != nil {
			t.Fatal(err)
		}
		t.Log("port:", table.minPort, table.maxPort)
		a.IsTrue(table.Match("192.168.1.100", 8079))
		a.IsTrue(table.Match("192.168.1.100", 8080))
		a.IsTrue(table.Match("192.168.1.100", 8081))
		a.IsFalse(table.Match("192.168.1.100", 8083))
	}

	{
		table := NewIPTable()
		table.IP = "*"
		table.Port = "8080-*"
		err := table.Init()
		if err != nil {
			t.Fatal(err)
		}
		t.Log("port:", table.minPort, table.maxPort)
		a.IsFalse(table.Match("192.168.1.100", 8079))
		a.IsTrue(table.Match("192.168.1.100", 8080))
		a.IsTrue(table.Match("192.168.1.100", 8081))
		a.IsTrue(table.Match("192.168.1.100", 8083))
	}

	{
		table := NewIPTable()
		table.IP = "*"
		table.Port = "*"
		err := table.Init()
		if err != nil {
			t.Fatal(err)
		}
		t.Log("port:", table.minPort, table.maxPort)
		a.IsTrue(table.Match("192.168.1.100", 8079))
		a.IsTrue(table.Match("192.168.1.100", 8080))
		a.IsTrue(table.Match("192.168.1.100", 8081))
		a.IsTrue(table.Match("192.168.1.100", 8083))
	}

	{
		table := NewIPTable()
		table.IP = "192.168.1.100"
		table.Port = "*"
		err := table.Init()
		if err != nil {
			t.Fatal(err)
		}
		t.Log("port:", table.minPort, table.maxPort)
		a.IsTrue(table.Match("192.168.1.100", 8080))
	}

	{
		table := NewIPTable()
		table.IP = "192.168.1.99-192.168.1.101"
		table.Port = "*"
		err := table.Init()
		if err != nil {
			t.Fatal(err)
		}
		t.Log("port:", table.minPort, table.maxPort)
		a.IsTrue(table.Match("192.168.1.100", 8080))
	}

	{
		table := NewIPTable()
		table.IP = "192.168.1.99/24"
		table.Port = "*"
		err := table.Init()
		if err != nil {
			t.Fatal(err)
		}
		t.Log("ip:", table.ipRange)
		a.IsTrue(table.Match("192.168.1.100", 8080))
		a.IsFalse(table.Match("192.168.2.100", 8080))
	}

	{
		table := NewIPTable()
		table.IP = "192.168.1.99/24"
		table.TimeTo = time.Now().Unix() - 10
		table.Port = "*"
		err := table.Init()
		if err != nil {
			t.Fatal(err)
		}
		a.IsFalse(table.Match("192.168.1.100", 8080))
		a.IsFalse(table.Match("192.168.2.100", 8080))
	}
}
