package teaproxy

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaevents"
	"github.com/gorilla/websocket"
	"github.com/iwind/TeaGo/logs"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// 调用Websocket
func (this *Request) callWebsocket(writer *ResponseWriter) error {
	if this.backend == nil {
		err := errors.New(this.requestPath() + ": no available backends for websocket")
		logs.Error(err)
		this.addError(err)
		this.serverError(writer)
		return err
	}

	upgrader := websocket.Upgrader{
		HandshakeTimeout: this.websocket.HandshakeTimeoutDuration(),
		CheckOrigin: func(r *http.Request) bool {
			if this.websocket.AllowAllOrigins {
				return true
			}
			origin := r.Header.Get("Origin")
			if len(origin) == 0 {
				return false
			}
			return this.websocket.MatchOrigin(origin)
		},
		Subprotocols: websocket.Subprotocols(this.raw),
	}

	// 自动补充Header
	this.raw.Header.Set("Connection", "upgrade")
	if len(this.raw.Header.Get("Upgrade")) == 0 {
		this.raw.Header.Set("Upgrade", "websocket")
	}

	// 接收客户端连接
	client, err := upgrader.Upgrade(this.responseWriter.Raw(), this.raw, nil)
	if err != nil {
		logs.Error(errors.New("upgrade: " + err.Error()))
		this.addError(errors.New("upgrade: " + err.Error()))
		return err
	}
	defer func() {
		_ = client.Close()
	}()

	if this.websocket.ForwardMode == teaconfigs.WebsocketForwardModeWebsocket {
		// 判断最大连接数
		if this.backend.MaxConns > 0 && this.backend.CurrentConns >= this.backend.MaxConns {
			this.serverError(writer)
			logs.Error(errors.New("too many connections"))
			this.addError(errors.New("too many connections"))
			return nil
		}

		// 增加连接数
		this.backend.IncreaseConn()
		defer this.backend.DecreaseConn()

		// 连接后端服务器
		scheme := "ws"
		if this.backend.Scheme == "https" {
			scheme = "wss"
		}
		host := this.raw.Host
		if this.backend.HasHost() {
			host = this.Format(this.backend.Host)
		}

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

		wsURL := url.URL{
			Scheme:   scheme,
			Host:     host,
			User:     this.raw.URL.User,
			Opaque:   this.raw.URL.Opaque,
			Path:     this.raw.URL.Path,
			RawQuery: this.raw.URL.RawQuery,
		}

		// TLS通讯
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		if this.backend.Cert != nil {
			obj := this.backend.Cert.CertObject()
			if obj != nil {
				tlsConfig.InsecureSkipVerify = false
				tlsConfig.Certificates = []tls.Certificate{*obj}
				if len(this.backend.Cert.ServerName) > 0 {
					tlsConfig.ServerName = this.backend.Cert.ServerName
				}
			}
		}

		// 超时时间
		connectionTimeout := this.backend.FailTimeoutDuration()
		if connectionTimeout <= 0 {
			connectionTimeout = 15 * time.Second
		}

		backendAddr := this.backend.Address
		if this.backend.HasAddrVariables() {
			backendAddr = this.Format(backendAddr)
		}
		dialer := websocket.Dialer{
			NetDial: func(network, addr string) (conn net.Conn, err error) {
				return net.DialTimeout(network, backendAddr, connectionTimeout)
			},
			TLSClientConfig:  tlsConfig,
			HandshakeTimeout: this.backend.FailTimeoutDuration(),
			Subprotocols:     websocket.Subprotocols(this.raw),
		}
		header := http.Header{}
		for k, v := range this.raw.Header {
			if strings.HasPrefix(k, "Sec-") || k == "Upgrade" || k == "Connection" {
				continue
			}
			header[k] = v
		}

		this.setProxyHeaders(header)

		// 自定义请求Header
		for _, h := range this.requestHeaders {
			if !h.On {
				continue
			}
			if h.HasVariables() {
				header[h.Name] = []string{this.Format(h.Value)}
			} else {
				header[h.Name] = []string{h.Value}
			}
		}

		server, resp, err := dialer.Dial(wsURL.String(), header)
		if err != nil {
			writer.statusCode = http.StatusInternalServerError
			_ = client.Close()

			if server != nil {
				_ = server.Close()
			}

			errString := ""
			if resp != nil && resp.Body != nil {
				data, _ := ioutil.ReadAll(resp.Body)
				errString = strconv.Itoa(resp.StatusCode) + ": " + string(bytes.TrimSpace(data))
			}
			err1 := errors.New(err.Error() + ": " + errString)
			logs.Error(err1)
			this.addError(err1)
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

				this.websocket.SetupScheduling(false)
			}

			return err
		}
		defer func() {
			_ = server.Close()
		}()

		// 设置关闭连接的处理函数
		clientIsClosed := false
		serverIsClosed := false
		client.SetCloseHandler(func(code int, text string) error {
			if serverIsClosed {
				return nil
			}
			serverIsClosed = true
			return server.Close()
		})

		// 从客户端接收数据
		go func() {
			for {
				messageType, message, err := client.ReadMessage()
				if err != nil {
					closeErr, ok := err.(*websocket.CloseError)
					if !ok && closeErr != nil && closeErr.Code != websocket.CloseGoingAway {
						logs.Error(err)
						this.addError(err)
					}
					clientIsClosed = true
					_ = client.Close()

					// 关闭Server
					if !serverIsClosed {
						serverIsClosed = true
						_ = server.Close()
					}
					break
				}
				_ = server.WriteMessage(messageType, message)
			}
		}()

		// 从后端服务器读取数据
		for {
			messageType, message, err := server.ReadMessage()
			if err != nil {
				closeErr, ok := err.(*websocket.CloseError)
				if !ok && closeErr != nil && closeErr.Code != websocket.CloseGoingAway {
					logs.Error(err)
					this.addError(err)
				}
				serverIsClosed = true
				_ = server.Close()

				// 关闭客户端
				if !clientIsClosed {
					clientIsClosed = true
					_ = client.Close()
				}
				break
			}
			_ = client.WriteMessage(messageType, message)
		}
	} else if this.websocket.ForwardMode == teaconfigs.WebsocketForwardModeHttp {
		messageQueue := make(chan []byte, 1024)
		quit := make(chan bool)
		go func() {
		FOR:
			for {
				select {
				case message := <-messageQueue:
					{
						this.raw.Method = http.MethodPut
						responseWriter := NewResponseWriter(nil)
						responseWriter.SetBodyCopying(true)
						this.raw.Body = ioutil.NopCloser(bytes.NewReader(message))
						this.raw.Header.Del("Upgrade")
						err := this.callBackend(responseWriter)
						if err != nil {
							continue FOR
						}
						if responseWriter.StatusCode() != http.StatusOK {
							logs.Error(errors.New(this.requestURI() + ": invalid response from backend: " + fmt.Sprintf("%d", responseWriter.StatusCode()) + " " + http.StatusText(responseWriter.StatusCode())))
							this.addError(errors.New(this.requestURI() + ": invalid response from backend: " + fmt.Sprintf("%d", responseWriter.StatusCode()) + " " + http.StatusText(responseWriter.StatusCode())))
							continue FOR
						}
						_ = client.WriteMessage(websocket.TextMessage, responseWriter.Body())
					}
				case <-quit:
					break FOR
				}
			}
		}()
		for {
			messageType, message, err := client.ReadMessage()
			if err != nil {
				closeErr, ok := err.(*websocket.CloseError)
				if !ok || closeErr.Code != websocket.CloseGoingAway {
					logs.Error(err)
					this.addError(err)
				}
				quit <- true
				break
			}
			if messageType == websocket.TextMessage || messageType == websocket.BinaryMessage {
				messageQueue <- message
			}
		}
	}

	return nil
}
