package actionutils

import (
	"github.com/iwind/TeaGo/actions"
	"net/http"
)

type ParentAction struct {
	actions.ActionObject
}

func (this *ParentAction) MainMenu(menuItem string) {
	this.Data["mainMenu"] = menuItem
}

func (this *ParentAction) SecondMenu(menuItem string) {
	this.Data["secondMenu"] = menuItem
}

func (this *ParentAction) ErrorPage(message string) {
	if this.Request.Method == http.MethodGet {
		this.WriteString(message)
	} else {
		this.Fail(message)
	}
}
