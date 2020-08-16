package agents

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestAppProcessesSource_Execute(t *testing.T) {
	source := NewAppProcessesSource()
	source.CmdlineKeyword = "mysql"
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}

	logs.PrintAsJSON(value, t)
}

func TestAppProcessesSource_Execute2(t *testing.T) {
	source := NewAppProcessesSource()
	source.CmdlineKeyword = "mysql123"
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}

	logs.PrintAsJSON(value, t)
}
