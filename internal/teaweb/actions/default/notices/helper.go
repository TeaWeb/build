package notices

import (
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teaweb/utils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"net/http"
)

type Helper struct {
}

func (this *Helper) BeforeAction(actionPtr actions.ActionWrapper) {
	action := actionPtr.Object()
	action.Data["teaMenu"] = "notices"

	if action.Request.Method == http.MethodGet {
		if !action.HasPrefix("/notices/badge") {
			count, err := teadb.NoticeDAO().CountAllUnreadNotices()
			if err != nil {
				logs.Error(err)
			}
			action.Data["countUnread"] = count
		}
	}

	// 操作按钮
	menuGroup := utils.NewMenuGroup()
	{
		menu := menuGroup.FindMenu("operations", "[操作]")
		menu.AlwaysActive = true
		menuGroup.AlwaysMenu = menu
		menu.Index = 10000
		menu.Add("通知", "", "/notices", true)
	}

	menuGroup.Sort()
	utils.SetSubMenu(action, menuGroup)
}
