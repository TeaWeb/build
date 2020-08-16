package mysql

import (
	"github.com/TeaWeb/build/internal/teaconfigs/db"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/iwind/TeaGo/actions"
	"strings"
)

type IndexAction actions.Action

// MySQL连接信息
func (this *IndexAction) RunGet(params struct{}) {
	this.Data["shouldRestart"] = shouldRestart
	this.Data["error"] = ""

	config, err := db.LoadMySQLConfig()
	if err != nil {
		this.Data["error"] = err.Error()
		config = db.DefaultMySQLConfig()
	} else {
		err = teadb.SharedDB().Test()
		if err != nil {
			this.Data["error"] = err.Error()
		}
	}
	config.Password = strings.Repeat("*", len(config.Password))

	this.Data["config"] = config

	this.Show()
}
