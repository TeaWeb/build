package teawaf

import (
	"github.com/TeaWeb/build/internal/teawaf/requests"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
)

type RuleConnector = string

const (
	RuleConnectorAnd = "and"
	RuleConnectorOr  = "or"
)

type RuleSet struct {
	Id          string        `yaml:"id" json:"id"`
	Code        string        `yaml:"code" json:"code"`
	On          bool          `yaml:"on" json:"on"`
	Name        string        `yaml:"name" json:"name"`
	Description string        `yaml:"description" json:"description"`
	Rules       []*Rule       `yaml:"rules" json:"rules"`
	Connector   RuleConnector `yaml:"connector" json:"connector"` // rules connector

	Action        ActionString `yaml:"action" json:"action"`
	ActionOptions maps.Map     `yaml:"actionOptions" json:"actionOptions"` // TODO TO BE IMPLEMENTED

	hasRules bool
}

func NewRuleSet() *RuleSet {
	return &RuleSet{
		Id: rands.HexString(16),
		On: true,
	}
}

func (this *RuleSet) Init() error {
	this.hasRules = len(this.Rules) > 0
	if this.hasRules {
		for _, rule := range this.Rules {
			err := rule.Init()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (this *RuleSet) AddRule(rule ...*Rule) {
	this.Rules = append(this.Rules, rule...)
}

func (this *RuleSet) MatchRequest(req *requests.Request) (b bool, err error) {
	if !this.hasRules {
		return false, nil
	}
	switch this.Connector {
	case RuleConnectorAnd:
		for _, rule := range this.Rules {
			b1, err1 := rule.MatchRequest(req)
			if err1 != nil {
				return false, err1
			}
			if !b1 {
				return false, nil
			}
		}
		return true, nil
	case RuleConnectorOr:
		for _, rule := range this.Rules {
			b1, err1 := rule.MatchRequest(req)
			if err1 != nil {
				return false, err1
			}
			if b1 {
				return true, nil
			}
		}
	default: // same as And
		for _, rule := range this.Rules {
			b1, err1 := rule.MatchRequest(req)
			if err1 != nil {
				return false, err1
			}
			if !b1 {
				return false, nil
			}
		}
		return true, nil
	}
	return
}

func (this *RuleSet) MatchResponse(req *requests.Request, resp *requests.Response) (b bool, err error) {
	if !this.hasRules {
		return false, nil
	}
	switch this.Connector {
	case RuleConnectorAnd:
		for _, rule := range this.Rules {
			b1, err1 := rule.MatchResponse(req, resp)
			if err1 != nil {
				return false, err1
			}
			if !b1 {
				return false, nil
			}
		}
		return true, nil
	case RuleConnectorOr:
		for _, rule := range this.Rules {
			b1, err1 := rule.MatchResponse(req, resp)
			if err1 != nil {
				return false, err1
			}
			if b1 {
				return true, nil
			}
		}
	default: // same as And
		for _, rule := range this.Rules {
			b1, err1 := rule.MatchResponse(req, resp)
			if err1 != nil {
				return false, err1
			}
			if !b1 {
				return false, nil
			}
		}
		return true, nil
	}
	return
}
