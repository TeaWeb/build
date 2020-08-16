package agents

import (
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestDockerSource_Execute(t *testing.T) {
	if !teatesting.RequireDocker() {
		return
	}

	source := NewDockerSource()
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(value, t)
}
