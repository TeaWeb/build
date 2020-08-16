package db

import "testing"

func TestSharedDBConfig(t *testing.T) {
	config := SharedDBConfig()
	t.Log(config.Type)
}

func TestDBConfig_Save(t *testing.T) {
	config := SharedDBConfig()
	config.Type = DBTypeMongo
	err := config.Save()
	if err != nil {
		t.Fatal(err)
	}
}
