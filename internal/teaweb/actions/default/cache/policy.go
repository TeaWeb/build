package cache

import (
	"github.com/TeaWeb/build/internal/teacache"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/actionutils"
	"github.com/iwind/TeaGo/maps"
)

type PolicyAction struct {
	actionutils.ParentAction
}

// 缓存策略详情
func (this *PolicyAction) Run(params struct {
	Filename string
}) {
	this.SecondMenu("policy")

	policy := shared.NewCachePolicyFromFile(params.Filename)
	if policy == nil {
		this.Fail("找不到Policy")
	}

	this.Data["policy"] = policy

	// 类型
	this.Data["type"] = teacache.FindType(policy.Type)

	// 正在使用此缓存策略的项目
	configItems := []maps.Map{}
	serverList, _ := teaconfigs.SharedServerList()
	if serverList != nil {
		for _, server := range serverList.FindAllServers() {

			if server.CachePolicy == policy.Filename {
				configItems = append(configItems, maps.Map{
					"type":   "server",
					"server": server.Description,
					"link":   "/proxy/cache?serverId=" + server.Id,
				})
			}

			for _, location := range server.Locations {
				if location.CachePolicy == policy.Filename {
					configItems = append(configItems, maps.Map{
						"type":     "location",
						"server":   server.Description,
						"location": location.Pattern,
						"link":     "/proxy/locations/cache?serverId=" + server.Id + "&locationId=" + location.Id,
					})
				}
			}
		}
	}

	this.Data["configItems"] = configItems

	this.Show()
}
