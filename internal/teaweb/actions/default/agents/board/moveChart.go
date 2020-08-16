package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/iwind/TeaGo/actions"
)

type MoveChartAction actions.Action

// 移动图表位置
func (this *MoveChartAction) Run(params struct {
	AgentId  string
	OldIndex int
	NewIndex int
}) {
	board := agents.NewAgentBoard(params.AgentId)
	if board == nil {
		this.Fail("找不到Board")
	}

	board.MoveChart(params.OldIndex, params.NewIndex)
	err := board.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
