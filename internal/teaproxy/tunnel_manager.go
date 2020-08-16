package teaproxy

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/logs"
	"sync"
)

var SharedTunnelManager = NewTunnelManager()

// Tunnel管理器
type TunnelManager struct {
	tunnelMap map[string][]*Tunnel // serverId => tunnels
	locker    sync.Mutex
}

// 获取新对象
func NewTunnelManager() *TunnelManager {
	return &TunnelManager{
		tunnelMap: map[string][]*Tunnel{},
	}
}

// 加载
func (this *TunnelManager) OnAttach(server *teaconfigs.ServerConfig) {
	this.locker.Lock()
	defer this.locker.Unlock()

	tunnels, ok := this.tunnelMap[server.Id]
	if ok {
		for _, tunnel := range tunnels {
			tunnel.Close()
		}
		delete(this.tunnelMap, server.Id)
	}

	if !server.On {
		return
	}

	if server.Tunnel != nil && server.Tunnel.On {
		tunnel := NewTunnel(server.Tunnel)
		this.tunnelMap[server.Id] = []*Tunnel{tunnel}
		go func() {
			server.Tunnel.SetIsActive(true)
			err := tunnel.Start()
			if err != nil {
				server.Tunnel.AddError(err.Error())
				server.Tunnel.SetIsActive(false)
				logs.Println("[tunnel]" + err.Error())
			}
		}()
	}
}

// 卸载
func (this *TunnelManager) OnDetach(server *teaconfigs.ServerConfig) {
	this.locker.Lock()
	defer this.locker.Unlock()

	tunnels, ok := this.tunnelMap[server.Id]
	if ok {
		for _, tunnel := range tunnels {
			tunnel.Close()
		}
		delete(this.tunnelMap, server.Id)
	}
}

// 查找Tunnel
func (this *TunnelManager) FindTunnel(serverId, tunnelId string) *Tunnel {
	this.locker.Lock()
	defer this.locker.Unlock()

	tunnels, ok := this.tunnelMap[serverId]
	if !ok {
		return nil
	}

	for _, tunnel := range tunnels {
		if tunnel.Id() == tunnelId {
			return tunnel
		}
	}

	return nil
}
