package apps

import (
	stringutil "github.com/iwind/TeaGo/utils/string"
	"testing"
)

func TestExportChartDataAction_extractTitles(t *testing.T) {
	a := new(ExportChartDataAction)

	{
		m := map[string]interface{}{}
		t.Log(a.extractTitles("", m))
	}

	{
		m := map[string]interface{}{
			"a": "1",
			"b": "2",
		}
		t.Log(a.extractTitles("", m))
	}

	{
		m := map[string]interface{}{
			"a": "1",
			"b": "2",
			"c": map[string]interface{}{
				"d": "3",
				"e": "4",
				"f": map[string]interface{}{
					"g": "5",
				},
			},
		}
		t.Log(stringutil.JSONEncodePretty(a.extractTitles("", m)))
	}
}

func TestExportChartDataAction_formatTime(t *testing.T) {
	a := new(ExportChartDataAction)
	t.Log(a.formatTime("20191011010203"))
}
