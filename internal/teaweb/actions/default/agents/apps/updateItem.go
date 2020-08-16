package apps

import (
	"encoding/json"
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaconfigs/forms"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"net/http"
	"regexp"
)

type UpdateItemAction actions.Action

// 添加监控项
func (this *UpdateItemAction) Run(params struct {
	AgentId string
	AppId   string
	ItemId  string
	From    string
}) {
	app := agentutils.InitAppData(this, params.AgentId, params.AppId, "monitor")

	item := app.FindItem(params.ItemId)
	if item == nil {
		this.Fail("找不到要修改的监控项")
	}
	if item.RecoverSuccesses <= 0 {
		item.RecoverSuccesses = 1
	}
	this.Data["item"] = item

	this.Data["from"] = params.From

	// 数据源
	hasPluginSource := false
	this.Data["sources"] = lists.Map(agents.AllDataSources(), func(k int, v interface{}) interface{} {
		m := v.(maps.Map)
		instance := m["instance"].(agents.SourceInterface)
		m["variables"] = instance.Variables()
		m["thresholds"] = instance.Thresholds()
		m["platforms"] = instance.Platforms()

		if m["category"] == agents.SourceCategoryPlugin {
			hasPluginSource = true
		}
		return m
	})
	this.Data["hasPluginSource"] = hasPluginSource

	this.Data["methods"] = []string{http.MethodGet, http.MethodPost, http.MethodPut}
	this.Data["dataFormats"] = agents.AllSourceDataFormats()
	this.Data["operators"] = agents.AllThresholdOperators()
	this.Data["noticeLevels"] = notices.AllNoticeLevels()
	this.Data["actions"] = agents.AllActions()

	groups1 := []*forms.Group{}
	groups2 := []*forms.Group{}
	css := ""
	javascript := ""

	for _, source := range agents.AllDataSources() {
		sourceInstance := source["instance"].(agents.SourceInterface)
		form := sourceInstance.Form()
		if form == nil {
			continue
		}
		form.ComposedAttrs = map[string]string{
			"v-show": "sourceCode == '" + sourceInstance.Code() + "'",
		}
		if sourceInstance.Code() == item.SourceCode {
			form.Init(item.SourceOptions)
		}
		form.Compose()

		css += form.CSS
		javascript += form.Javascript

		countGroups := len(form.Groups)
		if countGroups == 0 {
			continue
		} else if countGroups == 1 {
			groups1 = append(groups1, form.Groups[0])
		} else {
			groups1 = append(groups1, form.Groups[0])
			for i := 1; i < countGroups; i++ {
				groups2 = append(groups2, form.Groups[i])
			}
		}
	}

	this.Data["formGroups1"] = groups1
	this.Data["formGroups2"] = groups2
	this.Data["formCSS"] = css
	this.Data["formJavascript"] = javascript

	this.Show()
}

// 提交保存
func (this *UpdateItemAction) RunPost(params struct {
	AgentId string
	AppId   string
	ItemId  string

	Name       string
	SourceCode string
	On         bool

	DataFormat uint8
	Interval   uint

	CondParams         []string
	CondOps            []string
	CondValues         []string
	CondNoticeLevels   []uint
	CondNoticeMessages []string
	CondActions        []string
	CondMaxFails       []int

	RecoverSuccesses int

	Must *actions.Must
}) {
	if teaconst.DemoEnabled {
		this.Fail("演示版无法修改监控项")
	}

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

	params.Must.
		Field("name", params.Name).
		Require("请输入监控项名称").
		Field("sourceCode", params.SourceCode).
		Require("请选择数据源类型")

	item.On = params.On
	item.Name = params.Name

	// 数据源
	item.SourceCode = params.SourceCode
	item.SourceOptions = map[string]interface{}{}

	// 获取参数值
	instance := agents.FindDataSourceInstance(params.SourceCode, map[string]interface{}{})
	form := instance.Form()
	values, errField, err := form.ApplyRequest(this.Request)
	if err != nil {
		this.FailField(errField, err.Error())
	}

	values["dataFormat"] = params.DataFormat
	item.SourceOptions = values

	// 刷新间隔等其他选项
	item.Interval = fmt.Sprintf("%ds", params.Interval)
	item.RecoverSuccesses = params.RecoverSuccesses

	// 阈值设置
	item.Thresholds = []*agents.Threshold{}
	for index, param := range params.CondParams {
		if index < len(params.CondValues) &&
			index < len(params.CondOps) &&
			index < len(params.CondValues) &&
			index < len(params.CondNoticeLevels) &&
			index < len(params.CondNoticeMessages) &&
			index < len(params.CondMaxFails) {
			// 校验
			op := params.CondOps[index]
			value := params.CondValues[index]
			if op == agents.ThresholdOperatorRegexp || op == agents.ThresholdOperatorNotRegexp {
				_, err := regexp.Compile(value)
				if err != nil {
					this.Fail("阈值" + param + "正则表达式" + value + "校验失败：" + err.Error())
				}
			}

			// 动作
			actionJSON := params.CondActions[index]
			actionList := []map[string]interface{}{}
			err := json.Unmarshal([]byte(actionJSON), &actionList)
			if err != nil {
				logs.Error(err)
			}

			t := agents.NewThreshold()
			t.Param = param
			t.Operator = op
			t.Value = value
			t.NoticeLevel = types.Uint8(params.CondNoticeLevels[index])
			t.NoticeMessage = params.CondNoticeMessages[index]
			t.Actions = actionList
			t.MaxFails = params.CondMaxFails[index]
			item.AddThreshold(t)
		}
	}

	err = agent.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 通知更新
	agentutils.PostAgentEvent(agent.Id, agentutils.NewAgentEvent("UPDATE_ITEM", maps.Map{
		"appId":  app.Id,
		"itemId": params.ItemId,
	}))

	if app.IsSharedWithGroup {
		err = agentutils.SyncApp(agent.Id, agent.GroupIds, app, agentutils.NewAgentEvent("UPDATE_ITEM", maps.Map{
			"appId":  app.Id,
			"itemId": params.ItemId,
		}), nil)
		if err != nil {
			logs.Error(err)
		}
	}

	this.Success()
}
