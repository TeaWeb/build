package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/board/scripts"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
)

type PreviewItemChartAction actions.Action

// 预览图表
func (this *PreviewItemChartAction) RunPost(params struct {
	AgentId string
	AppId   string
	ItemId  string

	Name      string
	Columns   uint8
	ChartType string

	HTMLCode string `alias:"htmlCode"`

	PieParam string
	PieLimit int

	LineParams []string
	LineFills  []int
	LineColors []string
	LineNames  []string
	LineMax    float64

	URL string `alias:"urlURL"`

	JavascriptCode string

	Must *actions.Must
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到要修改的Agent")
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		this.Fail("找不到要修改的App")
	}

	item := app.FindItem(params.ItemId)
	if item == nil {
		this.Fail("找不到要操作的Item")
	}

	chart := widgets.NewChart()
	chart.Name = params.Name
	chart.On = true
	chart.Columns = params.Columns
	chart.Type = params.ChartType

	switch params.ChartType {
	case "html":
		options := &widgets.HTMLChart{}
		options.HTML = params.HTMLCode
		err := teautils.ObjectToMapJSON(options, &chart.Options)
		if err != nil {
			logs.Error(err)
		}
	case "url":
		options := &widgets.URLChart{}
		options.URL = params.URL
		err := teautils.ObjectToMapJSON(options, &chart.Options)
		if err != nil {
			logs.Error(err)
		}
	case "pie":
		options := &widgets.PieChart{}
		options.Param = params.PieParam
		options.Limit = params.PieLimit
		err := teautils.ObjectToMapJSON(options, &chart.Options)
		if err != nil {
			logs.Error(err)
		}
	case "line":
		options := &widgets.LineChart{}
		options.Max = params.LineMax
		for index, param := range params.LineParams {
			line := widgets.NewLine()
			line.Param = param
			if index < len(params.LineFills) {
				line.IsFilled = params.LineFills[index] > 0
			}
			if index < len(params.LineColors) {
				line.Color = params.LineColors[index]
			}
			if index < len(params.LineNames) {
				line.Name = params.LineNames[index]
			}
			options.AddLine(line)
		}
		err := teautils.ObjectToMapJSON(options, &chart.Options)
		if err != nil {
			logs.Error(err)
		}
	case "javascript":
		options := &widgets.JavascriptChart{}
		options.Code = params.JavascriptCode
		err := teautils.ObjectToMapJSON(options, &chart.Options)
		if err != nil {
			logs.Error(err)
		}
	}

	c, err := chart.AsObject()
	if err != nil {
		this.Fail("发现错误：" + err.Error())
	}

	code, err := c.AsJavascript(map[string]interface{}{
		"name":    params.Name,
		"columns": params.Columns,
	})

	dbEnabled := teadb.SharedDB().Test() == nil
	engine := scripts.NewEngine()
	engine.SetDBEnabled(dbEnabled)
	engine.SetCache(false)

	engine.SetContext(&scripts.Context{
		Agent: agent,
		App:   app,
		Item:  item,
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

	this.Data["charts" ] = engine.Charts()
	this.Data["output"] = engine.Output()
	this.Success()
}
