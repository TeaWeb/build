package agents

import (
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestPingSource_Execute(t *testing.T) {
	if !teatesting.RequireDNS() {
		return
	}

	source := NewPingSource()
	source.Host = "teaos.cn"
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}

	logs.PrintAsJSON(value, t)
}

func TestPingSource_Execute_Timeout(t *testing.T) {
	source := NewPingSource()
	source.Host = "123.123.123.123"
	value, err := source.Execute(nil)
	if err != nil {
		t.Log("error:", err.Error())
		return
	}

	logs.PrintAsJSON(value, t)
}

func TestPingSource_Execute_Unknown(t *testing.T) {
	source := NewPingSource()
	source.Host = "abcdefghi.com"
	value, err := source.Execute(nil)
	if err != nil {
		t.Log("error:", err.Error())
		return
	}

	logs.PrintAsJSON(value, t)
}
