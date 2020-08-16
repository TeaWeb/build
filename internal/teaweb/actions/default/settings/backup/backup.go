package backup

import "github.com/iwind/TeaGo/actions"

type BackupAction actions.Action

// 立即备份
func (this *BackupAction) Run(params struct{}) {
	err := backupTask()
	if err != nil {
		this.Fail("备份失败：" + err.Error())
	}
	this.Success()
}
