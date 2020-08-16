package agents

import (
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestSocketConnectivitySource_Execute(t *testing.T) {
	if !teatesting.RequireHTTPServer() {
		return
	}

	source := NewSocketConnectivitySource()
	source.Address = "127.0.0.1:9991"
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(value, t)
}

func TestSocketConnectivitySource_Execute_2(t *testing.T) {
	if !teatesting.RequirePort(27017) {
		return
	}

	source := NewSocketConnectivitySource()
	source.Address = "127.0.0.1:27017"
	source.Network = "tcp"
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(value, t)
}
