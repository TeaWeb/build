package waf

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/waf/wafutils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
)

type DeleteAction actions.Action

// 删除
func (this *DeleteAction) RunPost(params struct {
	WafId string
}) {
	if len(params.WafId) == 0 {
		this.Fail("请输入要删除的WAF ID")
	}

	filename := "waf." + params.WafId + ".conf"
	path := Tea.ConfigFile(filename)
	file := files.NewFile(path)
	if !file.Exists() {
		this.Fail("要删除的WAF不存在")
	}

	// 通知刷新
	if wafutils.IsPolicyUsed(params.WafId) {
		this.Fail("此策略正在被使用，不能删除，点击“详情”查看使用此WAF策略的项目")
	}

	wafList := teaconfigs.SharedWAFList()
	wafList.RemoveFile(filename)
	err := wafList.Save()
	if err != nil {
		this.Fail("删除失败：" + err.Error())
	}

	err = file.Delete()
	if err != nil {
		this.Fail("删除失败：" + err.Error())
	}

	this.Success()
}
