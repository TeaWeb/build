package agent

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/string"
)

type IpAction actions.Action

// 探测主机访问Master的IP
func (this *IpAction) Run(params struct{}) {
	this.WriteString(stringutil.JSONEncode(maps.Map{
		"ip": this.RequestRemoteIP(),
	}))
}
