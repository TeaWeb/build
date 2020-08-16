package board

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teastats"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

type ItemsAction actions.Action

// 数据指标
func (this *ItemsAction) RunGet(params struct {
	BoardType string
	ServerId  string
}) {
	this.Data["boardType"] = params.BoardType

	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	this.Data["server"] = maps.Map{
		"id": server.Id,
	}

	// 指标列表
	items := []interface{}{}
	for _, filter := range teastats.FindAllStatFilters() {
		if lists.ContainsString(server.StatItems, filter.GetString("code")) {
			continue
		}
		items = append(items, filter)
	}
	this.Data["items"] = items

	// 自己添加的运行的指标
	runningItems := []maps.Map{}
	for _, code := range server.StatItems {
		filter := teastats.FindSharedFilter(code)
		if filter == nil {
			continue
		}
		m := maps.Map{
			"code":        code,
			"name":        filter.Name(),
			"description": filter.Description(),
			"params":      filter.ParamVariables(),
			"values":      filter.ValueVariables(),
			"charts":      []maps.Map{},
		}
		runningItems = append(runningItems, m)
	}
	lists.Sort(runningItems, func(i int, j int) bool {
		m1 := runningItems[i]
		m2 := runningItems[j]
		code1 := m1.GetString("code")
		code2 := m2.GetString("code")

		return code1 < code2
	})
	this.Data["runningItems"] = runningItems

	// 图表引用的指标
	chartItemsMap := maps.Map{}
	for _, board := range []*teaconfigs.Board{server.RealtimeBoard, server.StatBoard} {
		if board == nil {
			continue
		}
		for _, c := range board.Charts {
			_, chart := c.FindChart()
			if chart == nil || !chart.On {
				continue
			}
			for _, code := range chart.Requirements {
				_, ok := chartItemsMap[code]
				if !ok {
					filter := teastats.FindSharedFilter(code)
					name := ""
					description := ""
					params := []*teastats.Variable{}
					values := []*teastats.Variable{}
					if filter != nil {
						name = filter.Name()
						description = filter.Description()

						index := strings.LastIndex(code, ".")
						if index > 0 {
							period := teastats.FindPeriodName(code[index+1:])
							if len(period) > 0 {
								name += "（" + period + "）"
							}
						}

						params = filter.ParamVariables()
						values = filter.ValueVariables()
					}

					chartItemsMap[code] = maps.Map{
						"code":        code,
						"name":        name,
						"description": description,
						"params":      params,
						"values":      values,
						"charts":      []maps.Map{},
					}
				}

				boardType := ""
				switch board {
				case server.RealtimeBoard:
					boardType = "realtime"
				case server.StatBoard:
					boardType = "stat"
				}

				chartItemsMap[code].(maps.Map)["charts"] = append(chartItemsMap[code].(maps.Map)["charts"].([]maps.Map), maps.Map{
					"name":      chart.Name,
					"widgetId":  c.WidgetId,
					"chartId":   c.ChartId,
					"boardType": boardType,
				})
			}
		}
	}

	chartItems := chartItemsMap.Values()
	lists.Sort(chartItems, func(i int, j int) bool {
		m1 := chartItems[i]
		m2 := chartItems[j]
		code1 := m1.(maps.Map).GetString("code")
		code2 := m2.(maps.Map).GetString("code")

		return code1 < code2
	})

	this.Data["chartItems"] = chartItems

	this.Show()
}
