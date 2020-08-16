package login

import (
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/iwind/TeaGo/actions"
	"strings"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct{}) {
	username := this.Session().GetString("username")
	this.Data["username"] = username
	this.Data["passwordMask"] = ""

	config := configs.SharedAdminConfig()

	user := config.FindUser(username)
	if user == nil {
		this.RedirectURL("/logout")
		return
	}

	encryptedType := ""
	if strings.HasPrefix(user.Password, "md5:") {
		encryptedType = "md5:"
	}
	this.Data["passwordMask"] = encryptedType + strings.Repeat("*", len(user.Password))
	this.Data["key"] = user.Key

	this.Show()
}
