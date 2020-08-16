package backup

import (
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/settings/backup/backuputils"
	"github.com/iwind/TeaGo/actions"
)

type DownloadAction actions.Action

// 下载
func (this *DownloadAction) Run(params struct {
	Filename string
}) {
	if teaconst.DemoEnabled {
		this.Fail("演示版无法下载")
	}

	backuputils.ActionDownloadFile(params.Filename, this.ResponseWriter, func() {
		this.WriteString("file not found")
	})
}
