package waf

import (
	"encoding/json"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teawaf"
	"github.com/TeaWeb/build/internal/teawaf/checkpoints"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/waf/wafutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"regexp"
	"strings"
)

type RuleUpdateAction actions.Action

// 修改规则集
func (this *RuleUpdateAction) RunGet(params struct {
	WafId   string
	GroupId string
	SetId   string
}) {
	waf := teaconfigs.SharedWAFList().FindWAF(params.WafId)
	if waf == nil {
		this.Fail("找不到WAF")
	}

	this.Data["config"] = maps.Map{
		"id":            waf.Id,
		"name":          waf.Name,
		"countInbound":  waf.CountInboundRuleSets(),
		"countOutbound": waf.CountOutboundRuleSets(),
		"inbound":       waf.Inbound,
	}

	group := waf.FindRuleGroup(params.GroupId)
	if group == nil {
		this.Fail("找不到分组")
	}

	this.Data["inbound"] = group.IsInbound
	this.Data["outbound"] = !group.IsInbound

	set := group.FindRuleSet(params.SetId)
	if set == nil {
		this.Fail("找不到规则集")
	}
	err := set.Init()
	if err != nil {
		this.Fail("规则校验失败：" + err.Error())
	}

	reg := regexp.MustCompile("^\\${([\\w.-]+)}$")
	this.Data["set"] = set
	this.Data["oldRules"] = lists.Map(set.Rules, func(k int, v interface{}) interface{} {
		rule := v.(*teawaf.Rule)

		prefix := ""
		param := ""
		result := reg.FindStringSubmatch(rule.Param)
		if len(result) > 0 {
			match := result[1]
			pieces := strings.SplitN(match, ".", 2)
			prefix = pieces[0]
			if len(pieces) == 2 {
				param = pieces[1]
			}
		}

		return maps.Map{
			"prefix":   prefix,
			"param":    param,
			"operator": rule.Operator,
			"value":    rule.Value,
			"case":     rule.IsCaseInsensitive,
			"options":  rule.CheckpointOptions,
		}
	})

	this.Data["group"] = group
	this.Data["connectors"] = []maps.Map{
		{
			"name":        "和 (AND)",
			"value":       teawaf.RuleConnectorAnd,
			"description": "所有规则都满足才视为匹配",
		},
		{
			"name":        "或 (OR)",
			"value":       teawaf.RuleConnectorOr,
			"description": "任一规则满足了就视为匹配",
		},
	}

	// check points
	checkpointList := []maps.Map{}
	for _, def := range checkpoints.AllCheckpoints {
		if (group.IsInbound && def.Instance.IsRequest()) || (!group.IsInbound && !def.Instance.IsRequest()) {
			checkpointList = append(checkpointList, maps.Map{
				"name":         def.Name,
				"prefix":       def.Prefix,
				"description":  def.Description,
				"hasParams":    def.HasParams,
				"paramOptions": def.Instance.ParamOptions(),
				"options": lists.Map(def.Instance.Options(), func(k int, v interface{}) interface{} {
					{
						option, ok := v.(*checkpoints.FieldOption)
						if ok {
							return maps.Map{
								"type":        option.Type(),
								"name":        option.Name,
								"maxLength":   option.MaxLength,
								"code":        option.Code,
								"rightLabel":  option.RightLabel,
								"value":       option.Value,
								"isRequired":  option.IsRequired,
								"size":        option.Size,
								"comment":     option.Comment,
								"placeholder": option.Placeholder,
							}
						}
					}

					{
						option, ok := v.(*checkpoints.OptionsOption)
						if ok {
							return maps.Map{
								"type":       option.Type(),
								"name":       option.Name,
								"code":       option.Code,
								"rightLabel": option.RightLabel,
								"value":      option.Value,
								"isRequired": option.IsRequired,
								"size":       option.Size,
								"comment":    option.Comment,
								"options":    option.Options,
							}
						}
					}

					return maps.Map{}
				}),
			})
		}
	}

	this.Data["checkpoints"] = checkpointList

	this.Data["operators"] = lists.Map(teawaf.AllRuleOperators, func(k int, v interface{}) interface{} {
		def := v.(*teawaf.RuleOperatorDefinition)
		return maps.Map{
			"name":        def.Name,
			"code":        def.Code,
			"description": def.Description,
			"case":        def.CaseInsensitive,
		}
	})

	this.Data["actions"] = lists.Map(teawaf.AllActions, func(k int, v interface{}) interface{} {
		def := v.(*teawaf.ActionDefinition)
		return maps.Map{
			"name":        def.Name,
			"description": def.Description,
			"code":        def.Code,
		}
	})

	this.Show()
}

// 提交测试或者保存
func (this *RuleUpdateAction) RunPost(params struct {
	WafId   string
	GroupId string
	SetId   string

	Name string

	RulePrefixes  []string
	RuleParams    []string
	RuleOperators []string
	RuleValues    []string
	RuleCases     []int
	RuleOptions   []string

	Connector string
	Action    string

	Test         bool
	TestPrefixes []string
	TestParams   []string
	TestValues   []string

	Must *actions.Must
}) {
	// waf
	wafList := teaconfigs.SharedWAFList()
	waf := wafList.FindWAF(params.WafId)
	if waf == nil {
		this.Fail("找不到WAF")
	}

	group := waf.FindRuleGroup(params.GroupId)
	if group == nil {
		this.Fail("找不到Group")
	}

	set := group.FindRuleSet(params.SetId)
	if set == nil {
		this.Fail("找不到规则集")
	}

	set.Name = params.Name
	set.ActionOptions = maps.Map{}
	set.Rules = []*teawaf.Rule{}
	for index, prefix := range params.RulePrefixes {
		if index < len(params.RuleParams) && index < len(params.RuleOperators) && index < len(params.RuleValues) && index < len(params.RuleCases) {
			rule := teawaf.NewRule()
			rule.Operator = params.RuleOperators[index]

			param := params.RuleParams[index]
			if len(param) > 0 {
				rule.Param = "${" + prefix + "." + param + "}"
			} else {
				rule.Param = "${" + prefix + "}"
			}
			rule.Value = params.RuleValues[index]
			rule.IsCaseInsensitive = params.RuleCases[index] == 1

			// 选项
			options := params.RuleOptions[index]
			if len(options) > 0 {
				arr := []maps.Map{}
				err := json.Unmarshal([]byte(options), &arr)
				if err != nil {
					logs.Error(err)
				} else {
					rule.CheckpointOptions = map[string]string{}
					for _, m := range arr {
						code := m.GetString("code")
						value := m.GetString("value")
						rule.CheckpointOptions[code] = value
					}
				}
			}

			// 校验
			err := rule.Init()
			if err != nil {
				this.Fail("校验规则 '" + rule.Param + " " + rule.Operator + " " + rule.Value + "' 失败，原因：" + err.Error())
			}

			set.AddRule(rule)
		}
	}
	set.Connector = params.Connector

	// action
	set.Action = params.Action
	set.ActionOptions = maps.Map{}
	for k, v := range this.ParamsMap {
		if len(v) == 0 {
			continue
		}
		index := strings.Index(k, "action_")
		if index > -1 {
			set.ActionOptions[k[len("action_"):]] = v[0]
		}
	}

	// 测试
	if params.Test {
		err := set.Init()
		if err != nil {
			this.Fail("校验错误：" + err.Error())
		}

		matchedIndex := -1
		breakIndex := -1
		matchLogs := []string{"start matching ...", "==="}
	Loop:
		for index, prefix := range params.TestPrefixes {
			if index < len(params.TestParams) && index < len(params.TestValues) {
				param := ""
				if len(params.TestParams[index]) == 0 {
					param = "${" + prefix + "}"
				} else {
					param = "${" + prefix + "." + params.TestParams[index] + "}"
				}

				breakIndex = index

				for _, rule := range set.Rules {
					if rule.Param == param {
						value := params.TestValues[index]
						if rule.Test(value) {
							matchLogs = append(matchLogs, "rule: "+rule.Param+" "+rule.Operator+" "+rule.Value+"\ncompare: "+value+"\nresult:true")

							if set.Connector == teawaf.RuleConnectorOr {
								matchedIndex = index
								break Loop
							}

							if set.Connector == teawaf.RuleConnectorAnd {
								matchedIndex = index
							}
						} else {
							matchLogs = append(matchLogs, "rule: "+rule.Param+" "+rule.Operator+" "+rule.Value+"\ncompare: "+value+"\nresult:false")

							if set.Connector == teawaf.RuleConnectorAnd {
								matchedIndex = -1
								break Loop
							}
						}
					}
				}
			}
		}

		this.Data["matchedIndex"] = matchedIndex
		this.Data["breakIndex"] = breakIndex
		this.Data["matchLogs"] = matchLogs
		this.Success()
	}

	// 保存
	params.Must.
		Field("name", params.Name).
		Require("请输入规则集名称")

	err := wafList.SaveWAF(waf)
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 通知刷新
	if wafutils.IsPolicyUsed(waf.Id) {
		proxyutils.NotifyChange()
	}

	this.Success()
}
