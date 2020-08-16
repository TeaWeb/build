package agents

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestFileChangeSource_Execute(t *testing.T) {
	source := NewFileChangeSource()
	source.Path = "/opt/test/cpu.sh"

	for i := 0; i < 2; i ++ {
		value, err := source.Execute(nil)
		if err != nil {
			t.Fatal(err)
		}

		logs.PrintAsJSON(value, t)
	}
}
