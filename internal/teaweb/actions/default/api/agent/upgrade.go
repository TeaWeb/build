package agent

import (
	"github.com/TeaWeb/build/internal/teaagent/agentconst"
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
)

type UpgradeAction actions.Action

// 升级信息
func (this *UpgradeAction) Run(params struct{}) {
	agent := this.Context.Get("agent").(*agents.AgentConfig)
	if !agent.On {
		this.Fail("agent is not on")
	}

	agentId := agent.Id
	agentVersion := this.Request.Header.Get("Tea-Agent-Version")
	agentOS := this.Request.Header.Get("Tea-Agent-Os")
	agentArch := this.Request.Header.Get("Tea-Agent-Arch")

	// 是否需要更新
	if agentId == "local" || len(agentVersion) == 0 || len(agentOS) == 0 || len(agentArch) == 0 {
		this.Fail("agentVersion, agentOS, agentArch should not be empty")
	}

	if agentVersion == agentconst.AgentVersion {
		this.Fail("agent is latest version")
	}

	agentFile := teautils.WebRoot() + "upgrade" + Tea.DS + agentconst.AgentVersion + Tea.DS + agentOS + Tea.DS + agentArch + Tea.DS
	if agentOS == "windows" {
		agentFile += "teaweb-agent.exe"
	} else {
		agentFile += "teaweb-agent"
	}

	file := files.NewFile(agentFile)
	if !file.Exists() {
		this.Fail("no upgrade file")
	}

	data, err := file.ReadAll()
	if err != nil {
		logs.Error(err)
		this.Fail(err.Error())
	}

	this.AddHeader("Tea-Agent-Version", agentconst.AgentVersion)
	this.Write(data)
}
