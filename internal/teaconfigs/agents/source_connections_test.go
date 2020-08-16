package agents

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestConnectionsSource_Execute(t *testing.T) {
	source := NewConnectionsSource()
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(value, t)
}
