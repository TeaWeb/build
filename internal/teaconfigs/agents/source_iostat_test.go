package agents

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
	"time"
)

func TestIOStatSource_Execute(t *testing.T) {
	source := NewIOStatSource()

	{
		value, err := source.Execute(nil)
		if err != nil {
			t.Fatal(err)
		}

		logs.PrintAsJSON(value, t)
	}

	time.Sleep(5 * time.Second)
	{
		value, err := source.Execute(nil)
		if err != nil {
			t.Fatal(err)
		}

		logs.PrintAsJSON(value, t)
	}
}
