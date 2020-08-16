package settings

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type UpdateAction actions.Action

// 修改全局代理设置
func (this *UpdateAction) RunGet(params struct{}) {
	this.Data["setting"] = teaconfigs.LoadProxySetting()
	this.Show()
}

func (this *UpdateAction) RunPost(params struct {
	MatchDomainStrictly bool

	Must *actions.Must
}) {
	setting := teaconfigs.LoadProxySetting()
	setting.MatchDomainStrictly = params.MatchDomainStrictly
	err := setting.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()

	this.Success()
}
