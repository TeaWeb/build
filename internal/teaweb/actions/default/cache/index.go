package cache

import (
	"github.com/TeaWeb/build/internal/teacache"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

// 缓存首页
func (this *IndexAction) Run(params struct{}) {
	policyList := []maps.Map{}
	config, _ := teaconfigs.SharedCacheConfig()
	for _, policy := range config.FindAllPolicies() {
		policyList = append(policyList, maps.Map{
			"name":     policy.Name,
			"filename": policy.Filename,
			"key":      policy.Key,

			"hasCapacity": policy.CapacitySize() > 0,
			"capacity":    policy.Capacity,

			"hasLife": policy.LifeDuration() > 0,
			"life":    policy.Life,

			"status": lists.Join(policy.Status, ", ", func(k int, v interface{}) interface{} {
				return v
			}),

			"hasMaxSize": policy.MaxDataSize() > 0,
			"maxSize":    policy.MaxSize,

			"type":     policy.Type,
			"typeName": teacache.FindTypeName(policy.Type),
			"options":  policy.Options,
		})
	}
	this.Data["policyList"] = policyList

	this.Show()
}
