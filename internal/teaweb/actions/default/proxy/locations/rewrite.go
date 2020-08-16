package locations

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/locations/locationutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type RewriteAction actions.Action

// 重写规则
func (this *RewriteAction) Run(params struct {
	ServerId   string
	LocationId string
}) {
	locationutils.SetCommonInfo(this, params.ServerId, params.LocationId, "rewrite")

	this.Data["queryParams"] = maps.Map{
		"serverId":   params.ServerId,
		"locationId": params.LocationId,
	}

	this.Show()
}
