package certutils

import (
	"github.com/TeaWeb/build/internal/teacluster"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"time"
)

// 检查ACME证书更新
func RenewACMECerts() {
	logs.Println("[acme]check acme certs requests")

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

	nodeIsChanged := false
	certList := teaconfigs.SharedSSLCertList()
	tasks := certList.Tasks
	if len(tasks) == 0 {
		return
	}

	taskIsChanged := false
	for _, task := range tasks {
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
				taskIsChanged = true
				continue
			}
			err = task.Request.Renew(client)
			if err != nil {
				task.RunAt = time.Now().Unix()
				task.RunError = err.Error()
				logs.Error(err)
				taskIsChanged = true
				continue
			}

			task.RunAt = time.Now().Unix()
			task.RunError = ""
			taskIsChanged = true

			// 更新证书
			for _, cert := range certList.Certs {
				if cert.TaskId != task.Id {
					continue
				}
				err := task.Request.WriteCertFile(Tea.ConfigFile(cert.CertFile))
				if err != nil {
					logs.Error(err)
				}

				err = task.Request.WriteKeyFile(Tea.ConfigFile(cert.KeyFile))
				if err != nil {
					logs.Error(err)
				}

				// 重新加载证书
				servers := teaproxy.SharedManager.FindAllServers()
				for _, server := range servers {
					certs := server.FindCerts(cert.Id)
					if len(certs) > 0 {
						for _, c := range certs {
							err := c.Validate()
							if err != nil {
								logs.Error(err)
							}
						}
					}
				}
			}
		}
	}

	// 保存修改
	if taskIsChanged {
		err := certList.Save()
		if err != nil {
			logs.Error(err)
		}
		nodeIsChanged = true
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
