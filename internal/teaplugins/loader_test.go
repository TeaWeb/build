package teaplugins

import (
	"github.com/TeaWeb/plugin/pkg/messages"
	plugins2 "github.com/TeaWeb/plugin/pkg/plugins"
	"github.com/iwind/TeaGo/files"
	"os"
	"testing"
)

func TestLoader_Load(t *testing.T) {
	path := os.Getenv("GOPATH") + "/src/github.com/TeaWeb/plugin/main/demo.plugin"
	f := files.NewFile(path)
	if !f.Exists() {
		return
	}
	loader := NewLoader(path)
	err := loader.Load()
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoader_CallAction(t *testing.T) {
	loader := NewLoader("")
	action := new(messages.RegisterPluginAction)
	action.Plugin = new(plugins2.Plugin)
	err := loader.CallAction(action, 1)
	if err != nil {
		t.Fatal(err)
	}
}
