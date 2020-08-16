package backend

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/scheduling"
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
)

type DataAction actions.Action

// 后端服务器数据
func (this *DataAction) Run(params struct {
	ServerId   string
	LocationId string
	Websocket  bool
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	// 自定义默认分组
	if len(server.RequestGroups) == 0 {
		group := teaconfigs.NewRequestGroup()
		group.IsDefault = true
		group.Id = "default"
		group.Name = "默认分组"
		server.AddRequestGroup(group)
		err := server.Save()
		if err != nil {
			logs.Error(err)
		}
	}

	backendList, err := server.FindBackendList(params.LocationId, params.Websocket)
	if err != nil {
		this.Fail(err.Error())
	}

	normalBackends := []*teaconfigs.BackendConfig{}
	backupBackends := []*teaconfigs.BackendConfig{}
	runningServer := teaproxy.SharedManager.FindServer(server.Id)
	for _, backend := range backendList.AllBackends() {
		// 是否下线以及错误次数
		if runningServer != nil {
			runningBackendList, err := runningServer.FindBackendList(params.LocationId, params.Websocket)
			if err == nil {
				runningBackend := runningBackendList.FindBackend(backend.Id)
				if runningBackend != nil {
					backend.IsDown = runningBackend.IsDown
					backend.DownTime = runningBackend.DownTime
					backend.CurrentFails = runningBackend.CurrentFails
					backend.CurrentConns = runningBackend.CurrentConns
				}
			}
		}

		if backend.IsBackup {
			backupBackends = append(backupBackends, backend)
		} else {
			normalBackends = append(normalBackends, backend)
		}
	}

	this.Data["normalBackends"] = normalBackends
	this.Data["backupBackends"] = backupBackends

	// 算法
	schedulingConfig := backendList.SchedulingConfig()
	if schedulingConfig == nil {
		this.Data["scheduling"] = scheduling.FindSchedulingType("random")
	} else {
		s := scheduling.FindSchedulingType(schedulingConfig.Code)
		if s == nil {
			this.Data["scheduling"] = scheduling.FindSchedulingType("random")
		} else {
			this.Data["scheduling"] = s
		}
	}

	// 分组
	if len(server.RequestGroups) > 0 {
		this.Data["groups"] = server.RequestGroups
	} else {
		this.Data["groups"] = []*teaconfigs.RequestGroup{}
	}

	this.Success()
}
