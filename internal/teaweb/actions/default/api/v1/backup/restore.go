package backup

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/settings/backup/backuputils"
	"github.com/iwind/TeaGo/actions"
)

type RestoreAction actions.Action

// 恢复
func (this *RestoreAction) RunGet(params struct {
	Filename string
}) {
	if !backuputils.ActionRestoreFile(params.Filename, func(message string) {
		apiutils.Fail(this, message)
	}) {
		return
	}
	apiutils.SuccessOK(this)
}
