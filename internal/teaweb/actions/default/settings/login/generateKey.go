package login

import (
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/iwind/TeaGo/actions"
)

type GenerateKeyAction actions.Action

// 为当前用户生成Key
func (this *GenerateKeyAction) Run(params struct{}) {
	username := this.Session().GetString("username")

	config := configs.SharedAdminConfig()
	user := config.FindUser(username)
	if user == nil {
		this.Fail("登录信息错误，请重新登录")
	}

	user.Key = user.GenerateKey()
	err := config.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
