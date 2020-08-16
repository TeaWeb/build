package settings

import (
	"github.com/TeaWeb/build/internal/teaconfigs/db"
	"github.com/TeaWeb/build/internal/teaweb/configs"
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

	action.Data["serverChanged"] = serverChanged

	if action.Spec.HasClassPrefix("profile") {
		action.Data["teaMenu"] = "settings.profile"
	} else {
		action.Data["teaMenu"] = "settings"
	}

	// 操作按钮
	menuGroup := utils.NewMenuGroup()
	{
		menu := menuGroup.FindMenu("operations", "[操作]")
		menu.AlwaysActive = true
		menuGroup.AlwaysMenu = menu
		menu.Index = 10000

		user := configs.SharedAdminConfig().FindActiveUser(action.Session().GetString("username"))
		if user.Granted(configs.AdminGrantAll) {
			menu.Add("管理界面", "", "/settings", action.Spec.HasClassPrefix("settings.IndexAction", "server."))
		}
		menu.Add("个人资料", "", "/settings/profile", action.Spec.HasClassPrefix("profile."))
		menu.Add("登录设置", "", "/settings/login", action.Spec.HasClassPrefix("login."))

		if user.Granted(configs.AdminGrantAll) {
			// mongodb管理
			if db.SharedDBConfig().Type == db.DBTypeMongo {
				menu.Add("MongoDB设置", "", "/settings/mongo", action.Spec.HasClassPrefix("mongo."))
			}

			// MySQL管理
			if db.SharedDBConfig().Type == db.DBTypeMySQL {
				menu.Add("MySQL设置", "", "/settings/mysql", action.Spec.HasClassPrefix("mysql."))
			}

			// mongodb管理
			if db.SharedDBConfig().Type == db.DBTypePostgres {
				menu.Add("PostgreSQL设置", "", "/settings/postgres", action.Spec.HasClassPrefix("postgres."))
			}

			// 数据库
			menu.Add("切换数据库", "", "/settings/database", action.Spec.HasClassPrefix("database."))

			// 备份
			menu.Add("备份", "", "/settings/backup", action.Spec.HasClassPrefix("backup."))
		}

		menu.Add("检查版本更新", "", "/settings/update", action.Spec.HasClassPrefix("update."))
	}

	menuGroup.Sort()
	utils.SetSubMenu(action, menuGroup)

	action.Data["teaTabbar"] = []maps.Map{}
}
