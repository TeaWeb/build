package teaproxy

// 调用代理
func (this *Request) callProxy(writer *ResponseWriter) error {
	backend := this.proxy.NextBackend(this.backendCall)

	if len(this.backendCall.ResponseCallbacks) > 0 {
		this.responseCallback = this.backendCall.CallResponseCallbacks
	}

	this.backend = backend
	return this.callBackend(writer)
}
