package proxy

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"net"
	"strings"
)

type LocalAddrsAction actions.Action

// 本地的网卡绑定的IP地址
func (this *LocalAddrsAction) RunPost(params struct{}) {
	// 本机可用的IP
	result := []maps.Map{}
	interfaces, _ := net.Interfaces()
	for _, i := range interfaces {
		addrList, _ := i.Addrs()
		for _, addr := range addrList {
			if addr.Network() != "ip+net" {
				continue
			}
			pieces := strings.Split(addr.String(), "/")
			result = append(result, maps.Map{
				"name": i.Name,
				"addr": pieces[0],
			})
		}
	}

	this.Data["result"] = result

	this.Success()
}
