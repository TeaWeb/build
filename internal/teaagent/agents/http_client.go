package teaagents

import (
	"crypto/tls"
	"net/http"
	"time"
)

// 公用的HTTP客户端
var HTTPClient = &http.Client{
	Timeout: 60 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConnsPerHost: 5,
	},
}

// 长时间的HTTP客户端
var LongHTTPClient = &http.Client{
	Timeout: 300 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConnsPerHost: 5,
	},
}
