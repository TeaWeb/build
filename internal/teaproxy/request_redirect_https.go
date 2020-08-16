package teaproxy

import (
	"net"
	"net/http"
)

func (this *Request) callRedirectToHttps(writer *ResponseWriter) {
	// 是否需要跳转到HTTPS
	if this.redirectToHttps && this.rawScheme == "http" {
		host := this.raw.Host

		host1, _, err := net.SplitHostPort(host)
		if err == nil {
			host = host1
		}

		// 是否有HTTPS
		if this.server.SSL != nil && this.server.SSL.On && len(this.server.SSL.Listen) > 0 {
			listen := this.server.SSL.Listen[0]
			_, port, err := net.SplitHostPort(listen)
			if err == nil {
				if port == "443" {
					u := "https://" + host + this.raw.RequestURI
					http.Redirect(writer, this.raw, u, http.StatusMovedPermanently)
					return
				} else {
					u := "https://" + host + ":" + port + this.raw.RequestURI
					http.Redirect(writer, this.raw, u, http.StatusMovedPermanently)
					return
				}
			} else {
				u := "https://" + host + this.raw.RequestURI
				http.Redirect(writer, this.raw, u, http.StatusMovedPermanently)
				return
			}
		}

		u := "https://" + host + this.raw.RequestURI
		http.Redirect(writer, this.raw, u, http.StatusMovedPermanently)
		return
	}
}
