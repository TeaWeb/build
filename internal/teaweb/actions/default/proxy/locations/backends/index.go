package backends

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/locations/locationutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

// 后端服务器
func (this *IndexAction) Run(params struct {
	ServerId   string
	LocationId string
}) {
	this.Data["queryParams"] = maps.Map{
		"serverId":   params.ServerId,
		"locationId": params.LocationId,
	}

	server, _ := locationutils.SetCommonInfo(this, params.ServerId, params.LocationId, "backends")
	this.Data["isTCP"] = server.IsTCP()
	this.Data["isHTTP"] = server.IsHTTP()

	this.Show()
}
