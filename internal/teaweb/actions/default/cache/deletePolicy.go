package cache

import (
	"github.com/TeaWeb/build/internal/teacache"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
)

type DeletePolicyAction actions.Action

// 删除缓存策略
func (this *DeletePolicyAction) Run(params struct {
	Filename string
}) {
	if len(params.Filename) == 0 {
		this.Fail("请指定要删除的缓存策略")
	}

	policy := shared.NewCachePolicyFromFile(params.Filename)
	if policy == nil {
		this.Fail("找不到要删除的缓存策略")
	}

	// 删除Server和Location中的cache
	isChanged := false
	serverList, _ := teaconfigs.SharedServerList()
	if serverList != nil {
		for _, server := range serverList.FindAllServers() {
			// 删除server中的缓存策略
			if server.CachePolicy == params.Filename {
				isChanged = true
				server.CachePolicy = ""

				for _, location := range server.Locations { // 删除Location中的缓存策略
					if location.CachePolicy == params.Filename {
						location.CachePolicy = ""
					}
				}

				err := server.Save()
				if err != nil {
					this.Fail("删除失败：" + err.Error())
				}
			} else { // 删除Location中的缓存策略
				serverChanged := false
				for _, location := range server.Locations {
					if location.CachePolicy == params.Filename {
						location.CachePolicy = ""
						isChanged = true
						serverChanged = true
					}
				}
				if serverChanged {
					err := server.Save()
					if err != nil {
						this.Fail("删除失败：" + err.Error())
					}
				}
			}
		}
	}

	config, _ := teaconfigs.SharedCacheConfig()
	config.DeletePolicy(params.Filename)
	err := config.Save()
	if err != nil {
		this.Fail("删除失败：" + err.Error())
	}

	err = policy.Delete()
	if err != nil {
		logs.Error(err)
	}

	if isChanged {
		proxyutils.NotifyChange()
	}

	// 重置缓存策略实例
	teacache.ResetCachePolicyManager(policy.Filename)

	this.Success()
}
