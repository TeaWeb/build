package board

import (
	"encoding/json"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/board/scripts"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

// 看板
func (this *IndexAction) Run(params struct {
	ServerId  string
	BoardType string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到要查看的代理服务")
	}

	if !server.IsHTTP() {
		this.RedirectURL("/proxy/detail?serverId=" + server.Id)
		return
	}

	this.Data["server"] = maps.Map{
		"id": server.Id,
	}

	if len(params.BoardType) == 0 {
		params.BoardType = "realtime"
	}
	this.Data["boardType"] = params.BoardType

	this.Data["errs"] = teaproxy.SharedManager.FindServerErrors(params.ServerId)

	this.Show()
}

// 看板数据
func (this *IndexAction) RunPost(params struct {
	ServerId string
	Type     string // realtime or stat
	Events   string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到要查看的代理服务")
	}

	if len(params.Type) == 0 {
		params.Type = "realtime"
	}

	var board *teaconfigs.Board
	switch params.Type {
	case "realtime":
		board = server.RealtimeBoard
	case "stat":
		board = server.StatBoard
	}

	// 初始化
	shouldSave := false
	if server.RealtimeBoard == nil {
		shouldSave = true

		board = teaconfigs.NewBoard()
		server.RealtimeBoard = board

		// 添加一些默认的图表
		board.AddChart("teaweb.proxy_status", "kTVuOEBm605H3AJS")
		board.AddChart("teaweb.locations", "cyZsvwR66oxcVpvj")
		board.AddChart("teaweb.bandwidth_realtime", "hiAUsteL1V6LG8zD")
		board.AddChart("teaweb.request_realtime", "APvSaVEoQ7VUvX4a")
		board.AddChart("teaweb.request_time", "g8SrxuMYwwhxNWwk")
		board.AddChart("teaweb.status_stat", "xnUsgQMSjWZ9MN7g")
		board.AddChart("teaweb.latest_errors", "RUCF1EbF4FpPMHpN")
	}
	if server.StatBoard == nil {
		shouldSave = true

		board = teaconfigs.NewBoard()
		server.StatBoard = board

		// 添加一些默认的图表
		board.AddChart("teaweb.request_stat", "QdC9mdfECEqdKllp")
		board.AddChart("teaweb.url_rank", "3TLvbynJiU6Ik5Yh")
		board.AddChart("teaweb.cost_rank", "U79hACpwPmMtmEgu")
		board.AddChart("teaweb.os_rank", "rJLuNQm6UeGJbyz6")
		board.AddChart("teaweb.browser_rank", "2yS1FKYOv0nIMc4I")
		board.AddChart("teaweb.region_rank", "nDx94UkqwzrdaNg1")
		board.AddChart("teaweb.province_rank", "chtZBWnb955NCre7")
	}

	if shouldSave {
		err := server.Save()
		if err != nil {
			logs.Error(err)
		}

		// 重启统计服务
		proxyutils.ReloadServerStats(server.Id)
	}

	if len(board.Charts) == 0 {
		this.Data["charts"] = []maps.Map{}
		this.Success()
	}

	engine := scripts.NewEngine()
	engine.SetMongo(teadb.SharedDB().Test() == nil)
	engine.SetContext(&scripts.Context{
		Server: server,
	})

	// 事件
	events := []interface{}{}
	if len(params.Events) > 0 {
		err := json.Unmarshal([]byte(params.Events), &events)
		if err != nil {
			logs.Error(err)
		}
	}

	for _, c := range board.Charts {
		_, chart := c.FindChart()
		if chart == nil || !chart.On {
			continue
		}

		obj, err := chart.AsObject()
		if err != nil {
			this.Fail(err.Error())
		}
		code, err := obj.AsJavascript(map[string]interface{}{
			"name":    chart.Name,
			"columns": chart.Columns,
			"id":      chart.Id,
			"events":  events,
		})
		if err != nil {
			this.Fail(err.Error())
		}

		widgetCode := `var widget = new widgets.Widget({
	"name": "看板",
	"requirements": ["mongo"]
});

widget.run = function () {
`
		widgetCode += "{\n" + code + "\n}\n"
		widgetCode += `
};
`

		err = engine.RunCode(widgetCode)
		if err != nil {
			this.Fail("运行错误：" + err.Error())
		}
	}

	this.Data["charts"] = engine.Charts()
	this.Data["output"] = engine.Output()

	this.Success()
}
