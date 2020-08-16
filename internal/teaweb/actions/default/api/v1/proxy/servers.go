package proxy

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type ServersAction actions.Action

// 代理服务列表
func (this *ServersAction) RunGet(params struct{}) {
	serverList, err := teaconfigs.SharedServerList()
	if err != nil {
		apiutils.Success(this, []interface{}{})
	}

	result := []maps.Map{}
	for _, server := range serverList.FindAllServers() {
		result = append(result, maps.Map{
			"config": server,
		})
	}
	apiutils.Success(this, result)
}
