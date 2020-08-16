package agent

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

// 读取Agent配置
func (this *IndexAction) Run(params struct{}) {
	agent := this.Context.Get("agent")
	if agent == nil {
		this.Fail("not found agent")
	}

	agentConfig, ok := agent.(*agents.AgentConfig)
	if !ok {
		this.Fail("invalid agent")
	}

	data, err := agentConfig.EncodeYAML()
	if err != nil {
		this.Fail("YAML encode error：" + err.Error())
	}

	this.Data["config"] = string(data)

	this.Success()
}
