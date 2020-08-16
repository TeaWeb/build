package profile

import (
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/time"
	"time"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct{}) {
	username := this.Session().GetString("username")
	user := configs.SharedAdminConfig().FindActiveUser(username)

	userMap := maps.Map{
		"name": user.Name,
		"tel":  user.Tel,
	}

	if user.CreatedAt > 0 {
		userMap["createdTime"] = timeutil.Format("Y-m-d H:i:s", time.Unix(user.CreatedAt, 0))
	} else {
		userMap["createdTime"] = "-"
	}

	// 登录
	if user.LoggedAt > 0 {
		userMap["loggedTime"] = timeutil.Format("Y-m-d H:i:s", time.Unix(user.LoggedAt, 0))

		if len(user.LoggedIP) > 0 {
			userMap["loggedIP"] = user.LoggedIP
		} else {
			userMap["loggedIP"] = "-"
		}
	} else {
		userMap["loggedTime"] = "-"
		userMap["loggedIP"] = "-"
	}

	// 头像
	userMap["avatar"] = user.Avatar

	this.Data["user"] = userMap

	this.Show()
}
