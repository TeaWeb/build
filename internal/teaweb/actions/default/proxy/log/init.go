package log

import (
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/TeaWeb/build/internal/teaweb/helpers"
	"github.com/iwind/TeaGo"
	"time"
)

// 是否正在获取请求内容
var requestBodyFetching = false
var requestBodyTime = time.Now()

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantLog,
			}).
			Helper(new(Helper)).
			Prefix("/proxy/log").
			Get("", new(IndexAction)).
			Get("/list", new(ListAction)).
			Get("/responseHeader/:logId/:day", new(ResponseHeaderAction)).
			Get("/requestHeader/:logId/:day", new(RequestHeaderAction)).
			Get("/responseBody/:logId/:day", new(ResponseBodyAction)).
			Get("/cookies/:logId/:day", new(CookiesAction)).
			Get("/history", new(HistoryAction)).
			GetPost("/day", new(DayAction)).
			EndAll()
	})

	// 请求Hook
	hook := &teaproxy.RequestHook{
		BeforeRequest: ProcessBeforeRequest,
		AfterRequest:  nil,
	}
	teaproxy.AddRequestHook(hook)
}

// 请求Hook
func ProcessBeforeRequest(req *teaproxy.Request, writer *teaproxy.ResponseWriter) bool {
	if requestBodyFetching && time.Since(requestBodyTime).Seconds() < 5 {
		req.SetIsWatching(true)
	} else {
		requestBodyFetching = false
		req.SetIsWatching(false)
	}
	return true
}
