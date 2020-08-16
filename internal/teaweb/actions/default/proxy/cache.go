package proxy

import (
	"github.com/TeaWeb/build/internal/teacache"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type CacheAction actions.Action

// 缓存设置
func (this *CacheAction) Run(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	this.Data["selectedTab"] = "cache"
	this.Data["server"] = server

	// 缓存策略
	this.Data["cachePolicy"] = nil
	if len(server.CachePolicy) > 0 {
		policy := shared.NewCachePolicyFromFile(server.CachePolicy)
		if policy != nil {
			this.Data["cachePolicy"] = policy
		}
	}
	this.Data["cachePolicyFile"] = server.CachePolicy

	cache, _ := teaconfigs.SharedCacheConfig()
	this.Data["cachePolicyList"] = lists.Map(cache.FindAllPolicies(), func(k int, v interface{}) interface{} {
		policy := v.(*shared.CachePolicy)
		return maps.Map{
			"filename": policy.Filename,
			"name":     policy.Name,
			"type":     teacache.FindTypeName(policy.Type),
		}
	})

	this.Show()
}
