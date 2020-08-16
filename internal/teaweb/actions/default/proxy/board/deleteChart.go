package board

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
)

type DeleteChartAction actions.Action

// 删除图表
func (this *DeleteChartAction) Run(params struct {
	ServerId string
	WidgetId string
	ChartId  string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到server")
	}

	widget := widgets.NewWidgetFromId(params.WidgetId)
	if widget == nil {
		this.Fail("找不到Widget")
	}

	chart := widget.FindChart(params.ChartId)
	if chart == nil {
		this.Fail("找不到Chart")
	}

	widget.RemoveChart(params.ChartId)
	if len(widget.Charts) > 0 {
		err := widget.Save()
		if err != nil {
			this.Fail("保存失败：" + err.Error())
		}
	} else {
		err := widget.Delete()
		if err != nil {
			this.Fail("保存失败：" + err.Error())
		}
	}

	// 移除所有Server中的相关记录
	chartId := params.ChartId
	for _, s := range teaconfigs.LoadServerConfigsFromDir(Tea.ConfigDir()) {
		contains := false
		if s.RealtimeBoard != nil && s.RealtimeBoard.HasChart(chartId) {
			contains = true
			s.RealtimeBoard.RemoveChart(chartId)
		}

		if s.StatBoard != nil && s.StatBoard.HasChart(chartId) {
			contains = true
			s.StatBoard.RemoveChart(chartId)
		}

		if contains {
			err := s.Save()
			if err != nil {
				logs.Error(err)
			}

			// 重启统计
			if len(chart.Requirements) > 0 {
				proxyutils.ReloadServerStats(s.Id)
			}
		}
	}

	this.Success()
}
