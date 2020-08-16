package servergroups

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

func (this *IndexAction) RunGet(params struct{}) {
	this.Data["selectedMenu"] = "index"

	groupList := teaconfigs.SharedServerGroupList()
	groups := groupList.Groups
	groupMaps := []maps.Map{}
	isChanged := false
	for _, group := range groups {
		serverMaps := []maps.Map{}
		for _, serverId := range group.ServerIds {
			server := teaconfigs.NewServerConfigFromId(serverId)
			if server == nil {
				isChanged = true
				continue
			}
			serverMaps = append(serverMaps, maps.Map{
				"id":          server.Id,
				"description": server.Description,
			})
		}

		groupMaps = append(groupMaps, maps.Map{
			"id":      group.Id,
			"name":    group.Name,
			"isOn":    group.IsOn,
			"servers": serverMaps,
		})
	}

	// 是否有未分组的
	serverList, err := teaconfigs.SharedServerList()
	if err != nil {
		logs.Error(err)
	} else {
		ungroupServerMaps := []maps.Map{}
		for _, server := range serverList.FindAllServers() {
			if !groupList.ContainsServer(server.Id) {
				ungroupServerMaps = append(ungroupServerMaps, maps.Map{
					"id":          server.Id,
					"description": server.Description,
				})
			}
		}
		if len(ungroupServerMaps) > 0 {
			groupMaps = append(groupMaps, maps.Map{
				"id":      "",
				"name":    "[未分组]",
				"isOn":    true,
				"servers": ungroupServerMaps,
			})
		}
	}

	this.Data["groups"] = groupMaps

	if isChanged {
		err := teaconfigs.SharedServerGroupList().Save()
		if err != nil {
			logs.Error(err)
		}
	}

	this.Show()
}
