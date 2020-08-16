package proxy

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
)

type IndexAction actions.Action

// 代理首页
func (this *IndexAction) Run(params struct {
}) {
	// 跳转到分组的第一个
	groupList := teaconfigs.SharedServerGroupList()
	serverId := ""
	for _, group := range groupList.Groups {
		if !group.IsOn {
			continue
		}
		if len(group.ServerIds) == 0 {
			continue
		}
		serverId = group.ServerIds[0]
		break
	}
	if len(serverId) > 0 {
		this.RedirectURL("/proxy/board?serverId=" + serverId)
		return
	}

	// 没有分组，就跳到所有的服务的第一个
	serverList, err := teaconfigs.SharedServerList()
	if err != nil {
		logs.Error(err)
		return
	}
	servers := serverList.FindAllServers()
	if len(servers) > 0 {
		this.RedirectURL("/proxy/board?serverId=" + servers[0].Id)
	} else {
		this.Show()
	}
}
