package teaproxy

import (
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/logs"
	"io"
	"math/rand"
	"net"
	"net/http"
	"sync"
)

// 隧道状态管理
type Tunnel struct {
	config *teaconfigs.TunnelConfig

	listener net.Listener

	locker      sync.Mutex
	connections []*TunnelConnection
}

// 获取新对象
func NewTunnel(config *teaconfigs.TunnelConfig) *Tunnel {
	return &Tunnel{
		config: config,
	}
}

// 获取当前Tunnel的ID
func (this *Tunnel) Id() string {
	return this.config.Id
}

// 启动
func (this *Tunnel) Start() error {
	if this.listener != nil {
		this.listener.Close()
		this.listener = nil
	}

	if this.config == nil {
		return errors.New("[tunnel]'config' should not be nil")
	}

	if !this.config.On {
		return errors.New("[tunnel]tunnel is not enabled")
	}

	listener, err := net.Listen("tcp", this.config.Endpoint)
	if err != nil {
		return err
	}
	this.listener = listener
	logs.Println("[tunnel]start", this.config.Endpoint)
	for {
		conn, err := this.listener.Accept()
		if err != nil {
			break
		}

		this.locker.Lock()
		tunnelConn := NewTunnelConnection(conn, this.config)
		tunnelConn.OnClose(func(tunnelConn *TunnelConnection) {
			this.locker.Lock()
			defer this.locker.Unlock()

			result := []*TunnelConnection{}
			for _, conn2 := range this.connections {
				if conn2 == tunnelConn {
					continue
				}
				result = append(result, conn2)
			}

			this.connections = result
		})

		this.connections = append(this.connections, tunnelConn)
		this.locker.Unlock()
	}

	return nil
}

// 发送请求
func (this *Tunnel) Write(req *http.Request) (resp *http.Response, err error) {
	return this.writeRequest(req, 0)
}

// 获取连接数量
func (this *Tunnel) CountConnections() int {
	this.locker.Lock()
	defer this.locker.Unlock()
	count := 0
	for _, conn := range this.connections {
		if conn.isAuthenticated {
			count++
		}
	}
	return count
}

// 发送请求，并记录尝试次数
func (this *Tunnel) writeRequest(req *http.Request, retries int) (resp *http.Response, err error) {
	this.locker.Lock()
	if len(this.connections) == 0 {
		this.locker.Unlock()
		return nil, errors.New("[tunnel]no clients to write data")
	}

	index := rand.Int() % len(this.connections)
	conn := this.connections[index]
	this.locker.Unlock()

	resp, err = conn.Write(req)
	if err != nil {
		conn.Close()
		this.removeConn(conn)

		if err == io.EOF || err == io.ErrUnexpectedEOF {
			for _, conn := range this.connections {
				conn.Close()
			}
			this.connections = []*TunnelConnection{}
		}

		this.locker.Lock()
		shouldRetry := len(this.connections) > 0 && retries < 3
		this.locker.Unlock()
		if shouldRetry {
			retries++
			return this.writeRequest(req, retries)
		}
	}

	return
}

// 关闭
func (this *Tunnel) Close() error {
	if this.listener == nil {
		return nil
	}
	err := this.listener.Close()
	this.listener = nil
	return err
}

// 移除connection
func (this *Tunnel) removeConn(conn *TunnelConnection) {
	this.locker.Lock()
	defer this.locker.Unlock()

	result := []*TunnelConnection{}
	for _, c := range this.connections {
		if c == conn {
			continue
		}
		result = append(result, c)
	}

	this.connections = result
}
