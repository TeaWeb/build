package agents

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
)

type AddAgentAction actions.Action

// 添加代理
func (this *AddAgentAction) Run(params struct{}) {
	this.Data["defaultGroupName"] = agents.SharedGroupList().FindDefaultGroup().Name
	this.Data["groups"] = agents.SharedGroupList().Groups

	this.Show()
}

// 提价保存
func (this *AddAgentAction) RunPost(params struct {
	Name                string
	Host                string
	GroupId             string
	AllowAllIP          bool
	IPs                 []string `alias:"ips"`
	On                  bool
	CheckDisconnections bool
	AutoUpdates         bool
	Must                *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入主机名").
		Field("host", params.Host).
		Require("请输入主机地址")

	agentList, err := agents.SharedAgentList()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	group := agents.SharedGroupList().FindGroup(params.GroupId)
	if group != nil {
		if group.MaxAgents > 0 && group.CountAgents >= group.MaxAgents {
			this.Fail("选择的分组不能超过最大Agent数量限制")
		}
		if !group.IsDateAvailable() {
			this.Fail(" 选择的分组不在有效期限内")
		}
	}

	agent := agents.NewAgentConfig()
	agent.On = params.On
	agent.Name = params.Name
	agent.Host = params.Host
	if len(params.GroupId) > 0 {
		agent.AddGroup(params.GroupId)
	}
	agent.AllowAll = params.AllowAllIP
	agent.Allow = params.IPs
	agent.Key = rands.HexString(32)
	agent.CheckDisconnections = params.CheckDisconnections
	agent.AutoUpdates = params.AutoUpdates
	agent.AddDefaultApps()
	err = agent.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	agentList.AddAgent(agent.Filename())
	err = agentList.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 重建索引
	err = agents.SharedGroupList().BuildIndexes()
	if err != nil {
		logs.Error(err)
	}

	this.Data["agentId"] = agent.Id

	// 通知更新
	agentutils.PostAgentEvent(agent.Id, agentutils.NewAgentEvent("ADD_AGENT", maps.Map{}))

	this.Success()
}
