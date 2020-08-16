package database

import (
	"github.com/TeaWeb/build/internal/teaconfigs/db"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

// 数据库设置首页
func (this *IndexAction) RunGet(params struct{}) {
	config := db.SharedDBConfig()
	this.Data["dbType"] = config.Type
	this.Data["dbTypeName"] = config.TypeName()

	this.Show()
}
