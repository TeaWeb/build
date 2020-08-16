package apps

import (
	"encoding/csv"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/board/scripts"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"strings"
)

type ExportChartDataAction actions.Action

// 导出图表数据
func (this *ExportChartDataAction) RunGet(params struct {
	Name     string
	AgentId  string
	AppId    string
	ItemId   string
	ChartId  string
	TimeType string
	TimePast string
	DayFrom  string
	DayTo    string

	Export string

	Must *actions.Must
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		this.Fail("找不到App")
	}

	item := app.FindItem(params.ItemId)
	if item == nil {
		this.Fail("找不到Item")
	}

	chart := item.FindChart(params.ChartId)
	if chart == nil {
		this.Fail("找不到Chart")
	}

	board := agents.NewAgentBoard(params.AgentId)
	if board == nil {
		this.Fail("找不到Board")
	}

	boardChart := board.FindChart(params.ChartId)
	if boardChart == nil {
		this.Fail("找不到BoardChart")
	}

	o, err := chart.AsObject()
	if err != nil {
		this.Fail("数据错误：" + err.Error())
	}

	code, err := o.AsJavascript(maps.Map{
		"name":    params.Name,
		"columns": chart.Columns,
	})
	if err != nil {
		this.Fail("数据错误：" + err.Error())
	}

	dbEnabled := teadb.SharedDB().Test() == nil
	engine := scripts.NewEngine()
	engine.SetDBEnabled(dbEnabled)
	engine.SetCache(false)

	if lists.ContainsString([]string{"data", "csv"}, params.Export) {
		engine.Exporting()
	}

	engine.SetContext(&scripts.Context{
		Agent:    agent,
		App:      app,
		Item:     item,
		TimeType: params.TimeType,
		TimePast: params.TimePast,
		DayFrom:  params.DayFrom,
		DayTo:    params.DayTo,
	})

	widgetCode := `var widget = new widgets.Widget({
	"name": "看板",
	"requirements": ["db"]
});

widget.run = function () {
`
	widgetCode += "{\n" + code + "\n}\n"
	widgetCode += `
};
`

	err = engine.RunCode(widgetCode)
	if err != nil {
		this.Fail("发生错误：" + err.Error())
	}

	this.Data["output"] = engine.Output()

	result := engine.Result()
	if !types.IsSlice(result) {
		this.Fail("只有数组数据才能导出")
	}

	exportTitles := []string{}
	lists.Each(result, func(k int, v interface{}) {
		m := map[string]interface{}{}
		err := teautils.ObjectToMapJSON(v, &m)
		if err != nil {
			return
		}

		value, ok := m["value"]
		if !ok || value == nil {
			return
		}

		valueMap, ok := value.(map[string]interface{})
		if !ok {
			return
		}

		titles := this.extractTitles("", valueMap)
		for _, title := range titles {
			if !lists.ContainsString(exportTitles, title) {
				exportTitles = append(exportTitles, title)
			}
		}
	})

	this.AddHeader("Content-Disposition", "attachment; filename=\"chart.data.csv\";")
	writer := csv.NewWriter(this.ResponseWriter)
	err = writer.Write(append([]string{"Time"}, exportTitles...))
	if err != nil {
		logs.Error(err)
	}

	lists.Each(result, func(k int, v interface{}) {
		m := map[string]interface{}{}
		err := teautils.ObjectToMapJSON(v, &m)
		if err != nil {
			return
		}

		value, ok := m["value"]
		if !ok || value == nil {
			return
		}

		valueMap, ok := value.(map[string]interface{})
		if !ok {
			return
		}

		lineValues := []string{}

		switch engine.Context().TimeUnit {
		case teaconfigs.TimeUnitSecond:
			lineValues = append(lineValues, this.formatTime(types.String(teautils.Get(m, []string{"timeFormat", "second"}))))
		case teaconfigs.TimeUnitMinute:
			lineValues = append(lineValues, this.formatTime(types.String(teautils.Get(m, []string{"timeFormat", "minute"}))))
		case teaconfigs.TimeUnitHour:
			lineValues = append(lineValues, this.formatTime(types.String(teautils.Get(m, []string{"timeFormat", "hour"}))))
		case teaconfigs.TimeUnitDay:
			lineValues = append(lineValues, this.formatTime(types.String(teautils.Get(m, []string{"timeFormat", "day"}))))
		case teaconfigs.TimeUnitMonth:
			lineValues = append(lineValues, this.formatTime(types.String(teautils.Get(m, []string{"timeFormat", "month"}))))
		case teaconfigs.TimeUnitYear:
			lineValues = append(lineValues, this.formatTime(types.String(teautils.Get(m, []string{"timeFormat", "year"}))))
		default:
			lineValues = append(lineValues, this.formatTime(types.String(teautils.Get(m, []string{"timeFormat", "minute"}))))
		}

		for _, title := range exportTitles {
			v2 := stringutil.JSONEncode(teautils.Get(valueMap, strings.Split(title, ".")))
			lineValues = append(lineValues, v2)
		}
		err = writer.Write(lineValues)
		if err != nil {
			logs.Error(err)
		}
	})

	writer.Flush()
}

func (this *ExportChartDataAction) extractTitles(prefix string, m map[string]interface{}) (result []string) {
	for k, v := range m {
		if v == nil {
			continue
		}

		key := k
		if len(prefix) != 0 {
			key = prefix + "." + key
		}

		if v1, ok := v.(map[string]interface{}); ok {
			result = append(result, this.extractTitles(key, v1)...)
			continue
		}

		if !lists.ContainsString(result, key) {
			result = append(result, key)
		}
	}

	return
}

func (this *ExportChartDataAction) formatTime(t string) string {
	if len(t) <= 4 { // year
		return t
	}
	if len(t) == 6 { // month
		return t[:4] + "-" + t[4:6]
	}
	if len(t) == 8 { // day
		return this.formatTime(t[:6]) + "-" + t[6:8]
	}
	if len(t) == 10 { // hour
		return this.formatTime(t[:8]) + " " + t[8:10]
	}
	if len(t) == 12 { // minute
		return this.formatTime(t[:10]) + ":" + t[10:12]
	}
	if len(t) == 14 { // second
		return this.formatTime(t[:12]) + ":" + t[12:14]
	}
	return t
}
