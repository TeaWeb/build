package teaproxy

import (
	"net/http"
)

// 自定义ServeMux
type HTTPServeMux struct {
	http.ServeMux
}

func (this *HTTPServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 解决因为URL中包含多个/而自动跳转的问题
	r.URL.Path = CleanPath(r.URL.Path)

	this.ServeMux.ServeHTTP(w, r)
}
