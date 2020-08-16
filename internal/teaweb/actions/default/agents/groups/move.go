package groups

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/iwind/TeaGo/actions"
)

type MoveAction actions.Action

// 移动分组位置
func (this *MoveAction) Run(params struct {
	FromIndex int
	ToIndex   int
}) {
	config := agents.SharedGroupList()
	config.Move(params.FromIndex, params.ToIndex)
	err := config.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
