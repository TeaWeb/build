package widgets

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestLoadAllWidgets(t *testing.T) {
	widgets := LoadAllWidgets()
	for _, widget := range widgets {
		t.Log(widget)
		logs.PrintAsJSON(widget)
	}
}

func TestNewWidgetFromId(t *testing.T) {
	{
		widget := NewWidgetFromId("LjPKIDfrkThkO3kX")
		t.Log(widget)
	}

	{
		widget := NewWidgetFromId("123456")
		t.Log(widget)
	}
}
