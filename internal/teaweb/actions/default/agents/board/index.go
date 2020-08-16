package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/board/scripts"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

// 看板首页
func (this *IndexAction) Run(params struct {
	AgentId string
}) {
	if len(params.AgentId) == 0 {
		params.AgentId = "local"
	}

	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	// 检查版本更新
	this.checkUpgrade(agent)

	this.Data["agentId"] = params.AgentId
	this.Data["tabbar"] = "board"

	this.Show()
}

// 数据
func (this *IndexAction) RunPost(params struct {
	AgentId string
}) {
	this.Data["error"] = ""

	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	// 添加默认App
	if agent != nil && !agent.AppsIsInitialized {
		agent.AddDefaultApps()
		err := agent.Save()

		// 通知更新
		if err == nil {
			agentutils.PostAgentEvent(agent.Id, agentutils.NewAgentEvent("UPDATE_AGENT", maps.Map{}))
		}
	}

	board := agents.NewAgentBoard(params.AgentId)
	if board == nil {
		this.Fail("无法读取Board配置")
	}

	if !teadb.SharedDB().IsAvailable() {
		this.Data["charts" ] = []interface{}{}
		this.Data["output"] = []string{}
		this.Data["error"] = "当前数据库不可用，无法展示图表"
		this.Success()
	}

	dbEnabled := teadb.SharedDB().Test() == nil
	if !dbEnabled {
		this.Data["charts" ] = []interface{}{}
		this.Data["output"] = []string{}
		this.Data["error"] = "当前数据库无法连接，无法展示图表"
		this.Success()
	}

	engine := scripts.NewEngine()
	engine.SetDBEnabled(dbEnabled)

	for _, c := range board.Charts {
		app := agent.FindApp(c.AppId)
		if app == nil || !app.On {
			continue
		}

		item := app.FindItem(c.ItemId)
		if item == nil || !item.On {
			continue
		}

		chart := item.FindChart(c.ChartId)
		if chart == nil || !chart.On {
			continue
		}

		// 设置
		if len(c.Name) > 0 {
			chart.Name = c.Name
		}

		o, err := chart.AsObject()
		if err != nil {
			logs.Error(err)
			continue
		}

		var chartName = chart.Name + "<span class=\"ops\">"
		if chart.SupportsTimeRange {
			chartName += "<a href=\"/agents/board/chart?agentId=" + agent.Id + "&appId=" + c.AppId + "&itemId=" + c.ItemId + "&chartId=" + c.ChartId + "\" title=\"更多时间选择\"><i class=\"icon clock small\"></i></a> &nbsp; "
		}
		chartName += "<a href=\"/agents/apps/itemValues?agentId=" + agent.Id + "&appId=" + app.Id + "&itemId=" + item.Id + "\" title=\"查看数值记录\"><i class=\"icon external small\"></i></a> &nbsp; <a href=\"\" title=\"从看板移除\" onclick=\"return Tea.Vue.removeChart('" + c.AppId + "', '" + c.ItemId + "', '" + c.ChartId + "')\"><i class=\"icon remove small\"></i></a></span>"
		code, err := o.AsJavascript(maps.Map{
			"name":    chartName,
			"columns": chart.Columns,
		})
		if err != nil {
			logs.Error(err)
			continue
		}

		ctx := &scripts.Context{
			Agent:    agent,
			App:      app,
			Item:     item,
			TimeType: c.TimeType,
			TimePast: c.TimePast,
			DayFrom:  c.DayFrom,
			DayTo:    c.DayTo,
		}
		engine.SetContext(ctx)

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
			if err != teadb.ErrorDBUnavailable {
				logs.Error(err)
			}
			continue
		}
	}

	this.Data["charts" ] = engine.Charts()
	this.Data["output"] = engine.Output()
	this.Success()
}

// 检查版本更新
// deprecated in v0.1.8
func (this *IndexAction) checkUpgrade(agent *agents.AgentConfig) {
	if len(agent.TeaVersion) == 0 { // 0.1.7之前
		isChanged := false
		for _, app := range agent.Apps {
			for _, item := range app.Items {
				err := item.Validate()
				if err != nil {
					logs.Error(err)
					continue
				}
				source := item.Source()
				if source == nil {
					continue
				}

				for _, oldChart := range item.Charts {
					for _, newChart := range source.Charts() {
						if oldChart.Id == newChart.Id {
							isChanged = true
							oldChart.SupportsTimeRange = newChart.SupportsTimeRange
							oldChart.Options = newChart.Options
							break
						}
					}
				}
			}
		}

		if isChanged {
			err := agent.Save()
			if err != nil {
				logs.Error(err)
			}
		}

		return
	}
}
