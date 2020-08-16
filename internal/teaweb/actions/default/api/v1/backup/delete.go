package backup

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/settings/backup/backuputils"
	"github.com/iwind/TeaGo/actions"
)

type DeleteAction actions.Action

// 删除备份文件
func (this *DeleteAction) RunGet(params struct {
	Filename string
}) {
	if len(params.Filename) == 0 {
		apiutils.Fail(this, "'filename' should not be empty")
		return
	}

	if !backuputils.ActionDeleteFile(params.Filename, func(err error) {
		apiutils.Fail(this, err.Error())
	}) {
		return
	}

	apiutils.SuccessOK(this)
}
