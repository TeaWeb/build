package rewrite

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type DeleteAction actions.Action

// 删除重写规则
func (this *DeleteAction) Run(params struct {
	ServerId   string
	LocationId string
	RewriteId  string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	rewriteList, err := server.FindRewriteList(params.LocationId)
	if err != nil {
		this.Fail(err.Error())
	}

	rewriteList.RemoveRewriteRule(params.RewriteId)
	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()

	this.Success()
}
