package plugins

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

	action.Data["teaMenu"] = "plugins"
	action.Data["teaTabbar"] = []maps.Map{}

	// 操作按钮
	menuGroup := utils.NewMenuGroup()
	{
		menu := menuGroup.FindMenu("operations", "[操作]")
		menu.AlwaysActive = true
		menuGroup.AlwaysMenu = menu
		menu.Index = 10000
		menu.Add("已安装插件", "", "/plugins", action.HasPrefix("/plugins"))
	}

	menuGroup.Sort()
	utils.SetSubMenu(action, menuGroup)
}
