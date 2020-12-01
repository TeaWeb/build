package waf

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teawaf"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/waf/wafutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
	"gopkg.in/yaml.v3"
)

type ImportAction actions.Action

// 导入
func (this *ImportAction) RunGet(params struct {
	WafId string
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
	}

	this.Show()
}

// 提交导入
func (this *ImportAction) RunPost(params struct {
	// step1
	File *actions.File

	// step2
	WafId         string
	GroupIds      []string
	OverwriteType bool
	OverwriteName bool

	Data string
	Step string

	Must *actions.Must
}) {
	if params.Step == "file" {
		if params.File == nil {
			this.Fail("请上传要导入的规则集文件")
		}

		data, err := params.File.Read()
		if err != nil {
			this.Fail("文件读取失败：" + err.Error())
		}

		waf := &teawaf.WAF{}
		err = yaml.Unmarshal(data, waf)
		if err != nil {
			this.Fail("文件内容分析失败：" + err.Error())
		}

		this.Data["data"] = string(data)
		groups := []maps.Map{}
		for _, group := range waf.Inbound {
			groups = append(groups, maps.Map{
				"id":        group.Id,
				"name":      "[入站]" + group.Name,
				"countSets": len(group.RuleSets),
			})
		}
		for _, group := range waf.Outbound {
			groups = append(groups, maps.Map{
				"id":        group.Id,
				"name":      "[出站]" + group.Name,
				"countSets": len(group.RuleSets),
			})
		}
		this.Data["groups"] = groups

		this.Success()
	} else if params.Step == "groups" { // 提交分组信息
		waf := &teawaf.WAF{}
		err := yaml.Unmarshal([]byte(params.Data), waf)
		if err != nil {
			this.Fail("文件内容分析失败：" + err.Error())
		}

		if len(params.GroupIds) == 0 {
			this.Fail("请选择要导入的规则分组")
		}

		wafList := teaconfigs.SharedWAFList()
		currentWAF := wafList.FindWAF(params.WafId)
		if currentWAF == nil {
			this.Fail("找不到当前的WAF")
		}

		countGroups := 0
		countSets := 0
		for _, groupId := range params.GroupIds {
			group := waf.FindRuleGroup(groupId)
			if group == nil {
				continue
			}

			// 删除老的分组
			if params.OverwriteType && len(group.Code) > 0 {
				oldGroup := currentWAF.FindRuleGroupWithCode(group.Code)
				if oldGroup != nil {
					currentWAF.RemoveRuleGroup(oldGroup.Id)
				}
			} else if params.OverwriteName && len(group.Name) > 0 {
				if group.IsInbound {
					oldGroup := this.findInboundGroupWithName(currentWAF, group.Name)
					if oldGroup != nil {
						currentWAF.RemoveRuleGroup(oldGroup.Id)
					}
				} else {
					oldGroup := this.findOutboundGroupWithName(currentWAF, group.Name)
					if oldGroup != nil {
						currentWAF.RemoveRuleGroup(oldGroup.Id)
					}
				}
			}

			// 添加新的分组
			countGroups++
			countSets += len(group.RuleSets)
			group.Id = rands.HexString(16) // 重新生成ID，避免和现有的ID冲突
			currentWAF.AddRuleGroup(group)
		}

		err = wafList.SaveWAF(currentWAF)
		if err != nil {
			this.Fail("保存失败：" + err.Error())
		}

		this.Data["countGroups"] = countGroups
		this.Data["countSets"] = countSets

		// 通知刷新
		if wafutils.IsPolicyUsed(currentWAF.Id) {
			proxyutils.NotifyChange()
		}

		this.Success()
	}
}

func (this *ImportAction) findInboundGroupWithName(waf *teawaf.WAF, name string) *teawaf.RuleGroup {
	for _, g := range waf.Inbound {
		if g.Name == name {
			return g
		}
	}
	return nil
}

func (this *ImportAction) findOutboundGroupWithName(waf *teawaf.WAF, name string) *teawaf.RuleGroup {
	for _, g := range waf.Outbound {
		if g.Name == name {
			return g
		}
	}
	return nil
}
