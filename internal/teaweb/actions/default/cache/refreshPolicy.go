package cache

import (
	"github.com/TeaWeb/build/internal/teacache"
	"github.com/TeaWeb/build/internal/teacluster"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/actionutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

type RefreshPolicyAction struct {
	actionutils.ParentAction
}

func (this *RefreshPolicyAction) RunGet(params struct {
	Filename string
}) {
	this.SecondMenu("refresh")
	policy := shared.NewCachePolicyFromFile(params.Filename)
	if policy == nil {
		this.ErrorPage("找不到Policy")
		return
	}

	this.Data["policy"] = policy

	this.Show()
}

func (this *RefreshPolicyAction) RunPost(params struct {
	Filename string
	Prefixes string
	Must     *actions.Must
}) {
	policy := shared.NewCachePolicyFromFile(params.Filename)
	if policy == nil {
		this.Data["result"] = "找不到Policy"
		this.Fail()
	}

	prefixes := []string{}
	if len(params.Prefixes) > 0 {
		lines := strings.Split(params.Prefixes, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			prefixes = append(prefixes, line)
		}
	}

	if len(prefixes) == 0 {
		this.Success()
	}

	manager := teacache.FindCachePolicyManager(params.Filename)
	if manager == nil {
		manager = teacache.NewManagerFromConfig(policy)
		defer func() {
			_ = manager.Close()
		}()
	}

	if manager == nil {
		this.Fail("找不到管理器")
	}

	count, err := manager.DeletePrefixes(prefixes)
	if err != nil {
		this.Fail("刷新失败：" + err.Error())
	}

	this.Data["count"] = count

	// 同步到集群
	action := new(teacluster.RunAction)
	action.Cmd = "cache.refresh"
	action.Data = maps.Map{
		"filename": params.Filename,
		"prefixes": prefixes,
	}
	err = teacluster.SharedManager.Write(action)
	if err != nil {
		this.Fail("同步到集群失败：" + err.Error())
	}

	this.Success()
}
