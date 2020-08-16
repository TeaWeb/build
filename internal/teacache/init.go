package teacache

import (
	"github.com/TeaWeb/build/internal/teaproxy"
)

func init() {
	hook := &teaproxy.RequestHook{
		BeforeRequest: ProcessBeforeRequest,
		AfterRequest:  ProcessAfterRequest,
	}
	teaproxy.AddRequestHook(hook)
}
