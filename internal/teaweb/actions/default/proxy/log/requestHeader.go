package log

import (
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/iwind/TeaGo/actions"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type RequestHeaderAction actions.Action

// 请求Header
func (this *RequestHeaderAction) Run(params struct {
	LogId string
	Day   string
}) {
	if len(params.Day) == 0 {
		params.Day = timeutil.Format("Ymd")
	}
	accessLog, err := teadb.AccessLogDAO().FindRequestHeaderAndBody(params.Day, params.LogId)
	if err != nil {
		this.Fail(err.Error())
	}
	if accessLog != nil {
		this.Data["headers"] = accessLog.Header
		this.Data["body"] = string(accessLog.RequestData)
	} else {
		this.Data["headers"] = map[string][]string{}
		this.Data["body"] = ""
	}

	this.Success()
}
