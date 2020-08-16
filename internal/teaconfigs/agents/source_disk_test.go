package agents

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestDiskSource_Execute(t *testing.T) {
	source := NewDiskSource()
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(value, t)
}

func TestDiskSource_Execute_Filter(t *testing.T) {
	source := NewDiskSource()
	source.ContainsAllMountPoints = true
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(value, t)
}
