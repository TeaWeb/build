package mongo

import (
	"github.com/iwind/TeaGo/actions"
)

type DataAction actions.Action

// 数据清理
func (this *DataAction) Run(params struct{}) {
	this.Show()
}
