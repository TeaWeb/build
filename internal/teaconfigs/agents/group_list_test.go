package agents

import (
	"github.com/iwind/TeaGo/logs"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"testing"
)

func TestSharedListConfig(t *testing.T) {
	config := SharedGroupList()
	t.Log(stringutil.JSONEncodePretty(config))
}

func TestSharedListConfig_Save(t *testing.T) {
	config := SharedGroupList()
	config.AddGroup(NewGroup("GROUP001"))
	err := config.Save()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSharedListConfig_Remove(t *testing.T) {
	config := SharedGroupList()
	config.RemoveGroup("0iqcKNj6zYCqfMoH")
	err := config.Save()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSharedListConfig2(t *testing.T) {
	config := SharedGroupList()
	//logs.PrintAsJSON(config, t)
	logs.PrintAsJSON(config.FindGroup(""), t)
}
