package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/string"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"strings"
	"time"
)

type ItemValuesAction actions.Action

// 监控项数据展示
func (this *ItemValuesAction) Run(params struct {
	AgentId string
	AppId   string
	ItemId  string
	Level   int
}) {
	app := agentutils.InitAppData(this, params.AgentId, params.AppId, "monitor")
	item := app.FindItem(params.ItemId)

	if item == nil {
		this.Fail("找不到要查看的Item")
	}

	this.Data["item"] = item
	this.Data["levels"] = notices.AllNoticeLevels()
	this.Data["selectedLevel"] = params.Level

	this.Show()
}

// 获取监控项数据
func (this *ItemValuesAction) RunPost(params struct {
	AgentId string
	AppId   string
	ItemId  string
	LastId  string
	Level   notices.NoticeLevel
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		this.Fail("找不到App")
	}

	item := app.FindItem(params.ItemId)
	if item == nil {
		this.Fail("找不到Item")
	}

	ones, err := teadb.AgentValueDAO().FindLatestItemValues(params.AgentId, params.AppId, params.ItemId, params.Level, params.LastId, 60)
	if err != nil {
		this.Fail("查询失败：" + err.Error())
	}

	source := item.Source()
	this.Data["values"] = lists.Map(ones, func(k int, v interface{}) interface{} {
		value := v.(*agents.Value)

		vars := []maps.Map{}
		if types.IsMap(value.Value) || types.IsSlice(value.Value) {
			if source != nil {
				for _, variable := range source.Variables() {
					if len(variable.Code) == 0 || strings.Index(variable.Code, "$") > -1 {
						continue
					}
					result := teautils.Get(value.Value, strings.Split(variable.Code, "."))
					vars = append(vars, maps.Map{
						"code":        variable.Code,
						"description": variable.Description,
						"value":       stringutil.JSONEncodePretty(result),
					})
				}
			}
		}

		return maps.Map{
			"id":          value.Id.Hex(),
			"costMs":      value.CostMs,
			"value":       value.Value,
			"error":       value.Error,
			"noticeLevel": notices.FindNoticeLevel(value.NoticeLevel),
			"threshold":   value.Threshold,
			"vars":        vars,
			"beginTime":   timeutil.Format("Y-m-d H:i:s", time.Unix(value.CreatedAt, 0)),
			"endTime":     timeutil.Format("Y-m-d H:i:s", time.Unix(value.Timestamp, 0)),
		}
	})
	this.Success()
}
