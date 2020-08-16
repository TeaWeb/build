package agents

import (
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestNginxStatusSource_Execute(t *testing.T) {
	if !teatesting.RequireNginxStatus() {
		return
	}

	source := NewNginxStatusSource()
	source.URL = "http://127.0.0.1:8888/nginx_status"
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(value, t)
}
