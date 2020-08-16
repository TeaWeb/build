package database

import (
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/iwind/TeaGo/actions"
)

type DeleteTableAction actions.Action

// 删除集合
func (this *DeleteTableAction) Run(params struct {
	Table string
}) {
	if len(params.Table) > 0 {
		err := teadb.SharedDB().DropTable(params.Table)
		if err != nil {
			this.Fail("删除失败：" + err.Error())
		}
	}
	this.Success()
}
