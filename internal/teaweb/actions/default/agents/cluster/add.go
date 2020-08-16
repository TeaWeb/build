package cluster

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
)

type AddAction actions.Action

// 批量添加
func (this *AddAction) Run(params struct{}) {
	// 检查安装工具
	{
		dirFile := files.NewFile(Tea.Root + "/web/installers")
		if dirFile.Exists() && len(dirFile.List()) > 0 {
			this.Data["checkInstaller"] = true
		} else {
			this.Data["checkInstaller"] = false
		}
	}

	//  Agent新版本
	{
		dirFile := files.NewFile(Tea.Root + "/web/upgrade")
		if dirFile.Exists() && len(dirFile.List()) > 0 {
			this.Data["checkUpgradeFiles"] = true
		} else {
			this.Data["checkUpgradeFiles"] = false
		}
	}

	// 分组信息
	this.Data["groups"] = agents.SharedGroupList().Groups

	this.Show()
}
