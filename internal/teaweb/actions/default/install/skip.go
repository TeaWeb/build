package install

import (
	"github.com/TeaWeb/build/internal/teaconfigs/db"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
)

type SkipAction actions.Action

// 跳过
func (this *SkipAction) RunPost(params struct{}) {
	config := db.SharedDBConfig()
	config.IsInitialized = true
	err := config.Save()
	if err != nil {
		logs.Error(err)
	}

	this.Success()
}
