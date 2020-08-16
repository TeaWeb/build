package waf

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"net/http"
)

type Helper struct {
}

// 相关Helper
func (this *Helper) BeforeAction(action *actions.ActionObject) {
	if action.Request.Method == http.MethodGet {
		proxyutils.AddServerMenu(action, false)

		action.Data["selectedMenu"] = "list"
		if action.HasPrefix("/proxy/waf/add") {
			action.Data["selectedMenu"] = "add"
		}

		action.Data["selectedSubMenu"] = "detail"
		if action.HasPrefix("/proxy/waf/rules", "/proxy/waf/group") {
			action.Data["selectedSubMenu"] = "rules"
		} else if action.HasPrefix("/proxy/waf/test") {
			action.Data["selectedSubMenu"] = "test"
		} else if action.HasPrefix("/proxy/waf/export") {
			action.Data["selectedSubMenu"] = "export"
		} else if action.HasPrefix("/proxy/waf/import") {
			action.Data["selectedSubMenu"] = "import"
		} else if action.HasPrefix("/proxy/waf/history") || action.HasPrefix("/proxy/waf/day") {
			action.Data["selectedSubMenu"] = "history"
		}

		action.Data["inbound"] = false
		action.Data["outbound"] = false
	}
}
