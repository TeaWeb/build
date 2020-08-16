package checkpoints

import (
	"github.com/TeaWeb/build/internal/teawaf/requests"
	"net/url"
	"strings"
)

type RequestCookiesCheckpoint struct {
	Checkpoint
}

func (this *RequestCookiesCheckpoint) RequestValue(req *requests.Request, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	var cookies = []string{}
	for _, cookie := range req.Cookies() {
		cookies = append(cookies, url.QueryEscape(cookie.Name)+"="+url.QueryEscape(cookie.Value))
	}
	value = strings.Join(cookies, "&")
	return
}

func (this *RequestCookiesCheckpoint) ResponseValue(req *requests.Request, resp *requests.Response, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param, options)
	}
	return
}
