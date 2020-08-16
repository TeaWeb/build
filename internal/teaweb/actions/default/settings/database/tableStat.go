package database

import (
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/iwind/TeaGo/actions"
)

type TableStatAction actions.Action

// 集合统计
func (this *TableStatAction) Run(params struct {
	Tables []string
}) {
	statMap, err := teadb.SharedDB().StatTables(params.Tables)
	if err != nil {
		this.Data["result"] = map[string]interface{}{}
		this.Fail("获取统计信息失败：" + err.Error())
	} else {
		this.Data["result"] = statMap
	}

	this.Success()
}
