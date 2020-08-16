package settings

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
)

type UpdateAction actions.Action

// 修改设置
func (this *UpdateAction) Run(params struct {
	AgentId string
}) {
	this.Data["selectedTab"] = "detail"

	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到要修改的Agent")
	}

	this.Data["agent"] = agent
	this.Data["groups"] = agents.SharedGroupList().Groups
	if len(agent.GroupIds) > 0 {
		this.Data["groupId"] = agent.GroupIds[0]
	} else {
		this.Data["groupId"] = "default"
	}

	this.Show()
}

// 提交保存
func (this *UpdateAction) RunPost(params struct {
	AgentId             string
	Name                string
	Host                string
	GroupId             string
	AllowAllIP          bool
	IPs                 []string `alias:"ips"`
	On                  bool
	Key                 string
	CheckDisconnections bool
	AutoUpdates         bool
	Must                *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入主机名").
		Field("host", params.Host).
		Require("请输入主机地址").
		Field("key", params.Key).
		Require("请输入密钥")

	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到要修改的Agent")
	}
	agent.On = params.On
	agent.Name = params.Name
	agent.Host = params.Host
	if len(params.GroupId) == 0 {
		agent.GroupIds = []string{}
	} else {
		agent.GroupIds = []string{params.GroupId}
	}
	agent.AllowAll = params.AllowAllIP
	agent.Allow = params.IPs
	agent.Key = params.Key
	agent.CheckDisconnections = params.CheckDisconnections
	agent.AutoUpdates = params.AutoUpdates
	err := agent.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 重建索引
	err = agents.SharedGroupList().BuildIndexes()
	if err != nil {
		logs.Error(err)
	}

	// 通知更新
	agentutils.PostAgentEvent(agent.Id, agentutils.NewAgentEvent("UPDATE_AGENT", maps.Map{}))

	this.Success()
}
