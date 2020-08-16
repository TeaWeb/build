package log

import (
	"github.com/TeaWeb/build/internal/teaweb/utils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"net/http"
)

type Helper struct {
}

func (this *Helper) BeforeAction(action *actions.ActionObject) {
	if action.Request.Method != http.MethodGet {
		return
	}

	action.Data["teaTabbar"] = []maps.Map{}
	action.Data["teaMenu"] = "log.runtime"

	// 操作按钮
	menuGroup := utils.NewMenuGroup()
	{
		menu := menuGroup.FindMenu("operations", "[操作]")
		menu.AlwaysActive = true
		menuGroup.AlwaysMenu = menu
		menu.Index = 10000
		menu.Add("系统日志", "", "/log/runtime", action.HasPrefix("/log/runtime"))
		menu.Add("操作日志", "", "/log/audit", action.HasPrefix("/log/audit"))
	}

	menuGroup.Sort()
	utils.SetSubMenu(action, menuGroup)
}
