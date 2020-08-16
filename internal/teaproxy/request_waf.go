package teaproxy

import (
	"github.com/TeaWeb/build/internal/teawaf"
	"github.com/iwind/TeaGo/logs"
	"net/http"
)

// call request waf
func (this *Request) callWAFRequest(writer *ResponseWriter) (blocked bool) {
	if this.waf == nil {
		return
	}
	goNext, group, ruleSet, err := this.waf.MatchRequest(this.raw, writer)
	if err != nil {
		logs.Error(err)
		return
	}

	if ruleSet != nil {
		if ruleSet.Action != teawaf.ActionAllow {
			this.SetAttr("waf_action", ruleSet.Action)
			this.SetAttr("waf_group", group.Id)
			this.SetAttr("waf_ruleset", ruleSet.Id)
			this.SetAttr("waf_ruleset_name", ruleSet.Name)
			this.SetAttr("waf_id", this.waf.Id)
		}
	}

	return !goNext
}

// call response waf
func (this *Request) callWAFResponse(resp *http.Response, writer *ResponseWriter) (blocked bool) {
	if this.waf == nil {
		return
	}

	goNext, group, ruleSet, err := this.waf.MatchResponse(this.raw, resp, writer)
	if err != nil {
		logs.Error(err)
		return
	}

	if ruleSet != nil {
		if ruleSet.Action != teawaf.ActionAllow {
			this.SetAttr("waf_action", ruleSet.Action)
			this.SetAttr("waf_group", group.Id)
			this.SetAttr("waf_ruleset", ruleSet.Id)
			this.SetAttr("waf_ruleset_name", ruleSet.Name)
			this.SetAttr("waf_id", this.waf.Id)
		}
	}

	return !goNext
}
