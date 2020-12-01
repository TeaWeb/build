package cluster

import (
	"github.com/TeaWeb/build/internal/teacluster"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/rands"
)

type UpdateAction actions.Action

// show update form
func (this *UpdateAction) RunGet(params struct{}) {
	this.Data["teaMenu"] = "cluster"

	node := teaconfigs.SharedNodeConfig()
	if node == nil {
		node = &teaconfigs.NodeConfig{}
		node.On = true
	}
	this.Data["node"] = node
	this.Data["roles"] = teaconfigs.AllNodeRoles()

	this.Show()
}

// submit form
func (this *UpdateAction) RunPost(params struct {
	Name          string
	Role          string
	ClusterId     string
	ClusterSecret string
	ClusterAddr   string
	On            bool

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入节点名称").
		Field("role", params.Role).
		Require("请选择当前节点角色").
		Expect(func() (message string, success bool) {
			if teaconfigs.ExistNodeRole(params.Role) {
				return "", true
			}
			return "请选择正确的角色", false
		}).
		Field("clusterId", params.ClusterId).
		Require("请输入集群ID").
		Field("clusterSecret", params.ClusterSecret).
		Require("请输入集群密钥").
		Field("clusterAddr", params.ClusterAddr).
		Require("请输入集群通讯地址")

	node := teaconfigs.SharedNodeConfig()
	if node == nil {
		node = &teaconfigs.NodeConfig{}
	}
	node.Name = params.Name
	node.Role = params.Role
	node.ClusterId = params.ClusterId
	node.ClusterSecret = params.ClusterSecret
	node.ClusterAddr = params.ClusterAddr
	node.On = params.On

	if len(node.Id) == 0 {
		node.Id = rands.HexString(16)
	}

	err := node.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	if !node.On {
		teacluster.SharedManager.SetIsChanged(false)
	}

	teacluster.SharedManager.Restart()

	this.Success()
}
