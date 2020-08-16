package db

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestLoadMySQLConfig(t *testing.T) {
	config, err := LoadMySQLConfig()
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(config, t)
}
