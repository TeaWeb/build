package board

import (
	"encoding/json"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teastats"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/board/scripts"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type TestAction actions.Action

// 测试图表
func (this *TestAction) Run(params struct {
	ServerId       string
	Name           string
	Description    string
	Columns        uint8
	Items          []string
	JavascriptCode string
	On             bool
	AutoGenerate   bool // 是否自动生成测试数据
	Events         string
	Must           *actions.Must
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	// 如果选中了指标，而且指标没有数据的话，则试着根据access log生成指标数据
	if params.AutoGenerate {
		if len(params.Items) > 0 {
			queue := teastats.NewQueue()
			queue.Start(server.Id)
			for _, item := range params.Items {
				instance := teastats.FindNewFilter(item)
				if instance == nil {
					continue
				}

				// 是否有数据
				v, err := teadb.ServerValueDAO().FindOneWithItem(server.Id, item)
				if err != nil {
					logs.Error(err)
				} else if v == nil {
					// 读取日志
					ones, err := teadb.AccessLogDAO().ListLatestAccessLogs(timeutil.Format("Ymd"), server.Id, "", false, 1000)
					if err != nil {
						logs.Error(err)
					} else {
						instance.Start(queue, item)
						for _, one := range ones {
							instance.Filter(one)
						}
						instance.Stop()
					}
				}
			}
			queue.Stop()
		}
	}

	// 图表信息
	chart := widgets.NewChart()
	chart.Id = "test_chart"
	chart.On = params.On
	chart.Name = params.Name
	chart.Description = params.Description
	chart.Columns = params.Columns
	chart.Requirements = params.Items
	chart.Type = "javascript"
	chart.Options = maps.Map{
		"code": params.JavascriptCode,
	}
	obj, err := chart.AsObject()
	if err != nil {
		this.Fail("运行错误：" + err.Error())
	}

	// 事件
	events := []interface{}{}
	if len(params.Events) > 0 {
		err := json.Unmarshal([]byte(params.Events), &events)
		if err != nil {
			logs.Error(err)
		}
	}

	code, err := obj.AsJavascript(map[string]interface{}{
		"name":    params.Name,
		"columns": params.Columns,
		"id":      chart.Id,
		"events":  events,
	})
	if err != nil {
		this.Fail("运行错误：" + err.Error())
	}

	engine := scripts.NewEngine()
	engine.SetMongo(teadb.SharedDB().Test() == nil)
	engine.SetContext(&scripts.Context{
		Server: server,
	})

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

	this.Data["charts"] = engine.Charts()
	this.Data["output"] = engine.Output()

	this.Success()
}
