package teaconfigs

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestWebConfig_LoadWebConfig(t *testing.T) {
	config, err := LoadWebConfig()
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(config, t)
}

func TestWebConfig_Save(t *testing.T) {
	config, err := LoadWebConfig()
	if err != nil {
		t.Fatal(err)
	}
	err = config.Save()
	if err != nil {
		t.Fatal(err)
	}
}
