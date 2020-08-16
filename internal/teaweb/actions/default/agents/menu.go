package agents

import "github.com/iwind/TeaGo/actions"

type MenuAction actions.Action

// 菜单
func (this *MenuAction) Run(params struct{}) {
	this.Success()
}
