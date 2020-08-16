package proxy

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type ServerAction actions.Action

// 单个代理服务
func (this *ServerAction) RunGet(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		apiutils.Fail(this, "not found")
		return
	}

	apiutils.Success(this, maps.Map{
		"config": server,
	})
}
