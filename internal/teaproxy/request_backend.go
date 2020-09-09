package teaproxy

import (
	"context"
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaevents"
	"github.com/iwind/TeaGo/logs"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 调用后端服务器
func (this *Request) callBackend(writer *ResponseWriter) error {
	// 是否为websocket请求
	if this.raw.Header.Get("Upgrade") == "websocket" {
		websocket := teaconfigs.NewWebsocketConfig()
		websocket.On = true
		websocket.AllowAllOrigins = true
		websocket.ForwardMode = teaconfigs.WebsocketForwardModeWebsocket
		websocket.Backends = []*teaconfigs.BackendConfig{this.backend}
		websocket.Origins = []string{}
		websocket.HandshakeTimeout = "5s"
		this.websocket = websocket
		err := websocket.Validate()
		if err != nil {
			logs.Error(err)
		}
		return this.callWebsocket(writer)
	}

	this.backend.IncreaseConn()
	defer this.backend.DecreaseConn()

	if len(this.backend.Address) == 0 {
		this.serverError(writer)
		logs.Error(errors.New("backend address should not be empty"))
		this.addError(errors.New("backend address should not be empty"))
		return nil
	}

	if this.backend.HasHost() {
		this.raw.Host = this.Format(this.backend.Host)
	}

	if len(this.raw.Host) > 0 {
		this.raw.URL.Host = this.raw.Host
	} else {
		this.raw.URL.Host = this.host
	}

	if len(this.backend.Scheme) > 0 && this.backend.Scheme != "http" && this.backend.Scheme != "ftp" {
		this.raw.URL.Scheme = this.backend.Scheme
	} else {
		this.raw.URL.Scheme = this.scheme
	}

	// new uri
	if this.backend.HasRequestURI() {
		uri := this.Format(this.backend.RequestPath())
		u, err := url.ParseRequestURI(uri)
		if err == nil {
			this.raw.URL.Path = CleanPath(u.Path)
			this.raw.URL.RawQuery = u.RawQuery

			args := this.Format(this.backend.RequestArgs())
			if len(args) > 0 {
				if len(u.RawQuery) > 0 {
					this.raw.URL.RawQuery += "&" + args
				} else {
					this.raw.URL.RawQuery += args
				}
			}
		}
	} else {
		u, err := url.ParseRequestURI(this.uri)
		if err == nil {
			this.raw.URL.Path = u.Path
			this.raw.URL.RawQuery = u.RawQuery
		}
	}

	// 设置代理相关的头部
	// 参考 https://tools.ietf.org/html/rfc7239
	this.setProxyHeaders(this.raw.Header)

	this.raw.Header.Set("Connection", "keep-alive")

	// 自定义请求Header
	if len(this.requestHeaders) > 0 {
		for _, header := range this.requestHeaders {
			if !header.On {
				continue
			}
			if header.HasVariables() {
				this.raw.Header[header.Name] = []string{this.Format(header.Value)}
			} else {
				this.raw.Header[header.Name] = []string{header.Value}
			}

			// 支持修改Host
			if header.Name == "Host" && len(header.Value) > 0 {
				this.raw.Host = header.Value
			}
		}
	}

	this.raw.RequestURI = ""

	var resp *http.Response = nil
	var err error = nil
	if this.backend.IsFTP() {
		client := SharedFTPClientPool.client(this, this.backend, this.location)
		resp, err = client.Do(this.raw)
	} else {
		client := SharedHTTPClientPool.client(this, this.backend, this.location)
		resp, err = client.Do(this.raw)
	}
	if err != nil {
		// 客户端取消请求，则不提示
		httpErr, ok := err.(*url.Error)
		if !ok || httpErr.Err != context.Canceled {
			// 如果超过最大失败次数，则下线
			if !this.backend.HasCheckURL() {
				currentFails := this.backend.IncreaseFails()
				if this.backend.MaxFails > 0 && currentFails >= this.backend.MaxFails {
					this.backend.IsDown = true
					this.backend.DownTime = time.Now()

					// 下线通知
					teaevents.Post(&teaconfigs.BackendDownEvent{
						Server:    this.server,
						Backend:   this.backend,
						Location:  this.location,
						Websocket: this.websocket,
					})

					if this.websocket != nil {
						this.websocket.SetupScheduling(false)
					} else {
						this.server.SetupScheduling(false)
					}
				}
			}

			this.serverError(writer)

			logs.Println("[proxy]'" + this.raw.URL.String() + "': " + err.Error())
			this.addError(err)
		} else {
			// 是否为客户端方面的错误
			isClientError := false
			if ok {
				if httpErr.Err == context.Canceled {
					isClientError = true
					this.addError(errors.New(httpErr.Op + " " + httpErr.URL + ": client closed the connection"))
					writer.WriteHeader(499) // 仿照nginx
				}
			}

			if !isClientError {
				this.serverError(writer)
				this.addError(err)
			}
		}
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
		return nil
	}

	// waf
	if this.waf != nil {
		if this.callWAFResponse(resp, writer) {
			err = resp.Body.Close()
			if err != nil {
				logs.Error(err)
			}
			return nil
		}
	}

	// 清除错误次数
	if resp.StatusCode >= 200 && !this.backend.HasCheckURL() {
		if !this.backend.IsDown && this.backend.CurrentFails > 0 {
			this.backend.CurrentFails = 0
		}
	}

	// 特殊页面
	if len(this.pages) > 0 && this.callPage(writer, resp.StatusCode) {
		err = resp.Body.Close()
		if err != nil {
			logs.Error(err)
		}
		return nil
	}

	// 忽略的Header
	ignoreHeaders := this.convertIgnoreHeaders()
	hasIgnoreHeaders := ignoreHeaders.Len() > 0

	// 设置Header
	hasCharset := len(this.charset) > 0
	if hasCharset {
		contentTypes, ok := resp.Header["Content-Type"]
		if ok && len(contentTypes) > 0 {
			contentType := contentTypes[0]
			if _, found := textMimeMap[contentType]; found {
				resp.Header["Content-Type"][0] = contentType + "; charset=" + this.charset
			}
		}
	}
	for k, v := range resp.Header {
		if k == "Connection" {
			continue
		}
		if hasIgnoreHeaders && ignoreHeaders.Has(strings.ToUpper(k)) {
			continue
		}
		for _, subV := range v {
			writer.Header().Add(k, subV)
		}
	}

	// 自定义响应Headers
	this.WriteResponseHeaders(writer, resp.StatusCode)

	// 响应回调
	if this.responseCallback != nil {
		this.responseCallback(writer)
	}

	// 是否需要刷新
	shouldFlush := this.raw.Header.Get("Accept") == "text/event-stream"

	// 准备
	writer.Prepare(resp.ContentLength)

	// 设置响应代码
	writer.WriteHeader(resp.StatusCode)

	pool := this.bytePool(resp.ContentLength)
	buf := pool.Get()
	if shouldFlush {
		for {
			n, readErr := resp.Body.Read(buf)
			if n > 0 {
				_, err = writer.Write(buf[:n])
				writer.Flush()
				if err != nil {
					break
				}
			}
			if readErr != nil {
				err = readErr
				break
			}
		}
	} else {
		_, err = io.CopyBuffer(writer, resp.Body, buf)
	}
	pool.Put(buf)

	err1 := resp.Body.Close()
	if err1 != nil {
		logs.Error(err1)
	}

	if err != nil {
		logs.Error(err)
		this.addError(err)
		return nil
	}

	return nil
}
