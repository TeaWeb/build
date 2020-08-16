package agents

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestDateSource_Execute(t *testing.T) {
	source := NewDateSource()
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(value, t)
}
