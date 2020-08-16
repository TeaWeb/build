package waf

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/build/internal/teawaf"
	"github.com/TeaWeb/build/internal/teawaf/checkpoints"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"regexp"
	"strings"
)

type TestAction actions.Action

// 测试
func (this *TestAction) RunGet(params struct {
	WafId   string
	Inbound bool
}) {
	waf := teaconfigs.SharedWAFList().FindWAF(params.WafId)
	if waf == nil {
		this.Fail("找不到WAF")
	}

	this.Data["inbound"] = params.Inbound

	this.Data["config"] = maps.Map{
		"id":            waf.Id,
		"name":          waf.Name,
		"countInbound":  waf.CountInboundRuleSets(),
		"countOutbound": waf.CountOutboundRuleSets(),
	}

	// 数据列表
	paramList := []string{}
	if params.Inbound {
		for _, group := range waf.Inbound {
			if !group.On {
				continue
			}
			for _, set := range group.RuleSets {
				if !set.On {
					continue
				}
				for _, rule := range set.Rules {
					if lists.ContainsString(paramList, rule.Param) {
						continue
					}
					paramList = append(paramList, rule.Param)
				}
			}
		}
	} else {
		for _, group := range waf.Outbound {
			if !group.On {
				continue
			}
			for _, set := range group.RuleSets {
				if !set.On {
					continue
				}
				for _, rule := range set.Rules {
					if lists.ContainsString(paramList, rule.Param) {
						continue
					}
					paramList = append(paramList, rule.Param)
				}
			}
		}
	}

	reg := regexp.MustCompile("^\\${([\\w.-]+)}$")
	this.Data["params"] = lists.Map(paramList, func(k int, v interface{}) interface{} {
		param := v.(string)

		prefix := ""
		result := reg.FindStringSubmatch(param)
		if len(result) > 0 {
			match := result[1]
			pieces := strings.SplitN(match, ".", 2)
			prefix = pieces[0]
			if len(pieces) == 2 {
				param = pieces[1]
			} else {
				param = ""
			}
		}

		checkpointName := ""
		if len(prefix) > 0 {
			checkpoint := checkpoints.FindCheckpointDefinition(prefix)
			if checkpoint != nil {
				def := checkpoints.FindCheckpointDefinition(prefix)
				if def != nil {
					checkpointName = def.Name
				}
			}
		}

		return maps.Map{
			"param":      strings.Trim(prefix+"."+param, "."),
			"fullParam":  types.String(v),
			"prefix":     prefix,
			"checkpoint": checkpointName,
		}
	})

	this.Show()
}

// 提交测试数据
func (this *TestAction) RunPost(params struct {
	WafId   string
	Params  []string
	Values  []string
	Inbound bool
}) {
	waf := teaconfigs.SharedWAFList().FindWAF(params.WafId)
	if waf == nil {
		this.Fail("找不到WAF")
	}

	valueMap := maps.Map{}
	for index, param := range params.Params {
		if index < len(params.Values) {
			valueMap[param] = params.Values[index]
		}
	}

	result := []string{}
	waf.Init()
	defer waf.Stop()

	matched := false
	setName := ""
	action := ""

	groups := []*teawaf.RuleGroup{}
	if params.Inbound {
		groups = waf.Inbound
	} else {
		groups = waf.Outbound
	}

Loop:
	for _, group := range groups {
		if !group.On {
			continue
		}
		result = append(result, "开始检查规则分组 '"+group.Name+"' "+fmt.Sprintf("%d 个规则集", len(group.RuleSets))+" ...")
		if len(group.RuleSets) == 0 {
			result = append(result, "　　跳过")
			continue
		}
		for _, set := range group.RuleSets {
			result = append(result, "　　开始检查规则集 '"+set.Name+"' "+fmt.Sprintf("%d 个规则 ...", len(set.Rules)))

			if len(set.Rules) == 0 {
				result = append(result, "　　　　跳过")
				continue
			}

			found := false
			if set.Connector == teawaf.RuleConnectorAnd {
				found = true
			}
			for _, rule := range set.Rules {
				value := teautils.ParseVariables(rule.Param, func(varName string) (value string) {
					v, _ := valueMap["${"+varName+"}"]
					return types.String(v)
				})
				if rule.Test(value) {
					if set.Connector == teawaf.RuleConnectorOr {
						found = true
					}
				} else {
					if set.Connector == teawaf.RuleConnectorAnd {
						found = false
					}
				}
			}
			if found {
				matched = true
				action = set.Action
				result = append(result, "　　　　匹配成功， 动作："+strings.ToUpper(set.Action))
				setName = set.Name
				break Loop
			} else {
				result = append(result, "　　　　 没有匹配的规则")
			}
		}
	}

	this.Data["result"] = result
	this.Data["action"] = action
	this.Data["setName"] = setName
	this.Data["matched"] = matched
	this.Success()
}
