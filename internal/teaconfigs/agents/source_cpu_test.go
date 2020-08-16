package agents

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
	"time"
)

func TestCPUSource_Execute(t *testing.T) {
	time.Sleep(1 * time.Second)

	before := time.Now()
	for i := 0; i < 3; i ++ {
		source := NewCPUSource()
		value, err := source.Execute(nil)
		if err != nil {
			t.Fatal(err)
		}
		logs.PrintAsJSON(value, t)
		t.Log(time.Since(before).Seconds(), "s")
	}
}
