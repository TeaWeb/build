package log

import (
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/iwind/TeaGo/actions"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type ResponseHeaderAction actions.Action

// 响应Header
func (this *ResponseHeaderAction) Run(params struct {
	LogId string
	Day   string
}) {
	if len(params.Day) == 0 {
		params.Day = timeutil.Format("Ymd")
	}
	accessLog, err := teadb.AccessLogDAO().FindResponseHeaderAndBody(params.Day, params.LogId)
	if err != nil {
		this.Fail(err.Error())
	}
	if accessLog != nil {
		this.Data["headers"] = accessLog.SentHeader
		this.Data["hasBody"] = len(accessLog.ResponseBodyData) > 0
	} else {
		this.Data["headers"] = map[string][]string{}
		this.Data["hasBody"] = false
	}

	this.Success()
}
