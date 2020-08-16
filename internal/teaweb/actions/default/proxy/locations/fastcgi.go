package locations

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/locations/locationutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type FastcgiAction actions.Action

// Fastcgi设置
func (this *FastcgiAction) Run(params struct {
	ServerId   string
	LocationId string
}) {
	locationutils.SetCommonInfo(this, params.ServerId, params.LocationId, "fastcgi")

	this.Data["queryParams"] = maps.Map{
		"serverId":   params.ServerId,
		"locationId": params.LocationId,
	}

	this.Show()
}
