package proxy

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"net/http"
)

type Helper struct {
}

func (this *Helper) BeforeAction(action *actions.ActionObject) {
	if action.Request.Method == http.MethodGet && !action.HasPrefix("/proxy/status") {
		if action.HasPrefix("/proxy/add") {
			proxyutils.AddServerMenu(action, false)
		} else {
			proxyutils.AddServerMenu(action, true)
		}
	}
}
