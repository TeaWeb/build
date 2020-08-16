package proxy

import (
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/iwind/TeaGo/actions"
)

type ClientDisconnectAction actions.Action

// 断开客户端连接
func (this *ClientDisconnectAction) RunPost(params struct {
	ServerId string
	Addr     string // client addr
}) {
	for _, listener := range teaproxy.SharedManager.FindServerListeners(params.ServerId) {
		listener.CloseTCPClient(params.Addr)
	}
	this.Success()
}
