package teaproxy

import (
	"context"
	"crypto/tls"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/logs"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// HTTP客户端池单例
var SharedHTTPClientPool = NewHTTPClientPool()

// 客户端池
type HTTPClientPool struct {
	clientsMap map[string]*http.Client // backend key => client
	locker     sync.RWMutex
}

// 获取新对象
func NewHTTPClientPool() *HTTPClientPool {
	return &HTTPClientPool{
		clientsMap: map[string]*http.Client{},
	}
}

// 根据地址获取客户端
func (this *HTTPClientPool) client(req *Request, backend *teaconfigs.BackendConfig, location *teaconfigs.LocationConfig) *http.Client {
	key := backend.UniqueKey()
	if location != nil {
		key = location.Id + "_" + key
	}

	backendAddr := backend.Address
	if backend.HasAddrVariables() {
		backendAddr = req.Format(backend.Address)
		key += "@" + backendAddr
	}

	this.locker.RLock()
	client, found := this.clientsMap[key]
	if found {
		this.locker.RUnlock()
		return client
	}
	this.locker.RUnlock()
	this.locker.Lock()

	maxConnections := int(backend.MaxConns)
	connectionTimeout := backend.FailTimeoutDuration()
	readTimeout := backend.ReadTimeoutDuration()
	idleTimeout := backend.IdleTimeoutDuration()
	idleConns := int(backend.IdleConns)

	// 超时时间
	if connectionTimeout <= 0 {
		connectionTimeout = 15 * time.Second
	}

	if idleTimeout <= 0 {
		idleTimeout = 2 * time.Minute
	}

	numberCPU := runtime.NumCPU()
	if numberCPU < 8 {
		numberCPU = 8
	}
	if maxConnections <= 0 {
		maxConnections = numberCPU
	}

	if idleConns <= 0 {
		idleConns = numberCPU
	}

	logs.Println("[proxy]setup backend '" + key + "', max connections:" + strconv.Itoa(maxConnections) + ", max idles:" + strconv.Itoa(idleConns))

	// TLS通讯
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	if backend.Cert != nil {
		obj := backend.Cert.CertObject()
		if obj != nil {
			tlsConfig.InsecureSkipVerify = false
			tlsConfig.Certificates = []tls.Certificate{*obj}
			if len(backend.Cert.ServerName) > 0 {
				tlsConfig.ServerName = backend.Cert.ServerName
			}
		}
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// 握手配置
			return (&net.Dialer{
				Timeout:   connectionTimeout,
				KeepAlive: 2 * time.Minute,
			}).DialContext(ctx, network, backendAddr)
		},
		MaxIdleConns:          0,
		MaxIdleConnsPerHost:   idleConns,
		MaxConnsPerHost:       maxConnections,
		IdleConnTimeout:       idleTimeout,
		ExpectContinueTimeout: 1 * time.Second,
		TLSHandshakeTimeout:   0, // 不限
		TLSClientConfig:       tlsConfig,
		Proxy:                 nil,
	}

	client = &http.Client{
		Timeout:   readTimeout,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	this.clientsMap[key] = client

	// 关闭老的
	if !backend.HasAddrVariables() {
		this.closeOldClient(key)
	}

	this.locker.Unlock()

	return client
}

// 关闭老的client
func (this *HTTPClientPool) closeOldClient(key string) {
	backendId := strings.Split(key, "@")[0]
	for key2, client := range this.clientsMap {
		backendId2 := strings.Split(key2, "@")[0]
		if backendId == backendId2 && key != key2 {
			teautils.CloseHTTPClient(client)
			delete(this.clientsMap, key2)
			logs.Println("[proxy]close backend '" + key2)
			break
		}
	}
}
