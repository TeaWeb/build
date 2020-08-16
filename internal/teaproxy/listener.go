package teaproxy

import (
	"crypto/tls"
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaplugins"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"golang.org/x/net/http2"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// 协议
type Scheme = uint8

const (
	SchemeHTTP   Scheme = 1
	SchemeHTTPS  Scheme = 2
	SchemeTCP    Scheme = 3
	SchemeTCPTLS Scheme = 4
)

// 对象pool
var requestPool = teautils.NewObjectPool(20480, func() interface{} {
	return &Request{
		isNew: true,
	}
})

// 代理服务监听器
type Listener struct {
	httpServer *http.Server // HTTP SERVER

	IsChanged bool // 标记是否改变，用来在其他地方重启改变的监听器

	Scheme  Scheme // http, https, tcp, tcp+tls
	Address string
	Error   error

	servers        []*teaconfigs.ServerConfig // 待启用的server
	currentServers []*teaconfigs.ServerConfig // 当前可用的Server
	namedServers   map[string]*NamedServer    // 域名 => server

	serversLocker      sync.RWMutex
	namedServersLocker sync.RWMutex

	// TCP
	tcpServer           net.Listener
	connectingTCPMap    map[net.Conn]*TCPClient
	connectingTCPLocker sync.Mutex
}

// 获取新对象
func NewListener() *Listener {
	return &Listener{
		namedServers:     map[string]*NamedServer{},
		connectingTCPMap: map[net.Conn]*TCPClient{},
	}
}

// 应用配置
func (this *Listener) ApplyServer(server *teaconfigs.ServerConfig) {
	this.serversLocker.Lock()
	defer this.serversLocker.Unlock()

	this.IsChanged = true

	isAvailable := false
	if this.Scheme == SchemeHTTP && server.Http {
		isAvailable = true
	} else if this.Scheme == SchemeHTTPS && server.SSL != nil && server.SSL.On {
		isAvailable = true
	} else if this.Scheme == SchemeTCP && server.TCP != nil && server.TCP.TCPOn {
		isAvailable = true
	} else if this.Scheme == SchemeTCPTLS && server.TCP != nil && server.SSL != nil && server.SSL.On {
		isAvailable = true
	}

	if !isAvailable {
		// 删除
		result := []*teaconfigs.ServerConfig{}
		for _, s := range this.servers {
			if s.Id == server.Id {
				continue
			}
			result = append(result, s)
		}
		this.servers = result

		return
	}

	found := false
	for index, s := range this.servers {
		if s.Id == server.Id {
			this.servers[index] = server
			found = true
			break
		}
	}
	if !found {
		this.servers = append(this.servers, server)
	}
}

// 删除配置
func (this *Listener) RemoveServer(serverId string) {
	this.serversLocker.Lock()
	defer this.serversLocker.Unlock()

	this.IsChanged = true
	result := []*teaconfigs.ServerConfig{}
	for _, s := range this.servers {
		if s.Id == serverId {
			continue
		}
		result = append(result, s)
	}
	this.servers = result
}

// 重置所有配置
func (this *Listener) Reset() {
	this.serversLocker.Lock()
	defer this.serversLocker.Unlock()

	this.IsChanged = true
	this.servers = []*teaconfigs.ServerConfig{}
}

// 判断是否包含某个配置
func (this *Listener) HasServer(serverId string) bool {
	this.serversLocker.RLock()
	defer this.serversLocker.RUnlock()

	for _, s := range this.servers {
		if s.Id == serverId {
			return true
		}
	}
	return false
}

// 是否包含配置
func (this *Listener) HasServers() bool {
	this.serversLocker.RLock()
	defer this.serversLocker.RUnlock()

	return len(this.servers) > 0
}

// 启动
func (this *Listener) Start() error {
	return this.Reload()
}

// 刷新
func (this *Listener) Reload() error {
	this.namedServersLocker.Lock()
	this.namedServers = map[string]*NamedServer{}
	this.namedServersLocker.Unlock()

	this.serversLocker.Lock()
	this.currentServers = this.servers
	hasServers := len(this.currentServers) > 0
	this.IsChanged = false
	this.Error = nil

	if !hasServers {
		defer this.serversLocker.Unlock()

		// 检查是否已启动
		return this.Shutdown()
	} else {
		this.serversLocker.Unlock()
	}

	var err error

	if this.Scheme == SchemeHTTP || this.Scheme == SchemeHTTPS { // HTTP
		err = this.startHTTPServer()
	} else if this.Scheme == SchemeTCP || this.Scheme == SchemeTCPTLS { // TCP
		err = this.startTCPServer()
	}

	return err
}

// 关闭
func (this *Listener) Shutdown() error {
	if this.Scheme == SchemeHTTP || this.Scheme == SchemeHTTPS { // HTTP
		if this.httpServer != nil {
			logs.Println("[proxy]shutdown listener on", this.Address)
			err := this.httpServer.Close()
			this.httpServer = nil
			return err
		}
	} else if this.Scheme == SchemeTCP || this.Scheme == SchemeTCPTLS {
		if this.tcpServer != nil {
			logs.Println("[proxy]shutdown listener on", this.Address)

			// 关闭listener
			err := this.tcpServer.Close()

			// 关闭现有连接
			this.connectingTCPLocker.Lock()
			for _, client := range this.connectingTCPMap {
				err1 := client.Close()
				if err1 != nil {
					logs.Error(err1)
				}
			}
			this.connectingTCPMap = map[net.Conn]*TCPClient{}
			this.connectingTCPLocker.Unlock()

			this.tcpServer = nil
			return err
		}
	}
	return nil
}

// 获取TCP连接列表
func (this *Listener) TCPClients(maxSize int) []*TCPClient {
	result := []*TCPClient{}

	if maxSize == 0 {
		return result
	}

	this.connectingTCPLocker.Lock()
	index := 0
	for _, client := range this.connectingTCPMap {
		if client.LConn() == nil || client.RConn() == nil {
			continue
		}

		index++
		result = append(result, client)
		if maxSize > 0 && index == maxSize-1 {
			break
		}
	}
	this.connectingTCPLocker.Unlock()
	lists.Sort(result, func(i int, j int) bool {
		c1 := result[i].LConn().RemoteAddr().String()
		c2 := result[j].LConn().RemoteAddr().String()
		return c1 < c2
	})
	return result
}

// 关闭某个TCP连接
func (this *Listener) CloseTCPClient(lAddr string) error {
	var err error
	this.connectingTCPLocker.Lock()
	for _, client := range this.connectingTCPMap {
		if client.LConn().RemoteAddr().String() == lAddr {
			err = client.Close()
			break
		}
	}
	this.connectingTCPLocker.Unlock()
	return err
}

// 启动HTTP Server
func (this *Listener) startHTTPServer() error {
	// 如果已经启动，则不做任何事情
	if this.httpServer != nil {
		return nil
	}

	defer func() {
		this.httpServer = nil
	}()

	var err error

	// 如果没启动，则启动
	httpHandler := new(HTTPServeMux)
	httpHandler.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		// QPS计算
		atomic.AddInt32(&qps, 1)

		// 处理
		this.handleHTTP(writer, req)
	})

	this.httpServer = &http.Server{
		Addr:        this.Address,
		Handler:     httpHandler,
		IdleTimeout: 2 * time.Minute,
	}
	this.httpServer.SetKeepAlivesEnabled(true)

	if this.Scheme == SchemeHTTP {
		logs.Println("[proxy]start listener on http", this.Address)
		err = this.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logs.Error(errors.New("[proxy]" + this.Address + ": " + err.Error()))
		} else {
			err = nil
		}
	}

	if this.Scheme == SchemeHTTPS {
		logs.Println("[proxy]start listener on https", this.Address)

		this.httpServer.TLSConfig = this.buildTLSConfig()

		// support http/2
		err = http2.ConfigureServer(this.httpServer, nil)
		if err != nil {
			logs.Error(err)
		}

		err = this.httpServer.ListenAndServeTLS("", "")
		if err != nil && err != http.ErrServerClosed {
			logs.Error(errors.New("[proxy]" + this.Address + ": " + err.Error()))
		} else {
			err = nil
		}
	}

	return err
}

// 启动TCP Server
func (this *Listener) startTCPServer() error {
	if this.tcpServer != nil {
		return nil
	}

	defer func() {
		this.tcpServer = nil
	}()

	var err error

	if this.Scheme == SchemeTCP {
		logs.Println("[proxy]start listener on tcp", this.Address)
		listener, err := net.Listen("tcp", this.Address)
		if err != nil {
			return errors.New("[proxy]tcp " + this.Address + ": " + err.Error())
		}
		this.tcpServer = listener
	} else if this.Scheme == SchemeTCPTLS {
		logs.Println("[proxy]start listener on tcp+tls", this.Address)
		listener, err := tls.Listen("tcp", this.Address, this.buildTLSConfig())
		if err != nil {
			return errors.New("[proxy]tcp " + this.Address + ": " + err.Error())
		}
		this.tcpServer = listener
	}

	// Accept
	server := this.tcpServer
	if server != nil {
		for {
			clientConn, err := server.Accept()
			if err != nil {
				break
			}

			var serverName = ""
			tlsConn, ok := clientConn.(*tls.Conn)
			if ok {
				go func() {
					err = tlsConn.Handshake()
					if err != nil {
						logs.Error(err)
						return
					}

					serverName = tlsConn.ConnectionState().ServerName
					this.connectTCPBackend(clientConn, serverName)
				}()
			} else {
				go this.connectTCPBackend(clientConn, serverName)
			}
		}
	}

	return err
}

// 处理请求
func (this *Listener) handleHTTP(writer http.ResponseWriter, rawRequest *http.Request) {
	// 插件过滤
	if teaplugins.HasRequestFilters {
		result, willContinue := teaplugins.FilterRequest(rawRequest)
		if !willContinue {
			return
		}
		rawRequest = result
	}

	// 域名
	reqHost := rawRequest.Host

	// TLS域名
	if this.isIP(reqHost) {
		if rawRequest.TLS != nil {
			serverName := rawRequest.TLS.ServerName
			if len(serverName) > 0 {
				// 端口
				index := strings.LastIndex(reqHost, ":")
				if index >= 0 {
					reqHost = serverName + reqHost[index:]
				} else {
					reqHost = serverName
				}
			}
		}
	}

	// 防止空Host
	if len(reqHost) == 0 {
		ctx := rawRequest.Context()
		if ctx != nil {
			addr := ctx.Value(http.LocalAddrContextKey)
			if addr != nil {
				reqHost = addr.(net.Addr).String()
			}
		}
	}

	domain, _, err := net.SplitHostPort(reqHost)
	if err != nil {
		domain = reqHost
	}
	server, serverName := this.findNamedServer(domain)
	if server == nil {
		// 严格匹配域名模式下，我们拒绝用户访问
		if teaconfigs.SharedProxySetting().MatchDomainStrictly {
			hijacker, ok := writer.(http.Hijacker)
			if ok {
				conn, _, _ := hijacker.Hijack()
				if conn != nil {
					_ = conn.Close()
					return
				}
			}
		}

		http.Error(writer, "404 page not found: '"+rawRequest.URL.String()+"'", http.StatusNotFound)
		return
	}

	// 包装新的请求
	req := requestPool.Get().(*Request)
	if req.isNew {
		req.isNew = false
		req.init(rawRequest)
		req.responseWriter = NewResponseWriter(writer)
	} else {
		req.reset(rawRequest)
		req.responseWriter.Reset(writer)
	}

	req.host = reqHost
	req.method = rawRequest.Method
	req.uri = rawRequest.URL.RequestURI()
	if this.Scheme == SchemeHTTP {
		req.rawScheme = "http"
	} else if this.Scheme == SchemeHTTPS {
		req.rawScheme = "https"
	} else {
		req.rawScheme = "http"
	}
	req.scheme = "http" // 转发后的scheme
	req.serverName = serverName
	req.serverAddr = this.Address
	req.root = server.Root
	req.index = server.Index
	req.charset = server.Charset

	// 配置请求
	err = req.configure(server, 0, false)
	if err != nil {
		req.serverError(req.responseWriter)
		logs.Error(errors.New(reqHost + rawRequest.URL.String() + ": " + err.Error()))

		// 返还request
		requestPool.Put(req)
		return
	}

	// 正向代理
	if server.ForwardHTTP != nil {
		err = req.Forward(req.responseWriter)
		if err != nil {
			logs.Error(errors.New(reqHost + rawRequest.URL.String() + ": " + err.Error()))
		}

		// 返还request
		requestPool.Put(req)

		return
	}

	// 处理请求
	err = req.call(req.responseWriter)
	if err != nil {
		// 已经在call()方法里处理过了这里不再重复
	}

	// 返还request
	requestPool.Put(req)
}

// 根据域名来查找匹配的域名
func (this *Listener) findNamedServer(name string) (serverConfig *teaconfigs.ServerConfig, serverName string) {
	// 读取缓存
	this.namedServersLocker.RLock()
	namedServer, found := this.namedServers[name]
	if found {
		this.namedServersLocker.RUnlock()
		return namedServer.Server, namedServer.Name
	}
	this.namedServersLocker.RUnlock()

	this.serversLocker.RLock()
	defer this.serversLocker.RUnlock()

	countServers := len(this.currentServers)
	if countServers == 0 {
		return nil, ""
	}

	// 只记录N个记录，防止内存耗尽
	maxNamedServers := 10240

	// 是否严格匹配域名
	matchDomainStrictly := teaconfigs.SharedProxySetting().MatchDomainStrictly

	// 如果只有一个server，则默认为这个
	if countServers == 1 && !matchDomainStrictly {
		server := this.currentServers[0]
		matchedName, matched := server.MatchName(name)
		if matched {
			if len(matchedName) > 0 {
				this.namedServersLocker.Lock()
				if len(this.namedServers) < maxNamedServers {
					this.namedServers[name] = &NamedServer{
						Name:   matchedName,
						Server: server,
					}
				}
				this.namedServersLocker.Unlock()
				return server, matchedName
			} else {
				return server, name
			}
		}

		// 匹配第一个域名
		firstName := server.FirstName()
		if len(firstName) > 0 {
			return server, firstName
		}
		return server, name
	}

	// 精确查找
	for _, server := range this.currentServers {
		if lists.ContainsString(server.Name, name) {
			this.namedServersLocker.Lock()
			if len(this.namedServers) < maxNamedServers {
				this.namedServers[name] = &NamedServer{
					Name:   name,
					Server: server,
				}
			}
			this.namedServersLocker.Unlock()
			return server, name
		}
	}

	// 模糊查找
	for _, server := range this.currentServers {
		if _, matched := server.MatchName(name); matched {
			this.namedServersLocker.Lock()
			if len(this.namedServers) < maxNamedServers {
				this.namedServers[name] = &NamedServer{
					Name:   name,
					Server: server,
				}
			}
			this.namedServersLocker.Unlock()
			return server, name
		}
	}

	// 如果没有找到，则匹配到第一个
	if matchDomainStrictly {
		return nil, name
	}

	server := this.currentServers[0]
	firstName := server.FirstName()
	if len(firstName) > 0 {
		this.namedServersLocker.Lock()
		if len(this.namedServers) < maxNamedServers {
			this.namedServers[name] = &NamedServer{
				Name:   firstName,
				Server: server,
			}
		}
		this.namedServersLocker.Unlock()
		return server, firstName
	}

	return server, name
}

// 根据域名匹配证书
func (this *Listener) matchSSL(domain string) (*teaconfigs.SSLConfig, *tls.Certificate, error) {
	this.serversLocker.RLock()
	defer this.serversLocker.RUnlock()

	// 如果域名为空，则取第一个
	// 通常域名为空是因为是直接通过IP访问的
	if len(domain) == 0 {
		if teaconfigs.SharedProxySetting().MatchDomainStrictly {
			return nil, nil, errors.New("[proxy]no tls server name matched")
		}

		if len(this.currentServers) > 0 && this.currentServers[0].SSL != nil {
			return this.currentServers[0].SSL, this.currentServers[0].SSL.FirstCert(), nil
		}
		return nil, nil, errors.New("[proxy]no tls server name found")
	}

	// 通过代理服务域名配置匹配
	server, _ := this.findNamedServer(domain)
	if server == nil || server.SSL == nil || !server.SSL.On {
		// 搜索所有的Server，通过SSL证书内容中的DNSName匹配
		for _, server := range this.currentServers {
			if server.SSL == nil || !server.SSL.On {
				continue
			}
			cert, ok := server.SSL.MatchDomain(domain)
			if ok {
				return server.SSL, cert, nil
			}
		}

		return nil, nil, errors.New("[proxy]no server found for '" + domain + "'")
	}

	// 证书是否匹配
	cert, ok := server.SSL.MatchDomain(domain)
	if ok {
		return server.SSL, cert, nil
	}

	return server.SSL, server.SSL.FirstCert(), nil
}

// 构造TLS配置
func (this *Listener) buildTLSConfig() *tls.Config {
	return &tls.Config{
		Certificates: nil,
		GetConfigForClient: func(info *tls.ClientHelloInfo) (config *tls.Config, e error) {
			ssl, _, err := this.matchSSL(info.ServerName)
			if err != nil {
				return nil, err
			}

			cipherSuites := ssl.TLSCipherSuites()
			if len(cipherSuites) == 0 {
				cipherSuites = nil
			}

			nextProto := []string{}
			if !ssl.HTTP2Disabled {
				nextProto = []string{http2.NextProtoTLS}
			}
			return &tls.Config{
				Certificates: nil,
				MinVersion:   ssl.TLSMinVersion(),
				CipherSuites: cipherSuites,
				GetCertificate: func(info *tls.ClientHelloInfo) (certificate *tls.Certificate, e error) {
					_, cert, err := this.matchSSL(info.ServerName)
					if err != nil {
						return nil, err
					}
					if cert == nil {
						return nil, errors.New("[proxy]no certs found for '" + info.ServerName + "'")
					}
					return cert, nil
				},
				ClientAuth: teaconfigs.GoSSLClientAuthType(ssl.ClientAuthType),
				ClientCAs:  ssl.CAPool(),

				NextProtos: nextProto,
			}, nil
		},
		GetCertificate: func(info *tls.ClientHelloInfo) (certificate *tls.Certificate, e error) {
			_, cert, err := this.matchSSL(info.ServerName)
			if err != nil {
				return nil, err
			}
			if cert == nil {
				return nil, errors.New("[proxy]no certs found for '" + info.ServerName + "'")
			}
			return cert, nil
		},
	}
}

// 连接TCP后端
func (this *Listener) connectTCPBackend(clientConn net.Conn, serverName string) {
	defer teautils.Recover()

	client := NewTCPClient(func() *teaconfigs.ServerConfig {
		this.serversLocker.RLock()

		if len(this.currentServers) == 0 {
			this.serversLocker.RUnlock()
			return nil
		}

		if len(serverName) == 0 {
			defer this.serversLocker.RUnlock()
			return this.currentServers[0]
		}

		this.serversLocker.RUnlock()

		server, _ := this.findNamedServer(serverName)
		return server
	}, clientConn)
	this.connectingTCPLocker.Lock()
	this.connectingTCPMap[clientConn] = client
	this.connectingTCPLocker.Unlock()

	client.Connect()

	this.connectingTCPLocker.Lock()
	delete(this.connectingTCPMap, clientConn)
	this.connectingTCPLocker.Unlock()
}

func (this *Listener) isIP(host string) bool {
	// IPv6
	if strings.Index(host, "[") > -1 {
		return true
	}

	for _, b := range host {
		if b >= 'a' && b <= 'z' {
			return false
		}
	}

	return true
}
