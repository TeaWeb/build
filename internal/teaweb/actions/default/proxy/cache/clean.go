package cache

import (
	"github.com/TeaWeb/build/internal/teacache"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/actions"
)

type CleanAction actions.Action

// 清除缓存
func (this *CleanAction) RunPost(params struct {
	Filename string
	Key      string
}) {
	policy := shared.NewCachePolicyFromFile(params.Filename)
	if policy == nil {
		this.Data["result"] = "找不到Policy"
		this.Fail()
	}

	manager := teacache.FindCachePolicyManager(params.Filename)
	if manager == nil {
		manager = teacache.NewManagerFromConfig(policy)
		defer manager.Close()
	}

	if manager == nil {
		this.Fail("找不到管理器")
	}

	err := manager.Delete(params.Key)
	if err != nil {
		this.Fail("ERROR:" + err.Error())
	}

	this.Success()
}
