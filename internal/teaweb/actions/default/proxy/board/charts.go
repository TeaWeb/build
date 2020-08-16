package board

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
	"github.com/TeaWeb/build/internal/teastats"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type ChartsAction actions.Action

// 图表
func (this *ChartsAction) Run(params struct {
	ServerId  string
	BoardType string
}) {
	if len(params.BoardType) == 0 {
		params.BoardType = "realtime"
	}

	this.Data["boardType"] = params.BoardType

	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	this.Data["server"] = maps.Map{
		"id": server.Id,
	}

	// 正在使用的图表
	usingCharts := []maps.Map{}
	if params.BoardType == "realtime" {
		if server.RealtimeBoard != nil {
			for _, c := range server.RealtimeBoard.Charts {
				widget, chart := c.FindChart()
				if chart == nil {
					continue
				}
				usingCharts = append(usingCharts, maps.Map{
					"id":           chart.Id,
					"name":         chart.Name,
					"description":  chart.Description,
					"requirements": chart.Requirements,
					"columns":      chart.Columns,
					"on":           chart.On,
					"widget": maps.Map{
						"id":      widget.Id,
						"author":  widget.Author,
						"version": widget.Version,
					},
				})
			}
		}
	} else {
		if server.StatBoard != nil {
			for _, c := range server.StatBoard.Charts {
				widget, chart := c.FindChart()
				if chart == nil {
					continue
				}

				usingCharts = append(usingCharts, maps.Map{
					"id":           chart.Id,
					"name":         chart.Name,
					"description":  chart.Description,
					"requirements": chart.Requirements,
					"columns":      chart.Columns,
					"on":           chart.On,
					"widget": maps.Map{
						"id":      widget.Id,
						"author":  widget.Author,
						"version": widget.Version,
					},
				})
			}
		}
	}
	this.Data["usingCharts"] = usingCharts

	// 所有的图表
	this.Data["widgets"] = lists.Map(widgets.LoadAllWidgets(), func(k int, v interface{}) interface{} {
		widget := v.(*widgets.Widget)

		return maps.Map{
			"id": widget.Id,
			"charts": lists.Map(widget.Charts, func(k int, v interface{}) interface{} {
				chart := v.(*widgets.Chart)
				isUsing := false
				if params.BoardType == "realtime" {
					if server.RealtimeBoard != nil {
						isUsing = server.RealtimeBoard.HasChart(chart.Id)
					}
				} else {
					if server.StatBoard != nil {
						isUsing = server.StatBoard.HasChart(chart.Id)
					}
				}

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

				return maps.Map{
					"id":          chart.Id,
					"name":        chart.Name,
					"description": chart.Description,
					"items":       items,
					"columns":     chart.Columns,
					"on":          chart.On,
					"isUsing":     isUsing,
					"widget": maps.Map{
						"id":      widget.Id,
						"author":  widget.Author,
						"version": widget.Version,
					},
				}
			}),
		}
	})

	this.Show()
}
