package helpers

import (
	"github.com/iwind/TeaGo/actions"
	"net/http"
)

type UserShouldAuth struct {
	action *actions.ActionObject
}

func (auth *UserShouldAuth) BeforeAction(actionPtr actions.ActionWrapper, paramName string) (goNext bool) {
	auth.action = actionPtr.Object()
	return true
}

// 存储用户名到SESSION
func (auth *UserShouldAuth) StoreUsername(username string, remember bool) {
	// 修改sid的时间
	if remember {
		cookie := &http.Cookie{
			Name:   "sid",
			Value:  auth.action.Session().Sid,
			Path:   "/",
			MaxAge: 14 * 86400,
		}
		auth.action.AddCookie(cookie)
	} else {
		cookie := &http.Cookie{
			Name:   "sid",
			Value:  auth.action.Session().Sid,
			Path:   "/",
			MaxAge: 0,
		}
		auth.action.AddCookie(cookie)
	}
	auth.action.Session().Write("username", username)
}

func (auth *UserShouldAuth) Logout() {
	auth.action.Session().Delete()
}
