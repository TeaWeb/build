package teaproxy

import (
	"bufio"
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/timers"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
	"time"
)

// 隧道连接
type TunnelConnection struct {
	conn   net.Conn
	reader *bufio.Reader
	locker *sync.Mutex

	secret          string
	isAuthenticated bool

	isClosed     bool
	closeHandler func(tunnelConn *TunnelConnection)

	ticker *teautils.Ticker
}

// 获取新对象
func NewTunnelConnection(conn net.Conn, tunnelConfig *teaconfigs.TunnelConfig) *TunnelConnection {
	tunnelConn := &TunnelConnection{
		conn:            conn,
		reader:          bufio.NewReader(conn),
		locker:          &sync.Mutex{},
		isAuthenticated: true,
	}

	if len(tunnelConfig.Secret) > 0 {
		// 认证
		tunnelConn.secret = tunnelConfig.Secret
		tunnelConn.isAuthenticated = false
		timers.Delay(5*time.Second, func(timer *time.Timer) {
			if !tunnelConn.isAuthenticated {
				tunnelConn.Close()
			}
		})

		// Ping
		tunnelConn.ticker = teautils.Every(30*time.Second, func(ticker *teautils.Ticker) {
			err := tunnelConn.Ping()
			if err != nil {
				tunnelConn.Close()
			}
		})

		tunnelConn.auth()
	}
	return tunnelConn
}

// 判断是否已认证
func (this *TunnelConnection) IsAuthenticated() bool {
	return this.isAuthenticated
}

// 发送请求
func (this *TunnelConnection) Write(req *http.Request) (*http.Response, error) {
	if !this.isAuthenticated {
		return nil, errors.New("[tunnel]not been authenticated")
	}

	if this.reader == nil {
		return nil, errors.New("[tunnel]no tunnel reader")
	}

	this.locker.Lock()

	data, err := httputil.DumpRequest(req, true)
	_, err = this.conn.Write(data)
	if err != nil {
		this.locker.Unlock()
		return nil, err
	}

	resp, err := http.ReadResponse(this.reader, req)
	if err != nil {
		this.locker.Unlock()
		return resp, err
	}
	resp.Body = &TunnelResponseBody{
		ReadCloser: resp.Body,
		locker:     this.locker,
	}
	return resp, nil
}

// 远程地址
func (this *TunnelConnection) RemoteAddr() string {
	return this.conn.RemoteAddr().String()
}

// 设置关闭回调
func (this *TunnelConnection) OnClose(handler func(tunnelConn *TunnelConnection)) {
	this.closeHandler = handler
}

// Ping客户端
func (this *TunnelConnection) Ping() error {
	if this.isClosed || !this.isAuthenticated {
		return nil
	}

	req, err := http.NewRequest(http.MethodGet, "/$$TEA/ping", nil)
	if err != nil {
		logs.Error(err)
		return nil
	}

	resp, err := this.Write(req)
	if err != nil {
		return err
	}

	resp.Body.Close()

	return nil
}

// 关闭
func (this *TunnelConnection) Close() error {
	this.isClosed = true

	err := this.conn.Close()
	if this.closeHandler != nil {
		this.closeHandler(this)
	}

	if this.ticker != nil {
		this.ticker.Stop()
	}

	return err
}

// 认证
func (this *TunnelConnection) auth() {
	data, _, err := this.reader.ReadLine()
	if err != nil {
		this.Close()
		return
	}
	this.isAuthenticated = strings.TrimSpace(string(data)) == this.secret
	if !this.isAuthenticated {
		this.Close()
		return
	}
}
