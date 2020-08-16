package proxy

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
)

type DeleteAction actions.Action

// 删除
func (this *DeleteAction) Run(params struct {
	ServerId string
}) {
	this.Data["server"] = maps.Map{
		"id": params.ServerId,
	}

	this.Show()
}

func (this *DeleteAction) RunPost(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	// 检查有没有被引用
	serverList, err := teaconfigs.SharedServerList()
	if err != nil {
		this.Fail("删除失败：" + err.Error())
	}
	servers := serverList.FindAllServers()
	for _, s := range servers {
		if s.Id == server.Id {
			continue
		}
		if description, referred := s.RefersProxy(server.Id); referred {
			this.Fail("有别的代理服务在引用此代理服务：" + s.Description + "[" + description + "]，请删除引用后再次尝试")
		}
	}

	// 删除统计数据
	err = teadb.ServerValueDAO().DropServerTable(server.Id)
	if err != nil {
		this.Fail("删除统计数据失败：" + err.Error())
	}

	// 从list中删除
	serverList.RemoveServer(server.Filename)
	err = serverList.Save()
	if err != nil {
		this.Fail("删除失败：" + err.Error())
	}

	err = server.Delete()
	if err != nil {
		logs.Error(err)
		this.Fail("配置文件删除失败")
	}

	// 重启
	teaproxy.SharedManager.RemoveServer(server.Id)
	proxyutils.NotifyChange()

	this.Success()
}
