package fastcgi

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type DeleteAction actions.Action

// 删除Fastcgi
func (this *DeleteAction) Run(params struct {
	ServerId   string
	LocationId string
	FastcgiId  string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	fastcgiList, err := server.FindFastcgiList(params.LocationId)
	if err != nil {
		this.Fail(err.Error())
	}

	fastcgiList.RemoveFastcgi(params.FastcgiId)
	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()

	this.Success()
}
