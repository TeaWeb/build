package board

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
	"github.com/TeaWeb/build/internal/teastats"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

type ChartAction actions.Action

// 图表详情
func (this *ChartAction) Run(params struct {
	ServerId  string
	WidgetId  string
	ChartId   string
	BoardType string
}) {
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
	this.Data["canUpdate"] = Tea.IsTesting() || !strings.HasPrefix(widget.Id, "teaweb.")

	items := []maps.Map{}
	for _, r := range chart.Requirements {
		filter := teastats.FindSharedFilter(r)
		if filter != nil {
			items = append(items, maps.Map{
				"name": filter.Name(),
				"code": r,
			})
		}
	}
	this.Data["items"] = items

	this.Show()
}
