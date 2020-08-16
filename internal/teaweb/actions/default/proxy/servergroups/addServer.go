package servergroups

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"net/http"
)

type AddServerAction actions.Action

func (this *AddServerAction) RunGet(params struct {
	GroupId string
}) {
	groupList := teaconfigs.SharedServerGroupList()
	group := groupList.Find(params.GroupId)
	if group == nil {
		this.Error("not found", http.StatusNotFound)
		return
	}

	this.Data["group"] = group

	// 所有代理服务列表
	serverList, err := teaconfigs.SharedServerList()
	if err != nil {
		this.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	serverMaps := []maps.Map{}
	for _, server := range serverList.FindAllServers() {
		if groupList.ContainsServer(server.Id) {
			continue
		}
		serverMaps = append(serverMaps, maps.Map{
			"id":          server.Id,
			"description": server.Description,
		})
	}
	this.Data["servers"] = serverMaps

	this.Show()
}

func (this *AddServerAction) RunPost(params struct {
	GroupId   string
	ServerIds []string

	Must *actions.Must
}) {
	groupList := teaconfigs.SharedServerGroupList()
	group := groupList.Find(params.GroupId)
	if group == nil {
		this.Error("not found", http.StatusNotFound)
		return
	}

	group.Add(params.ServerIds...)
	err := groupList.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
