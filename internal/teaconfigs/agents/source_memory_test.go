package agents

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestMemorySource_Execute(t *testing.T) {
	source := NewMemorySource()
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}

	logs.PrintAsJSON(value, t)
}
