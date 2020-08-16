package api

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"testing"
)

func TestAPIScript_Save(t *testing.T) {
	s := NewAPIScript()
	s.Code = "var api = API(\"$\");"
	err := s.Save()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("success", s.Filename)

	if len(s.Filename) > 0 {
		files.NewFile(Tea.ConfigFile(s.Filename)).Delete()
	}
}
