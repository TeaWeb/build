package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
)

type ChartsAction actions.Action

// 图表列表
func (this *ChartsAction) Run(params struct {
	AgentId string
}) {
	this.Data["agentId"] = params.AgentId
	this.Data["tabbar"] = "charts"

	charts := []maps.Map{}

	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	board := agents.NewAgentBoard(params.AgentId)
	if board == nil {
		this.Fail("无法创建看板的配置文件")
	}

	chartMapping := map[string]maps.Map{}
	for _, app := range agent.Apps {
		if !app.On {
			continue
		}

		for _, item := range app.Items {
			if !item.On {
				continue
			}

			if len(item.Charts) == 0 {
				continue
			}

			for _, chart := range item.Charts {
				if !chart.On {
					continue
				}

				isUsing := board.FindChart(chart.Id) != nil
				info := maps.Map{
					"id":       chart.Id,
					"name":     chart.Name,
					"typeName": widgets.FindChartTypeName(chart.Type),
					"app": maps.Map{
						"id":   app.Id,
						"name": app.Name,
					},
					"item": maps.Map{
						"id":   item.Id,
						"name": item.Name,
					},
					"isUsing": isUsing,
					"columns": chart.Columns,
				}
				chartMapping[chart.Id] = info
				charts = append(charts, info)
			}
		}
	}

	this.Data["charts"] = charts

	usingCharts := []maps.Map{}
	if board != nil {
		hasDeleted := false
		for _, c := range board.Charts {
			info, found := chartMapping[c.ChartId]
			if found {
				usingCharts = append(usingCharts, info)
			} else {
				hasDeleted = true
			}
		}

		// 删除已删除的
		if hasDeleted {
			leftCharts := []*agents.BoardChart{}
			for _, c := range board.Charts {
				_, found := chartMapping[c.ChartId]
				if found {
					leftCharts = append(leftCharts, c)
				}
			}
			board.Charts = leftCharts
			err := board.Save()
			if err != nil {
				logs.Error(err)
			}
		}
	}
	this.Data["usingCharts"] = usingCharts

	this.Show()
}
