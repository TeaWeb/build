package locations

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/locations/locationutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type HeadersAction actions.Action

// 自定义Http Header
func (this *HeadersAction) Run(params struct {
	ServerId   string // 必填
	LocationId string
}) {
	locationutils.SetCommonInfo(this, params.ServerId, params.LocationId, "headers")

	this.Data["headerQuery"] = maps.Map{
		"serverId":   params.ServerId,
		"locationId": params.LocationId,
	}

	this.Show()
}
