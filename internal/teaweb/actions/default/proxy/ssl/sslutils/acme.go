package sslutils

import (
	"github.com/TeaWeb/build/internal/teacluster"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"time"
)

// 重载证书
func ReloadACMECert(serverId string, taskId string) (errs []error) {
	if len(serverId) == 0 || len(taskId) == 0 {
		return
	}

	// 更新目前正在使用的Server
	runningServer := teaproxy.SharedManager.FindServer(serverId)
	if runningServer == nil || runningServer.SSL == nil || len(runningServer.SSL.Certs) == 0 {
		return
	}

	// 查找正在使用此任务的证书
	for _, cert := range runningServer.SSL.Certs {
		if cert.TaskId == taskId {
			err := cert.Validate()
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return
}

// 检查ACME证书更新
func RenewACMECerts() {
	logs.Println("[acme]check acme requests")

	// skip slave node
	node := teaconfigs.SharedNodeConfig()
	if node != nil && node.On && !node.IsMaster() {
		return
	}
	nodeDirty := false
	if node != nil && teacluster.SharedManager.IsActive() {
		// 集群节点状态
		nodeDirty = teacluster.SharedManager.IsChanged()
	}

	serverList, err := teaconfigs.SharedServerList()
	if err != nil {
		return
	}

	nodeIsChanged := false

	for _, server := range serverList.FindAllServers() {
		if server.SSL == nil || !server.SSL.On || len(server.SSL.CertTasks) == 0 {
			continue
		}
		serverIsChanged := false
		for _, task := range server.SSL.CertTasks {
			if !task.On {
				continue
			}
			if task.Request == nil {
				continue
			}
			date := task.Request.CertDate()
			if len(date[1]) == 0 {
				continue
			}
			if timeutil.Format("Y-m-d") >= date[1] {
				client, err := task.Request.Client()
				if err != nil {
					task.RunAt = time.Now().Unix()
					task.RunError = err.Error()
					logs.Error(err)
					serverIsChanged = true
					continue
				}
				err = task.Request.Renew(client)
				if err != nil {
					task.RunAt = time.Now().Unix()
					task.RunError = err.Error()
					logs.Error(err)
					serverIsChanged = true
					continue
				}

				task.RunAt = time.Now().Unix()
				task.RunError = ""
				serverIsChanged = true

				// 更新证书
				found := false
				for _, cert := range server.SSL.Certs {
					if cert.TaskId == task.Id {
						err = task.Request.WriteCertFile(Tea.ConfigFile(cert.CertFile))
						if err != nil {
							logs.Error(err)
						}

						err = task.Request.WriteKeyFile(Tea.ConfigFile(cert.KeyFile))
						if err != nil {
							logs.Error(err)
						}

						found = true
					}
				}

				// 重新加载证书
				if found {
					errs := ReloadACMECert(server.Id, task.Id)
					for _, err2 := range errs {
						logs.Println("[acme]reload acme task:", err2.Error())
					}
				}
			}
		}

		if serverIsChanged {
			nodeIsChanged = true
			server.Save()
		}
	}

	// 如果先前节点没有变更，则自动推送到集群
	if !nodeDirty && nodeIsChanged {
		node := teaconfigs.SharedNodeConfig()
		if node != nil && node.On && node.IsMaster() && teacluster.SharedManager.IsActive() {
			teacluster.SharedManager.PushItems()
			teacluster.SharedManager.SetIsChanged(false)
		}
	}
}
