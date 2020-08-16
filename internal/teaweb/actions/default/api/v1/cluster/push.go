package cluster

import (
	"github.com/TeaWeb/build/internal/teacluster"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
)

type PushAction actions.Action

// push到cluster
func (this *PushAction) RunGet(params struct{}) {
	node := teaconfigs.SharedNodeConfig()
	if node == nil {
		apiutils.Fail(this, "the node has not been configured")
		return
	}

	if !teacluster.SharedManager.IsActive() {
		apiutils.Fail(this, "the node is not connecting to cluster")
		return
	}

	if node.IsMaster() {
		teacluster.SharedManager.PushItems()
		teacluster.SharedManager.SetIsChanged(false)
	}

	apiutils.SuccessOK(this)
}
