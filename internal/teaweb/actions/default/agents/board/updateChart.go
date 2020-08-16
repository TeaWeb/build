package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/iwind/TeaGo/actions"
)

type UpdateChartAction actions.Action

// 修改图表设置
func (this *UpdateChartAction) RunPost(params struct {
	AgentId string
	AppId   string
	ItemId  string
	ChartId string

	Name     string
	TimeType string
	TimePast string
	DayFrom  string
	DayTo    string
}) {
	board := agents.NewAgentBoard(params.AgentId)
	if board == nil {
		this.Fail("找不到Agent")
	}

	chart := board.FindChart(params.ChartId)
	if chart == nil {
		this.Fail("找不到要修改的图表")
	}

	chart.Name = params.Name
	chart.TimeType = params.TimeType

	switch params.TimeType {
	case "past":
		chart.TimePast = params.TimePast
	case "range":
		chart.DayFrom = params.DayFrom
		chart.DayTo = params.DayTo
	}

	err := board.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
