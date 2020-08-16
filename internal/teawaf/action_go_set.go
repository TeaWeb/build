package teawaf

import (
	"github.com/TeaWeb/build/internal/teawaf/requests"
	"github.com/iwind/TeaGo/logs"
	"net/http"
)

type GoSetAction struct {
	GroupId string `yaml:"groupId" json:"groupId"`
	SetId   string `yaml:"setId" json:"setId"`
}

func (this *GoSetAction) Perform(waf *WAF, request *requests.Request, writer http.ResponseWriter) (allow bool) {
	group := waf.FindRuleGroup(this.GroupId)
	if group == nil || !group.On {
		return true
	}
	set := group.FindRuleSet(this.SetId)
	if set == nil || !set.On {
		return true
	}

	b, err := set.MatchRequest(request)
	if err != nil {
		logs.Error(err)
		return true
	}
	if !b {
		return true
	}
	actionObject := FindActionInstance(set.Action, set.ActionOptions)
	if actionObject == nil {
		return true
	}
	return actionObject.Perform(waf, request, writer)
}
