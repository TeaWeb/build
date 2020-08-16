package agents

import (
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestTeaWebSource_Execute(t *testing.T) {
	if !teatesting.RequireTeaWebServer() {
		return
	}

	source := NewTeaWebSource()
	source.API = "http://127.0.0.1:7777/api/monitor?TeaKey=z8O4MuXixbKH6aiVyZigYTxxovRblR3u"
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(value, t)
}
