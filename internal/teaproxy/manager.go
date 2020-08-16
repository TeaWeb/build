package teaproxy

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"gopkg.in/yaml.v3"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"io/ioutil"
	"net"
	"path/filepath"
	"strings"
	"sync"
)

// 共享的管理对象
var SharedManager = NewManager()

// 管理器
type Manager struct {
	listeners  map[string]*Listener // scheme://address => listener
	oldServers map[string]*teaconfigs.ServerConfig
	servers    map[string]*teaconfigs.ServerConfig

	done   chan bool
	locker sync.RWMutex
}

// 获取新对象
func NewManager() *Manager {
	return &Manager{
		listeners: map[string]*Listener{},
		servers:   map[string]*teaconfigs.ServerConfig{},
		done:      make(chan bool),
	}
}

//  启动
func (this *Manager) Start() error {
	configsDir := Tea.ConfigDir()

	serverFiles := []string{}
	serverList, err := teaconfigs.SharedServerList()
	if err != nil {
		files, err := filepath.Glob(configsDir + Tea.DS + "*.proxy.conf")
		if err != nil {
			return err
		}
		serverFiles = files
	} else {
		for _, f := range serverList.Files {
			serverFiles = append(serverFiles, Tea.ConfigFile(f))
		}
	}

	this.servers = map[string]*teaconfigs.ServerConfig{}
	for _, configFile := range serverFiles {
		if strings.Contains(configFile, "server.sample.www.proxy.conf") { // 跳过示例配置
			continue
		}

		configData, err := ioutil.ReadFile(configFile)
		if err != nil {
			logs.Error(err)
			continue
		}

		server := &teaconfigs.ServerConfig{}
		err = yaml.Unmarshal(configData, server)
		if err != nil {
			logs.Error(err)
			continue
		}

		// 老的server
		oldServer, foundOldServer := this.oldServers[server.Id]
		if foundOldServer {
			server.CloneState(oldServer)
		}

		err = server.Validate()
		if err != nil {
			logs.Error(err)
			continue
		}

		this.servers[server.Id] = server
		this.ApplyServer(server)
	}

	err = this.Reload()
	if err != nil {
		return err
	}

	return nil
}

// 重启
func (this *Manager) Restart() error {
	this.locker.Lock()
	this.oldServers = this.servers
	this.servers = map[string]*teaconfigs.ServerConfig{}
	for _, listener := range this.listeners {
		listener.Reset()
	}
	this.locker.Unlock()
	return this.Start()
}

// 添加Server
func (this *Manager) ApplyServer(server *teaconfigs.ServerConfig) {
	this.locker.Lock()
	defer this.locker.Unlock()

	this.servers[server.Id] = server

	// old servers
	oldServer, ok := this.oldServers[server.Id]
	if ok {
		if server.Version > oldServer.Version {
			this.detachServer(oldServer)
			this.attachServer(server)
		}
	} else {
		this.attachServer(server)
	}

	keys := []string{}

	if server.IsHTTP() { // HTTP
		// HTTP
		if server.Http {
			for _, address := range server.ParseListenAddresses() {
				// 是否有端口
				if shared.RegexpDigitNumber.MatchString(address) {
					address = ":" + address
				} else if _, _, err := net.SplitHostPort(address); err != nil {
					address += ":80"
				}

				if len(address) > 0 {
					keys = append(keys, "http://"+address)
				}
			}
		}

		// HTTPS
		if server.SSL != nil && server.SSL.On {
			err := server.SSL.Validate()
			if err != nil {
				logs.Error(err)
			}
			for _, address := range server.SSL.ParseListenAddresses() {
				// 是否有端口
				if shared.RegexpDigitNumber.MatchString(address) {
					address = ":" + address
				} else if _, _, err := net.SplitHostPort(address); err != nil {
					address += ":443"
				}

				if len(address) > 0 {
					keys = append(keys, "https://"+address)
				}
			}
		}
	} else if server.IsTCP() { // TCP
		// TCP
		if server.TCP.TCPOn {
			for _, address := range server.ParseListenAddresses() {
				// 是否有端口
				if shared.RegexpDigitNumber.MatchString(address) {
					address = ":" + address
				} else if _, _, err := net.SplitHostPort(address); err != nil {
					logs.Println("[proxy]invalid tcp address: '" + address + "'")
					continue
				}

				if len(address) > 0 {
					keys = append(keys, "tcp://"+address)
				}
			}
		}

		// TCP+TLS
		if server.SSL != nil && server.SSL.On {
			err := server.SSL.Validate()
			if err != nil {
				logs.Error(err)
			}
			for _, address := range server.SSL.ParseListenAddresses() {
				// 是否有端口
				if shared.RegexpDigitNumber.MatchString(address) {
					address = ":" + address
				} else if _, _, err := net.SplitHostPort(address); err != nil {
					logs.Println("[proxy]invalid tcp+tls address: '" + address + "'")
					continue
				}

				if len(address) > 0 {
					keys = append(keys, "tcp+tls://"+address)
				}
			}
		}
	}

	// 删除
	for _, listener := range this.listeners {
		if listener.HasServer(server.Id) {
			listener.RemoveServer(server.Id)
		}
	}

	// 添加
	for _, key := range keys {
		pieces := strings.SplitN(key, "://", 2)
		scheme := pieces[0]
		address := pieces[1]

		if !server.On {
			continue
		}
		if scheme == "http" && !server.Http {
			continue
		}
		if scheme == "https" && (server.SSL == nil || !server.SSL.On) {
			continue
		}
		if scheme == "tcp" && (server.TCP == nil || !server.TCP.TCPOn) {
			continue
		}
		if scheme == "tcp+tls" && (server.TCP == nil || server.SSL == nil || !server.SSL.On) {
			continue
		}

		listener, found := this.listeners[key]
		if found {
			listener.ApplyServer(server)
		} else {
			listener := NewListener()
			if scheme == "http" {
				listener.Scheme = SchemeHTTP
			} else if scheme == "https" {
				listener.Scheme = SchemeHTTPS
			} else if scheme == "tcp" {
				listener.Scheme = SchemeTCP
			} else if scheme == "tcp+tls" {
				listener.Scheme = SchemeTCPTLS
			}
			listener.Address = address
			listener.ApplyServer(server)
			this.listeners[key] = listener
		}
	}
}

// 删除Server
func (this *Manager) RemoveServer(serverId string) {
	this.locker.Lock()
	defer this.locker.Unlock()

	server, ok := this.servers[serverId]
	if ok {
		this.detachServer(server)
		delete(this.servers, serverId)
	}

	for _, listener := range this.listeners {
		if listener.HasServer(serverId) {
			err := listener.Reload()
			if err != nil {
				logs.Error(err)
			}
		}
	}
}

// 查找Server
func (this *Manager) FindServer(serverId string) *teaconfigs.ServerConfig {
	this.locker.RLock()
	defer this.locker.RUnlock()

	server, found := this.servers[serverId]
	if found {
		return server
	}
	return nil
}

// 查找所有的Server
func (this *Manager) FindAllServers() []*teaconfigs.ServerConfig {
	this.locker.RLock()
	defer this.locker.RUnlock()

	result := []*teaconfigs.ServerConfig{}
	for _, server := range this.servers {
		result = append(result, server)
	}
	return result
}

// 查找Server相关错误
func (this *Manager) FindServerErrors(serverId string) []string {
	this.locker.RLock()
	defer this.locker.RUnlock()

	errs := []string{}
	for _, listener := range this.listeners {
		if !listener.HasServer(serverId) {
			continue
		}
		if listener.Error != nil {
			errs = append(errs, listener.Error.Error())
		}
	}

	return errs
}

//  查找Server对应的Listener
func (this *Manager) FindServerListeners(serverId string) []*Listener {
	this.locker.RLock()
	defer this.locker.RUnlock()

	result := []*Listener{}
	for _, listener := range this.listeners {
		if !listener.HasServer(serverId) {
			continue
		}
		result = append(result, listener)
	}

	return result
}

// 重载配置
func (this *Manager) Reload() error {
	this.locker.Lock()
	defer this.locker.Unlock()

	// 清理空的listener
	for key, listener := range this.listeners {
		if !listener.HasServers() {
			err := listener.Reload()
			if err != nil {
				return err
			}
			delete(this.listeners, key)
		}
	}

	// 启动新的有变化的listener
	for _, listener := range this.listeners {
		if listener.IsChanged {
			go func(listener *Listener) {
				err := listener.Reload()
				if err != nil {
					listener.Error = err
				}
			}(listener)
		}
	}

	return nil
}

// 等待
func (this *Manager) Wait() {
	<-this.done
}

// 停止
func (this *Manager) Shutdown() error {
	this.done <- true

	this.locker.Lock()
	defer this.locker.Unlock()

	for _, listener := range this.listeners {
		err := listener.Shutdown()
		if err != nil {
			return err
		}
	}

	return nil
}

// attach server
func (this *Manager) attachServer(server *teaconfigs.ServerConfig) {
	server.OnAttach()
	SharedTunnelManager.OnAttach(server)
}

// detach server
func (this *Manager) detachServer(server *teaconfigs.ServerConfig) {
	server.OnDetach()
	SharedTunnelManager.OnDetach(server)
}
