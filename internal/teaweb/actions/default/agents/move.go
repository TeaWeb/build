package agents

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/iwind/TeaGo/actions"
)

type MoveAction actions.Action

// 交换位置
func (this *MoveAction) Run(params struct {
	FromId string
	ToId   string
}) {
	config, err := agents.SharedAgentList()
	if err != nil {
		this.Fail("保存错误：" + err.Error())
	}
	config.MoveAgent(params.FromId, params.ToId)
	err = config.Save()
	if err != nil {
		this.Fail("保存错误：" + err.Error())
	}

	this.Success()
}
