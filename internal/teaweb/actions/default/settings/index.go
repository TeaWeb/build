package settings

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct{}) {
	this.Data["error"] = ""

	server, err := teaconfigs.LoadWebConfig()
	if err != nil {
		this.Data["error"] = "读取配置错误：" + err.Error()
		this.Show()
		return
	}

	this.Data["server"] = server

	// admin
	admin := configs.SharedAdminConfig()
	if admin.Security == nil {
		admin.Security = configs.NewAdminSecurity()
	}
	this.Data["security"] = admin.Security
	this.Data["passwordEncryptTypeText"] = admin.Security.PasswordEncryptTypeText()

	this.Show()
}
