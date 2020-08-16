package locations

import (
	"github.com/TeaWeb/build/internal/teacache"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/locations/locationutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type CacheAction actions.Action

// 缓存设置
func (this *CacheAction) Run(params struct {
	ServerId   string
	LocationId string
}) {
	_, location := locationutils.SetCommonInfo(this, params.ServerId, params.LocationId, "cache")

	// 缓存策略
	this.Data["cachePolicy"] = nil
	if len(location.CachePolicy) > 0 {
		policy := shared.NewCachePolicyFromFile(location.CachePolicy)
		if policy != nil {
			this.Data["cachePolicy"] = policy
		}
	}
	this.Data["cachePolicyFile"] = location.CachePolicy
	if !location.CacheOn {
		this.Data["cachePolicyFile"] = "none"
	}

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
