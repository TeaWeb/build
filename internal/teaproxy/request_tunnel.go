package teaproxy

import (
	"errors"
	"github.com/iwind/TeaGo/logs"
	"io"
	"strings"
)

func (this *Request) callTunnel(writer *ResponseWriter) error {
	if this.tunnel == nil {
		return errors.New("tunnel config should not be nil")
	}

	tunnel := SharedTunnelManager.FindTunnel(this.server.Id, this.tunnel.Id)
	if tunnel == nil {
		return errors.New("tunnel should not be nil")
	}

	this.setProxyHeaders(this.raw.Header)

	resp, err := tunnel.Write(this.raw)
	if err != nil {
		this.serverError(writer)
		this.addError(err)
		logs.Println("[tunnel]\""+this.raw.RequestURI+"\":\n", err.Error())
		return err
	}
	defer resp.Body.Close()

	// 忽略的Header
	ignoreHeaders := this.convertIgnoreHeaders()
	hasIgnoreHeaders := ignoreHeaders.Len() > 0

	// 设置Header
	hasCharset := len(this.charset) > 0
	for k, v := range resp.Header {
		if k == "Connection" {
			continue
		}
		if hasIgnoreHeaders && ignoreHeaders.Has(strings.ToUpper(k)) {
			continue
		}
		for _, subV := range v {
			// 字符集
			if hasCharset && k == "Content-Type" {
				if _, found := textMimeMap[subV]; found {
					if !strings.Contains(subV, "charset=") {
						subV += "; charset=" + this.charset
					}
				}
			}

			writer.Header().Add(k, subV)
		}
	}

	// 自定义响应Headers
	this.WriteResponseHeaders(writer, resp.StatusCode)

	// 准备
	writer.Prepare(resp.ContentLength)

	// 设置响应代码
	writer.WriteHeader(resp.StatusCode)

	pool := this.bytePool(resp.ContentLength)
	buf := pool.Get()
	_, err = io.CopyBuffer(writer, resp.Body, buf)
	pool.Put(buf)
	if err != nil {
		this.addError(err)
		logs.Println("[tunnel]write response: " + err.Error())
		return err
	}

	return nil
}
