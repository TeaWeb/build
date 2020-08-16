package teaproxy

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/jlaffaye/ftp"
	"runtime"
	"strings"
	"sync"
)

// FTP客户端池单例
var SharedFTPClientPool = NewFTPClientPool()

// FTP客户端Pool
type FTPClientPool struct {
	clientMap map[string]*FTPClient // key => client
	locker    sync.Mutex
}

// 获取新的客户端Pool
func NewFTPClientPool() *FTPClientPool {
	return &FTPClientPool{
		clientMap: map[string]*FTPClient{},
	}
}

// 通过Backend配置FTP客户端
func (this *FTPClientPool) client(req *Request, backend *teaconfigs.BackendConfig, location *teaconfigs.LocationConfig) *FTPClient {
	key := backend.UniqueKey()
	if location != nil {
		key = location.Id + "_" + key
	}

	backendAddr := backend.Address
	if backend.HasAddrVariables() {
		backendAddr = req.Format(backend.Address)
		key += "@" + backendAddr
	}

	this.locker.Lock()
	defer this.locker.Unlock()
	client, ok := this.clientMap[key]
	if ok {
		return client
	}

	// 关闭以前的连接
	if !backend.HasAddrVariables() {
		this.closeOldClients(key)
	}

	if backend.FTP == nil {
		backend.FTP = &teaconfigs.FTPBackendConfig{}
	}
	numberCPU := runtime.NumCPU()
	if numberCPU < 8 {
		numberCPU = 8
	}
	if backend.IdleConns <= 0 {
		backend.IdleConns = int32(numberCPU)
	}
	client = &FTPClient{
		pool: &FTPConnectionPool{
			addr:           backendAddr,
			username:       backend.FTP.Username,
			password:       backend.FTP.Password,
			dir:            backend.FTP.Dir,
			timeout:        backend.FailTimeoutDuration(),
			c:              make(chan *ftp.ServerConn, backend.IdleConns),
			maxConnections: int64(backend.MaxConns),
		},
	}
	this.clientMap[key] = client

	return client
}

// 关闭老的client
func (this *FTPClientPool) closeOldClients(key string) {
	backendId := strings.Split(key, "@")[0]
	for key2, client := range this.clientMap {
		backendId2 := strings.Split(key2, "@")[0]
		if backendId == backendId2 && key != key2 {
			go func() {
				_ = client.Close()
			}()
			delete(this.clientMap, key2)
			break
		}
	}
}
