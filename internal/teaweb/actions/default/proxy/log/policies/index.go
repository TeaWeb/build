package policies

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/tealogs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

type IndexAction actions.Action

// 日志策略
func (this *IndexAction) RunGet(params struct{}) {
	this.Data["policyList"] = lists.Map(teaconfigs.SharedAccessLogStoragePolicyList().FindAllPolicies(), func(k int, v interface{}) interface{} {
		policy := v.(*teaconfigs.AccessLogStoragePolicy)
		format, _ := policy.Options["format"]
		return maps.Map{
			"id":         policy.Id,
			"name":       policy.Name,
			"on":         policy.On,
			"typeName":   tealogs.FindStorageTypeName(policy.Type),
			"formatName": tealogs.FindStorageFormatName(types.String(format)),
		}
	})

	this.Show()
}
