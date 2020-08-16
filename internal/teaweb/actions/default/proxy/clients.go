package proxy

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teageo"
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"net"
)

type ClientsAction actions.Action

// 客户端连接管理
func (this *ClientsAction) RunGet(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	if server.Index == nil {
		server.Index = []string{}
	}

	this.Data["selectedTab"] = "clients"
	this.Data["server"] = server

	if !server.IsTCP() {
		this.RedirectURL("/proxy/detail?serverId=" + params.ServerId)
		return
	}

	this.Show()
}

// 客户端连接数据
func (this *ClientsAction) RunPost(params struct {
	ServerId string
	Size     int

	Must *actions.Must
}) {
	if params.Size <= 0 {
		params.Size = 100
	}

	clients := []maps.Map{}
	for _, listener := range teaproxy.SharedManager.FindServerListeners(params.ServerId) {
		for _, pair := range listener.TCPClients(params.Size) {
			clientAddr := pair.LConn().RemoteAddr().String()
			ip, _, err := net.SplitHostPort(clientAddr)
			location := ""
			if err == nil {
				ipObj := net.ParseIP(ip)
				if ipObj != nil {
					record, err := teageo.DB.City(ipObj)
					if err == nil {
						if _, ok := record.Country.Names["zh-CN"]; ok {
							location = teageo.ConvertName(record.Country.Names["zh-CN"])
						}
						if _, ok := record.City.Names["zh-CN"]; ok {
							location += " " + teageo.ConvertName(record.City.Names["zh-CN"])
						}
					}
				}
			}

			clients = append(clients, maps.Map{
				"clientAddr":     clientAddr,
				"clientLocation": location,
				"serverAddr":     pair.LConn().LocalAddr().String(),
				"backendAddr":    pair.RConn().RemoteAddr().String(),
				"readSpeed":      pair.ReadSpeed(),
				"writeSpeed":     pair.WriteSpeed(),
			})
		}
	}

	this.Data["clients"] = clients
	this.Success()
}
