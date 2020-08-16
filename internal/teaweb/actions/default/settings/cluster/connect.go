package cluster

import (
	"github.com/TeaWeb/build/internal/teacluster"
	"github.com/iwind/TeaGo/actions"
)

type ConnectAction actions.Action

// 连接到集群
func (this *ConnectAction) RunPost(params struct{}) {
	teacluster.SharedManager.Restart()
	this.Success()
}
