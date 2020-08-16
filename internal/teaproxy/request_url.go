package teaproxy

import (
	"errors"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/logs"
	"io"
	"net/http"
	"time"
)

func (this *Request) callURL(writer *ResponseWriter, method string, url string, host string) error {
	req, err := http.NewRequest(method, url, this.raw.Body)
	if err != nil {
		return err
	}

	// 修改Host
	if len(host) > 0 {
		req.Host = this.Format(host)
	}

	// 添加当前Header
	req.Header = this.raw.Header

	// 代理头部
	this.setProxyHeaders(req.Header)

	// 自定义请求Header
	if len(this.requestHeaders) > 0 {
		for _, header := range this.requestHeaders {
			if !header.On {
				continue
			}
			if header.HasVariables() {
				req.Header[header.Name] = []string{this.Format(header.Value)}
			} else {
				req.Header[header.Name] = []string{header.Value}
			}
		}
	}

	var client = teautils.SharedHttpClient(60 * time.Second)
	resp, err := client.Do(req)
	if err != nil {
		logs.Error(errors.New(req.URL.String() + ": " + err.Error()))
		this.addError(err)
		this.serverError(writer)
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Header
	this.WriteResponseHeaders(writer, resp.StatusCode)

	writer.AddHeaders(resp.Header)
	writer.Prepare(resp.ContentLength)

	// 设置响应代码
	writer.WriteHeader(resp.StatusCode)

	// 输出内容
	pool := this.bytePool(resp.ContentLength)
	buf := pool.Get()
	_, err = io.CopyBuffer(writer, resp.Body, buf)
	pool.Put(buf)

	return err
}
