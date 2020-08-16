package agents

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestJSONFileSource_Execute(t *testing.T) {
	a := assert.NewAssertion(t)

	source := NewFileSource()
	a.IsNotNil(source.Validate())

	source.Path = Tea.ConfigFile("server.conf")
	source.DataFormat = SourceDataFormatYAML
	a.IsNil(source.Validate())

	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(value)
	}
}
