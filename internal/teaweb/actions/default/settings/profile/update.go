package profile

import (
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/time"
	"time"
)

type UpdateAction actions.Action

func (this *UpdateAction) Run(params struct{}) {
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

func (this *UpdateAction) RunPost(params struct {
	Name string
	Tel  string
	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入你的姓名")

	username := this.Session().GetString("username")

	adminConfig := configs.SharedAdminConfig()
	user := adminConfig.FindActiveUser(username)
	user.Name = params.Name
	user.Tel = params.Tel
	adminConfig.Save()

	this.Next("/settings/profile", nil).Success("保存成功")
}
