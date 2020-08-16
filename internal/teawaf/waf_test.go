package teawaf

import (
	"github.com/iwind/TeaGo/assert"
	"net/http"
	"testing"
)

func TestWAF_MatchRequest(t *testing.T) {
	a := assert.NewAssertion(t)

	set := NewRuleSet()
	set.Name = "Name_Age"
	set.Connector = RuleConnectorAnd
	set.Rules = []*Rule{
		{
			Param:    "${arg.name}",
			Operator: RuleOperatorEqString,
			Value:    "lu",
		},
		{
			Param:    "${arg.age}",
			Operator: RuleOperatorEq,
			Value:    "20",
		},
	}
	set.Action = ActionBlock

	group := NewRuleGroup()
	group.AddRuleSet(set)
	group.IsInbound = true

	waf := NewWAF()
	waf.AddRuleGroup(group)
	err := waf.Init()
	if err != nil {
		t.Fatal(err)
	}

	waf.OnAction(func(action ActionString) (goNext bool) {
		return action != ActionBlock
	})

	req, err := http.NewRequest(http.MethodGet, "http://teaos.cn/hello?name=lu&age=20", nil)
	if err != nil {
		t.Fatal(err)
	}
	goNext, _, set, err := waf.MatchRequest(req, nil)
	if err != nil {
		t.Fatal(err)
	}
	if set == nil {
		t.Log("not match")
		return
	}
	t.Log("goNext:", goNext, "set:", set.Name)
	a.IsFalse(goNext)
}
