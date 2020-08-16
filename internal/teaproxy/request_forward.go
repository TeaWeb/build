package teaproxy

import (
	"crypto/tls"
	"net/http"
	"sync"
)

var mitmCache = map[string]*tls.Config{} // host => config
var mitmLocker = &sync.Mutex{}

// 正向代理
// TODO 支持WebSocket
func (this *Request) Forward(writer *ResponseWriter) error {
	if len(this.raw.URL.Scheme) == 0 {
		this.rawScheme = "https"
	}

	this.setProxyHeaders(this.raw.Header)

	enableMitm := this.server.ForwardHTTP.EnableMITM

	proxy := &ForwardProxy{
		req:    this,
		writer: writer,
	}
	if this.method == http.MethodConnect { // connect
		if enableMitm {
			return proxy.forwardMitm()
		} else {
			return proxy.forwardConnect()
		}

	} else { // http
		return proxy.forwardHTTP()
	}
}
