package backup

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"regexp"
	"time"
)

type CleanAction actions.Action

// 清除N天以前的备份文件
func (this *CleanAction) RunPost(params struct{}) {
	oldDay := timeutil.Format("Ymd", time.Now().AddDate(0, 0, -30)) + ".zip"

	reg := regexp.MustCompile(`^\d{8}\.zip$`)
	dir := files.NewFile(Tea.Root + "/backups/")
	var lastErr error = nil
	count := 0
	if dir.Exists() {
		for _, f := range dir.List() {
			if !reg.MatchString(f.Name()) {
				continue
			}
			if f.Name() < oldDay {
				err := f.Delete()
				if err != nil {
					lastErr = err
					break
				}
				count++
			}
		}
	}

	this.Data["count"] = count
	if lastErr != nil {
		this.Fail("清除失败：" + lastErr.Error())
	} else {
		this.Success()
	}
}
