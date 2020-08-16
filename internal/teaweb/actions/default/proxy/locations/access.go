package locations

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/locations/locationutils"
	"github.com/iwind/TeaGo/actions"
)

type AccessAction actions.Action

// 访问控制
func (this *AccessAction) Run(params struct {
	ServerId   string
	LocationId string
}) {
	_, location := locationutils.SetCommonInfo(this, params.ServerId, params.LocationId, "access")

	this.Data["policy"] = location.AccessPolicy

	this.Show()
}
