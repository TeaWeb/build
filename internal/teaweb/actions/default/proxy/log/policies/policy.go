package policies

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/tealogs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

type PolicyAction actions.Action

// 策略详情
func (this *PolicyAction) RunGet(params struct {
	PolicyId string
}) {
	policy := teaconfigs.NewAccessLogStoragePolicyFromId(params.PolicyId)
	if policy == nil {
		this.Fail("找不到策略")
	}

	format, _ := policy.Options["format"]
	template, _ := policy.Options["template"]

	this.Data["policy"] = maps.Map{
		"id":           policy.Id,
		"on":           policy.On,
		"name":         policy.Name,
		"format":       types.String(format),
		"formatName":   tealogs.FindStorageFormatName(types.String(format)),
		"type":         policy.Type,
		"typeName":     tealogs.FindStorageTypeName(policy.Type),
		"templateCode": types.String(template),
		"options":      policy.Options,
		"cond":         policy.Cond,
	}

	this.Data["configItems"] = FindAllUsingPolicy(policy.Id)

	// syslog
	this.Data["syslogPriorities"] = tealogs.SyslogStoragePriorities

	this.Show()
}
