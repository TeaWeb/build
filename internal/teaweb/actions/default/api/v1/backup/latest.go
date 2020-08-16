package backup

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/settings/backup/backuputils"
	"github.com/iwind/TeaGo/actions"
)

type LatestAction actions.Action

// 下载最新的备份文件
func (this *LatestAction) RunGet(params struct{}) {
	backupFiles := backuputils.ActionListFiles()
	if len(backupFiles) == 0 {
		apiutils.Fail(this, "no backup files")
		return
	}

	file := backupFiles[0]
	backuputils.ActionDownloadFile(file.GetString("name"), this.ResponseWriter, func() {
		apiutils.Fail(this, "file not found")
	})
}
