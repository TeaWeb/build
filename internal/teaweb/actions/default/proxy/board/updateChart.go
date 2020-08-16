package board

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
	"github.com/TeaWeb/build/internal/teastats"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateChartAction actions.Action

// 修改图表
func (this *UpdateChartAction) Run(params struct {
	From      string
	ServerId  string
	WidgetId  string
	ChartId   string
	BoardType string
}) {
	this.Data["from"] = params.From
	this.Data["boardType"] = params.BoardType

	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到server")
	}
	this.Data["server"] = maps.Map{
		"id": server.Id,
	}

	widget := widgets.NewWidgetFromId(params.WidgetId)
	if widget == nil {
		this.Fail("找不到Widget")
	}

	this.Data["widget"] = widget

	chart := widget.FindChart(params.ChartId)
	if chart == nil {
		this.Fail("找不到Chart")
	}

	this.Data["chart"] = chart
	this.Data["items"] = teastats.FindAllStatFilters()

	this.Show()
}

// 保存提交
func (this *UpdateChartAction) RunPost(params struct {
	ServerId       string
	WidgetId       string
	ChartId        string
	Name           string
	Description    string
	Columns        uint8
	Items          []string
	JavascriptCode string
	On             bool
	Must           *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入名称")

	widget := widgets.NewWidgetFromId(params.WidgetId)
	if widget == nil {
		this.Fail("找不到Widget")
	}

	chart := widget.FindChart(params.ChartId)
	if chart == nil {
		this.Fail("找不到Chart")
	}

	oldRequirements := append([]string{}, chart.Requirements...)

	chart.On = params.On
	chart.Name = params.Name
	chart.Description = params.Description
	chart.Columns = params.Columns
	chart.Requirements = params.Items
	chart.Type = "javascript"
	chart.Options = maps.Map{
		"code": params.JavascriptCode,
	}

	err := widget.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 重启统计指标
	if !this.equalSlice(oldRequirements, chart.Requirements) {
		for _, s := range teaconfigs.LoadServerConfigsFromDir(Tea.ConfigDir()) {
			if (s.RealtimeBoard != nil && s.RealtimeBoard.HasChart(chart.Id)) || (s.StatBoard != nil && s.StatBoard.HasChart(chart.Id)) {
				proxyutils.ReloadServerStats(s.Id)
			}
		}
	}

	this.Success()
}

// 判断两个string slice是否相同
func (this *UpdateChartAction) equalSlice(slice1 []string, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for k, v := range slice1 {
		if v != slice2[k] {
			return false
		}
	}
	return true
}
