package helpers

import (
	"github.com/TeaWeb/build/internal/teacluster"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/iwind/TeaGo/actions"
	"net/http"
)

// 认证拦截
type UserMustAuth struct {
	Username string
	Grant    string
}

func (this *UserMustAuth) BeforeAction(actionPtr actions.ActionWrapper, paramName string) (goNext bool) {
	var action = actionPtr.Object()

	// 安全
	action.AddHeader("X-Frame-Options", "SAMEORIGIN")
	action.AddHeader("Content-Security-Policy", "default-src 'self' data:; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'")

	// 检查IP
	if !configs.SharedAdminConfig().AllowIP(action.RequestRemoteIP()) {
		action.ResponseWriter.WriteHeader(http.StatusForbidden)
		action.WriteString(teaconst.TeaProductName + " Access Forbidden")
		return false
	}

	var session = action.Session()
	var username = session.GetString("username")
	if len(username) == 0 {
		this.login(action)
		return false
	}

	// 检查用户是否存在
	user := configs.SharedAdminConfig().FindActiveUser(username)
	if user == nil {
		this.login(action)
		return false
	}

	if len(this.Grant) > 0 {
		if !user.Granted(this.Grant) {
			action.WriteString("Permission Denied")
			return false
		}
	}

	this.Username = username

	// 初始化内置方法
	action.ViewFunc("teaTitle", func() string {
		return action.Data["teaTitle"].(string)
	})

	// 初始化变量
	modules := []map[string]interface{}{
		/**{
			"code":     "lab",
			"menuName": "实验室",
			"icon":     "medapps",
		},**/
	}

	if user.Granted(configs.AdminGrantProxy) {
		modules = append(modules, map[string]interface{}{
			"code":     "proxy",
			"menuName": "代理",
			"icon":     "paper plane outline",
		})
	}

	if teaconst.PlusEnabled {
		if user.Granted(configs.AdminGrantApi) {
			modules = append(modules, map[string]interface{}{
				"code":     "plus.apis",
				"menuName": "API+",
				"icon":     "shekel sign",
			})
		}

		/**if user.Granted(configs.AdminGrantQ) {
			modules = append(modules, map[string]interface{}{
				"code":     "plus.q",
				"menuName": "测试小Q+",
				"icon":     "dog",
			})
		}**/
	}

	// 附加功能
	if user.Granted(configs.AdminGrantAgent) {
		modules = append(modules, map[string]interface{}{
			"code":     "agents",
			"menuName": "监控",
			"icon":     "server",
		})
	}

	if user.Granted(configs.AdminGrantPlugin) {
		modules = append(modules, map[string]interface{}{
			"code":     "plugins",
			"menuName": "插件",
			"icon":     "puzzle piece",
		})
	}

	if teaconst.PlusEnabled {
		if user.Granted(configs.AdminGrantTeam) {
			modules = append(modules, map[string]interface{}{
				"code":     "plus.team",
				"menuName": "团队+",
				"icon":     "users",
			})
		}
	}

	if user.Granted(configs.AdminGrantLog) {
		modules = append(modules, map[string]interface{}{
			"code":     "log.runtime",
			"menuName": "系统日志",
			"icon":     "history",
		})
	}

	if teaconst.PlusEnabled {
		action.Data["teaTitle"] = teaconst.TeaProductName + "+管理平台"
		action.Data["teaName"] = teaconst.TeaProductName + "+"
	} else {
		action.Data["teaTitle"] = teaconst.TeaProductName + "管理平台"
		action.Data["teaName"] = teaconst.TeaProductName
	}

	if len(user.Name) == 0 {
		action.Data["teaUsername"] = username
	} else {
		action.Data["teaUsername"] = user.Name
	}

	action.Data["teaUserAvatar"] = user.Avatar

	action.Data["teaMenu"] = ""
	action.Data["teaModules"] = modules
	action.Data["teaSubMenus"] = []map[string]interface{}{}
	action.Data["teaTabbar"] = []map[string]interface{}{}
	action.Data["teaVersion"] = teaconst.TeaVersion
	action.Data["teaIsSuper"] = user.Granted(configs.AdminGrantAll)
	action.Data["teaDemoEnabled"] = teaconst.DemoEnabled
	action.Data["teaClusterActive"] = teacluster.SharedManager.IsActive()
	action.Data["teaClusterEnabled"] = teacluster.ClusterEnabled
	action.Data["teaClusterIsChanged"] = teacluster.SharedManager.IsChanged()

	return true
}

func (this *UserMustAuth) login(action *actions.ActionObject) {
	security := configs.SharedAdminConfig().Security
	if security != nil {
		action.RedirectURL(security.NewLoginURL())
	} else {
		action.RedirectURL("/login")
	}
}
