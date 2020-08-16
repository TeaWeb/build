package backup

import (
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/settings/backup/backuputils"
	"github.com/iwind/TeaGo/actions"
)

type DeleteAction actions.Action

// 删除备份
func (this *DeleteAction) Run(params struct {
	File string
}) {
	if teaconst.DemoEnabled {
		this.Fail("演示版无法删除")
	}

	if len(params.File) == 0 {
		this.Fail("请指定要删除的备份文件")
	}

	if !backuputils.ActionDeleteFile(params.File, func(err error) {
		this.Fail("删除失败：" + err.Error())
	}) {
		return
	}

	this.Success()
}
