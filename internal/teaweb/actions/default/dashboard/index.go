package dashboard

import (
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

// 仪表板
func (this *IndexAction) Run(params struct{}) {
	this.Data["teaMenu"] = "dashboard"
	this.Show()
}

