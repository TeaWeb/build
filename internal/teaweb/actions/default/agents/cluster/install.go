package cluster

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"sync"
)

type InstallAction actions.Action

// 安装
func (this *InstallAction) Run(params struct {
	Hosts    []string
	Master   string
	Dir      string
	GroupId  string
	Port     int
	Username string
	AuthType string
	Password string
	Key      string
}) {
	// 禁止demo
	if teaconst.DemoEnabled {
		this.Fail("DEMO版本无法操作")
	}

	// 检查数量限制
	group := agents.SharedGroupList().FindGroup(params.GroupId)
	if group != nil {
		if group.MaxAgents > 0 && group.CountAgents >= group.MaxAgents {
			this.Fail("选择的分组不能超过最大Agent数量限制")
		}
		if !group.IsDateAvailable() {
			this.Fail(" 选择的分组不在有效期限内")
		}
	}

	if len(params.Hosts) == 0 {
		this.Data["states"] = []interface{}{}
		this.Success()
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(params.Hosts))

	states := []maps.Map{}
	stateLocker := sync.Mutex{}

	for _, host := range params.Hosts {
		func(addr string) {
			defer wg.Done()

			installer := agentutils.NewInstaller()
			installer.Port = params.Port
			installer.Host = addr
			installer.Master = params.Master
			installer.GroupId = params.GroupId
			installer.Dir = params.Dir
			installer.AuthUsername = params.Username
			installer.AuthType = params.AuthType
			installer.AuthPassword = params.Password
			installer.AuthKey = []byte(params.Key)
			err := installer.Start()
			result := maps.Map{
				"addr":        addr,
				"name":        installer.HostName,
				"ip":          installer.HostIP,
				"isInstalled": installer.IsInstalled,
				"hasError":    err != nil,
			}
			if err != nil {
				result["result"] = err.Error()
			} else {
				result["result"] = "安装成功"
			}
			stateLocker.Lock()
			states = append(states, result)
			stateLocker.Unlock()
		}(host)
	}

	wg.Wait()

	this.Data["states"] = states

	// 重建索引
	err := agents.SharedGroupList().BuildIndexes()
	if err != nil {
		logs.Error(err)
	}

	this.Success()
}
