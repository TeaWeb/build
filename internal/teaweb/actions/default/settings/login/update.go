package login

import (
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/iwind/TeaGo/actions"
)

type UpdateAction actions.Action

func (this *UpdateAction) Run(params struct{}) {
	username := this.Session().GetString("username")
	this.Data["username"] = username

	this.Show()
}

func (this *UpdateAction) RunPost(params struct {
	Username  string
	Password  string
	Password2 string
	Must      *actions.Must
}) {
	params.Must.
		Field("username", params.Username).
		Require("请输入登录用户名").
		Match("^[a-zA-Z0-9]{1,20}$", "用户名只能包含英文字母、数字")

	config := configs.SharedAdminConfig()

	username := this.Session().GetString("username")
	if username != params.Username && config.ContainsActiveUser(params.Username) {
		this.FailField("username", "此用户名已经被使用，请换一个")
	}

	if len(params.Password) > 0 {
		params.Must.
			Field("password", params.Password).
			Match("^[a-zA-Z0-9]{1,20}$", "密码只能包含英文字母、数字").
			Equal(params.Password2, "两次输入的密码不一致")
	}

	var found = false
	for _, user := range config.Users {
		if user.Username == username {
			user.Username = params.Username

			if len(params.Password) > 0 {
				user.Password = config.EncryptPassword(params.Password)
			}

			found = true
		}
	}

	if !found {
		this.RedirectURL("/logout")
		return
	}

	err := config.Save()
	if err != nil {
		this.Fail("文件保存失败，请检查'configs/admin.conf'文件的写入权限")
	}

	// 修改SESSION中的username
	this.Session().Write("username", params.Username)

	this.Next("/settings/login", nil).Success("保存成功")
}
