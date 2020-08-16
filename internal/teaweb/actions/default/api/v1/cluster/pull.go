package cluster

import (
	"github.com/TeaWeb/build/internal/teacluster"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
)

type PullAction actions.Action

// 从cluster pull
func (this *PullAction) RunGet(params struct{}) {
	node := teaconfigs.SharedNodeConfig()
	if node == nil {
		apiutils.Fail(this, "the node has not been configured")
		return
	}

	if !teacluster.SharedManager.IsActive() {
		apiutils.Fail(this, "the node is not connecting to cluster")
		return
	}

	if !node.IsMaster() {
		teacluster.SharedManager.PullItems()
		teacluster.SharedManager.SetIsChanged(false)
	}

	apiutils.SuccessOK(this)
}
