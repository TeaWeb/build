package websocket

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/locations/locationutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

// Websocket信息
func (this *IndexAction) Run(params struct {
	ServerId   string
	LocationId string
}) {
	this.Data["queryParams"] = maps.Map{
		"serverId":   params.ServerId,
		"locationId": params.LocationId,
		"websocket":  1,
	}

	server, location := locationutils.SetCommonInfo(this, params.ServerId, params.LocationId, "websocket")

	if location.Websocket == nil {
		this.Data["websocket"] = nil
	} else {
		this.Data["websocket"] = maps.Map{
			"on":               location.Websocket.On,
			"allowAllOrigins":  location.Websocket.AllowAllOrigins,
			"origins":          location.Websocket.Origins,
			"handshakeTimeout": location.Websocket.HandshakeTimeout,
			"forwardMode":      location.Websocket.ForwardModeSummary(),
		}
	}

	this.Data["isTCP"] = server.IsTCP()
	this.Data["isHTTP"] = server.IsHTTP()

	this.Show()
}
