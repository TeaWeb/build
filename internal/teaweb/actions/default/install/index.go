package install

import (
	"github.com/iwind/TeaGo/actions"
	"runtime"
)

type IndexAction actions.Action

// 安装
func (this *IndexAction) RunGet(params struct {
}) {
	this.Data["os"] = runtime.GOOS
	this.Show()
}
