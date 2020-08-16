package agents

import (
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestMySQLSource_Execute(t *testing.T) {
	if !teatesting.RequireMySQL() {
		return
	}

	source := NewMySQLSource()
	source.TimeoutSeconds = 10
	source.Addr = "127.0.0.1"
	source.Username = "root"
	source.Password = "123456"
	source.DatabaseName = "teaweb"
	source.SQL = "SELECT * FROM tea_accessLogs"
	values, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(values, t)
}
