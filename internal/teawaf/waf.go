package teawaf

import (
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/build/internal/teawaf/checkpoints"
	"github.com/TeaWeb/build/internal/teawaf/requests"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/rands"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"reflect"
)

type WAF struct {
	Id             string                `yaml:"id" json:"id"`
	On             bool                  `yaml:"on" json:"on"`
	Name           string                `yaml:"name" json:"name"`
	Inbound        []*RuleGroup          `yaml:"inbound" json:"inbound"`
	Outbound       []*RuleGroup          `yaml:"outbound" json:"outbound"`
	CreatedVersion string                `yaml:"createdVersion" json:"createdVersion"`
	Cond           []*shared.RequestCond `yaml:"cond" json:"cond"`

	ActionBlock *BlockAction `yaml:"actionBlock" json:"actionBlock"` // action block config

	IPTables []*IPTable `yaml:"ipTables" json:"ipTables"` // IP table list

	hasInboundRules  bool
	hasOutboundRules bool
	onActionCallback func(action ActionString) (goNext bool)

	checkpointsMap map[string]checkpoints.CheckpointInterface // prefix => checkpoint
}

func NewWAF() *WAF {
	return &WAF{
		Id: rands.HexString(16),
		On: true,
	}
}

func NewWAFFromFile(path string) (waf *WAF, err error) {
	if len(path) == 0 {
		return nil, errors.New("'path' should not be empty")
	}
	file := files.NewFile(path)
	if !file.Exists() {
		return nil, errors.New("'" + path + "' not exist")
	}

	reader, err := file.Reader()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	waf = &WAF{}
	err = reader.ReadYAML(waf)
	if err != nil {
		return nil, err
	}
	return waf, nil
}

func (this *WAF) Init() error {
	// checkpoint
	this.checkpointsMap = map[string]checkpoints.CheckpointInterface{}
	for _, def := range checkpoints.AllCheckpoints {
		instance := reflect.New(reflect.Indirect(reflect.ValueOf(def.Instance)).Type()).Interface().(checkpoints.CheckpointInterface)
		instance.Init()
		this.checkpointsMap[def.Prefix] = instance
	}

	// rules
	this.hasInboundRules = len(this.Inbound) > 0
	this.hasOutboundRules = len(this.Outbound) > 0

	if this.hasInboundRules {
		for _, group := range this.Inbound {
			// finder
			for _, set := range group.RuleSets {
				for _, rule := range set.Rules {
					rule.SetCheckpointFinder(this.FindCheckpointInstance)
				}
			}

			err := group.Init()
			if err != nil {
				return err
			}
		}
	}

	if this.hasOutboundRules {
		for _, group := range this.Outbound {
			// finder
			for _, set := range group.RuleSets {
				for _, rule := range set.Rules {
					rule.SetCheckpointFinder(this.FindCheckpointInstance)
				}
			}

			err := group.Init()
			if err != nil {
				return err
			}
		}
	}

	// validate conds
	if len(this.Cond) > 0 {
		for _, cond := range this.Cond {
			err := cond.Validate()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *WAF) AddRuleGroup(ruleGroup *RuleGroup) {
	if ruleGroup.IsInbound {
		this.Inbound = append(this.Inbound, ruleGroup)
	} else {
		this.Outbound = append(this.Outbound, ruleGroup)
	}
}

func (this *WAF) RemoveRuleGroup(ruleGroupId string) {
	if len(ruleGroupId) == 0 {
		return
	}

	{
		result := []*RuleGroup{}
		for _, group := range this.Inbound {
			if group.Id == ruleGroupId {
				continue
			}
			result = append(result, group)
		}
		this.Inbound = result
	}

	{
		result := []*RuleGroup{}
		for _, group := range this.Outbound {
			if group.Id == ruleGroupId {
				continue
			}
			result = append(result, group)
		}
		this.Outbound = result
	}
}

func (this *WAF) FindRuleGroup(ruleGroupId string) *RuleGroup {
	if len(ruleGroupId) == 0 {
		return nil
	}
	for _, group := range this.Inbound {
		if group.Id == ruleGroupId {
			return group
		}
	}
	for _, group := range this.Outbound {
		if group.Id == ruleGroupId {
			return group
		}
	}
	return nil
}

func (this *WAF) FindRuleGroupWithCode(ruleGroupCode string) *RuleGroup {
	if len(ruleGroupCode) == 0 {
		return nil
	}
	for _, group := range this.Inbound {
		if group.Code == ruleGroupCode {
			return group
		}
	}
	for _, group := range this.Outbound {
		if group.Code == ruleGroupCode {
			return group
		}
	}
	return nil
}

func (this *WAF) MoveInboundRuleGroup(fromIndex int, toIndex int) {
	if fromIndex < 0 || fromIndex >= len(this.Inbound) {
		return
	}
	if toIndex < 0 || toIndex >= len(this.Inbound) {
		return
	}
	if fromIndex == toIndex {
		return
	}

	group := this.Inbound[fromIndex]
	result := []*RuleGroup{}
	for i := 0; i < len(this.Inbound); i++ {
		if i == fromIndex {
			continue
		}
		if fromIndex > toIndex && i == toIndex {
			result = append(result, group)
		}
		result = append(result, this.Inbound[i])
		if fromIndex < toIndex && i == toIndex {
			result = append(result, group)
		}
	}

	this.Inbound = result
}

func (this *WAF) MoveOutboundRuleGroup(fromIndex int, toIndex int) {
	if fromIndex < 0 || fromIndex >= len(this.Outbound) {
		return
	}
	if toIndex < 0 || toIndex >= len(this.Outbound) {
		return
	}
	if fromIndex == toIndex {
		return
	}

	group := this.Outbound[fromIndex]
	result := []*RuleGroup{}
	for i := 0; i < len(this.Outbound); i++ {
		if i == fromIndex {
			continue
		}
		if fromIndex > toIndex && i == toIndex {
			result = append(result, group)
		}
		result = append(result, this.Outbound[i])
		if fromIndex < toIndex && i == toIndex {
			result = append(result, group)
		}
	}

	this.Outbound = result
}

func (this *WAF) MatchRequest(rawReq *http.Request, writer http.ResponseWriter) (goNext bool, group *RuleGroup, set *RuleSet, err error) {
	if !this.hasInboundRules {
		return true, nil, nil, nil
	}

	req := requests.NewRequest(rawReq)

	// validate captcha
	if rawReq.URL.Path == "/WAFCAPTCHA" {
		captchaValidator.Run(req, writer)
		return
	}

	// match rules
	for _, group := range this.Inbound {
		if !group.On {
			continue
		}
		b, set, err := group.MatchRequest(req)
		if err != nil {
			return true, nil, nil, err
		}
		if b {
			if this.onActionCallback == nil {
				if set.Action == ActionBlock && this.ActionBlock != nil {
					return this.ActionBlock.Perform(this, req, writer), group, set, nil
				} else {
					actionObject := FindActionInstance(set.Action, set.ActionOptions)
					if actionObject == nil {
						return true, group, set, errors.New("no action called '" + set.Action + "'")
					}
					goNext := actionObject.Perform(this, req, writer)
					return goNext, group, set, nil
				}
			} else {
				goNext = this.onActionCallback(set.Action)
			}
			return goNext, group, set, nil
		}
	}
	return true, nil, nil, nil
}

func (this *WAF) MatchResponse(rawReq *http.Request, rawResp *http.Response, writer http.ResponseWriter) (goNext bool, group *RuleGroup, set *RuleSet, err error) {
	if !this.hasOutboundRules {
		return true, nil, nil, nil
	}
	req := requests.NewRequest(rawReq)
	resp := requests.NewResponse(rawResp)
	for _, group := range this.Outbound {
		if !group.On {
			continue
		}
		b, set, err := group.MatchResponse(req, resp)
		if err != nil {
			return true, nil, nil, err
		}
		if b {
			if this.onActionCallback == nil {
				if set.Action == ActionBlock && this.ActionBlock != nil {
					return this.ActionBlock.Perform(this, req, writer), group, set, nil
				} else {
					actionObject := FindActionInstance(set.Action, set.ActionOptions)
					if actionObject == nil {
						return true, group, set, errors.New("no action called '" + set.Action + "'")
					}
					goNext := actionObject.Perform(this, req, writer)
					return goNext, group, set, nil
				}
			} else {
				goNext = this.onActionCallback(set.Action)
			}
			return goNext, group, set, nil
		}
	}
	return true, nil, nil, nil
}

// save to file path
func (this *WAF) Save(path string) error {
	if len(path) == 0 {
		return errors.New("path should not be empty")
	}
	if len(this.CreatedVersion) == 0 {
		this.CreatedVersion = teaconst.TeaVersion
	}
	data, err := yaml.Marshal(this)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}

func (this *WAF) ContainsGroupCode(code string) bool {
	if len(code) == 0 {
		return false
	}
	for _, group := range this.Inbound {
		if group.Code == code {
			return true
		}
	}
	for _, group := range this.Outbound {
		if group.Code == code {
			return true
		}
	}
	return false
}

func (this *WAF) Copy() *WAF {
	waf := &WAF{
		Id:       this.Id,
		On:       this.On,
		Name:     this.Name,
		Inbound:  this.Inbound,
		Outbound: this.Outbound,
	}
	return waf
}

func (this *WAF) CountInboundRuleSets() int {
	count := 0
	for _, group := range this.Inbound {
		count += len(group.RuleSets)
	}
	return count
}

func (this *WAF) CountOutboundRuleSets() int {
	count := 0
	for _, group := range this.Outbound {
		count += len(group.RuleSets)
	}
	return count
}

func (this *WAF) OnAction(onActionCallback func(action ActionString) (goNext bool)) {
	this.onActionCallback = onActionCallback
}

func (this *WAF) FindCheckpointInstance(prefix string) checkpoints.CheckpointInterface {
	instance, ok := this.checkpointsMap[prefix]
	if ok {
		return instance
	}
	return nil
}

// start
func (this *WAF) Start() {
	for _, checkpoint := range this.checkpointsMap {
		checkpoint.Start()
	}
}

// call stop() when the waf was deleted
func (this *WAF) Stop() {
	for _, checkpoint := range this.checkpointsMap {
		checkpoint.Stop()
	}
}

// merge with template
func (this *WAF) MergeTemplate() (changedItems []string) {
	changedItems = []string{}

	// compare versions
	if !Tea.IsTesting() && this.CreatedVersion == teaconst.TeaVersion {
		return
	}
	this.CreatedVersion = teaconst.TeaVersion

	template := Template()
	groups := []*RuleGroup{}
	groups = append(groups, template.Inbound...)
	groups = append(groups, template.Outbound...)
	for _, group := range groups {
		oldGroup := this.FindRuleGroupWithCode(group.Code)
		if oldGroup == nil {
			group.Id = rands.HexString(16)
			this.AddRuleGroup(group)
			changedItems = append(changedItems, "+group "+group.Name)
			continue
		}

		// check rule sets
		for _, set := range group.RuleSets {
			oldSet := oldGroup.FindRuleSetWithCode(set.Code)
			if oldSet == nil {
				oldGroup.AddRuleSet(set)
				changedItems = append(changedItems, "+group "+group.Name+" rule set:"+set.Name)
			} else if len(oldSet.Rules) < len(set.Rules) {
				oldSet.Rules = set.Rules
				changedItems = append(changedItems, "*group "+group.Name+" rule set:"+set.Name)
			}
		}
	}
	return
}

// match keyword
func (this *WAF) MatchKeyword(keyword string) (matched bool, name string, tags []string) {
	if teautils.MatchKeyword(this.Name, keyword) {
		matched = true
		name = this.Name
	}
	return
}

// match cond
func (this *WAF) MatchConds(formatter func(string) string) bool {
	if len(this.Cond) == 0 {
		return true
	}
	for _, cond := range this.Cond {
		if !cond.Match(formatter) {
			return false
		}
	}
	return true
}
