package logout

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/build/internal/teaweb/helpers"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct {
	Auth *helpers.UserShouldAuth
}) {
	params.Auth.Logout()
	this.RedirectURL("/")
}
