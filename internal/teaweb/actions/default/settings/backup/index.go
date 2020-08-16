package backup

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/settings/backup/backuputils"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

// 备份列表
func (this *IndexAction) Run(params struct{}) {
	this.Data["files"] = backuputils.ActionListFiles()

	this.Data["shouldRestart"] = backuputils.ShouldRestart()

	this.Show()
}
