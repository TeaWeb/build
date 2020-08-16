package v1

import (
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
)

type Helper struct {
}

func (this *Helper) BeforeAction(action actions.ActionWrapper) bool {
	if teaconst.DemoEnabled {
		action.Object().Fail("can not call api under demo mode")
		return false
	}
	apiutils.ValidateUser(action)
	return true
}
