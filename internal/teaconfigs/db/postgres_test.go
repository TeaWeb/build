package db

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestLoadPostgresConfig(t *testing.T) {
	config, err := LoadPostgresConfig()
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(config, t)
}
