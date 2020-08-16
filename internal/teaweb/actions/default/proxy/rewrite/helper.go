package rewrite

import (
	"github.com/iwind/TeaGo/actions"
	"net/http"
)

type Helper struct {
}

func (this *Helper) BeforeAction(action *actions.ActionObject) {
	if action.Request.Method != http.MethodGet {
		return
	}

	if action.HasParam("locationId") {
		action.Data["selectedTab"] = "location"
	}
}
