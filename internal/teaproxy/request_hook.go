package teaproxy

// 外部Hook
var requestHooks = []*RequestHook{}

// 请求Hook定义
type RequestHook struct {
	BeforeRequest func(req *Request, writer *ResponseWriter) (goNext bool)
	AfterRequest  func(req *Request, writer *ResponseWriter) (goNext bool)
}

// 添加Hook
func AddRequestHook(hook *RequestHook) {
	requestHooks = append(requestHooks, hook)
}

// 执行Before Hook
func CallRequestBeforeHook(req *Request, writer *ResponseWriter) (goNext bool) {
	for _, hook := range requestHooks {
		if hook.BeforeRequest == nil {
			continue
		}
		b := hook.BeforeRequest(req, writer)
		if !b {
			return false
		}
	}
	return true
}

// 执行After Hook
func CallRequestAfterHook(req *Request, writer *ResponseWriter) (goNext bool) {
	for i := len(requestHooks) - 1; i >= 0; i -- {
		hook := requestHooks[i]
		if hook.AfterRequest == nil {
			continue
		}
		b := hook.AfterRequest(req, writer)
		if !b {
			return false
		}
	}
	return true
}
