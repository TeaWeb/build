package backup

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/settings/backup/backuputils"
	"github.com/iwind/TeaGo/actions"
)

type FilesAction actions.Action

// 文件名
func (this *FilesAction) RunGet(params struct{}) {
	apiutils.Success(this, backuputils.ActionListFiles())
}
