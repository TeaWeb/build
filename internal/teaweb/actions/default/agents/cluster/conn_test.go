package cluster

import (
	"encoding/json"
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"testing"
)

func TestConnAction_Run(t *testing.T) {
	if teatesting.IsGlobal() {
		return
	}

	t1 := actions.NewTesting(new(ConnAction))
	t1.Params(actions.Params{
		"hosts": []string{"192.168.2.33", "127.0.0.1", "teaos.cn"},
		"port":  []string{"22"},
	})
	result := t1.Run(t).Data

	m := maps.Map{}
	err := json.Unmarshal(result, &m)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(m)
}
