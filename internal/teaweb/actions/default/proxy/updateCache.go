package proxy

import (
	"github.com/TeaWeb/build/internal/teacache"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateCacheAction actions.Action

// 更新缓存设置
func (this *UpdateCacheAction) Run(params struct {
	ServerId string
	Policy   string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到要操作的代理服务")
	}

	if len(params.Policy) > 0 {
		policy := shared.NewCachePolicyFromFile(params.Policy)
		if policy == nil {
			this.Fail("找不到要使用的缓存策略")
		}
		this.Data["policy"] = maps.Map{
			"name":     policy.Name,
			"typeName": teacache.FindTypeName(policy.Type),
			"type":     policy.Type,
			"key":      policy.Key,
		}
	} else {
		this.Data["policy"] = maps.Map{
			"name":     "",
			"typeName": "",
			"type":     "",
		}
	}

	server.CacheOn = true
	server.CachePolicy = params.Policy
	err := server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()

	this.Success()
}
