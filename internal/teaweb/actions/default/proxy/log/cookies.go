package log

import (
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/iwind/TeaGo/actions"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type CookiesAction actions.Action

func (this *CookiesAction) Run(params struct {
	LogId string
	Day   string
}) {
	if len(params.Day) == 0 {
		params.Day = timeutil.Format("Ymd")
	}

	accessLog, err := teadb.AccessLogDAO().FindAccessLogCookie(params.Day, params.LogId)
	if err != nil {
		this.Fail(err.Error())
	}

	if accessLog != nil {
		this.Data["cookies"] = accessLog.Cookie
		this.Data["count"] = len(accessLog.Cookie)
	} else {
		this.Data["cookies"] = map[string]string{}
		this.Data["count"] = 0
	}

	this.Success()
}
