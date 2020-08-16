package cluster

import (
	"github.com/TeaWeb/build/internal/teacluster"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type NodeAction actions.Action

// 节点信息
func (this *NodeAction) RunGet(params struct{}) {
	config := teaconfigs.SharedNodeConfig()
	if config == nil {
		apiutils.Fail(this, "not a node yet")
		return
	}
	apiutils.Success(this, maps.Map{
		"isActive":  teacluster.SharedManager.IsActive(),
		"isChanged": teacluster.SharedManager.IsChanged(),
		"config":    config,
	})
}
