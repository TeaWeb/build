package agents

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestProcessesSource_Execute(t *testing.T) {
	source := NewProcessesSource()
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(value, t)
}
