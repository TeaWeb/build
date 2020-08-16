package websocket

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/locations/locationutils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateAction actions.Action

// 更改Websocket
func (this *UpdateAction) Run(params struct {
	From       string
	ServerId   string
	LocationId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	this.Data["server"] = server

	location := server.FindLocation(params.LocationId)
	if location == nil {
		this.Fail("找不到要修改的路径配置")
	}

	this.Data["selectedTab"] = "location"
	this.Data["selectedSubTab"] = "websocket"

	this.Data["location"] = maps.Map{
		"on":          location.On,
		"id":          location.Id,
		"pattern":     location.PatternString(),
		"fastcgi":     location.Fastcgi,
		"headers":     location.Headers,
		"cachePolicy": location.CachePolicy,
		"rewrite":     location.Rewrite,
		"backends":    location.Backends,
		"wafId":       location.WafId,
		"wafOn":       location.WAFOn,
	}
	this.Data["from"] = params.From

	hasWebsocket := false
	if location.Websocket == nil {
		location.Websocket = teaconfigs.NewWebsocketConfig()
		location.Websocket.ForwardMode = teaconfigs.WebsocketForwardModeWebsocket
	} else {
		hasWebsocket = true
	}
	location.Websocket.Validate()

	this.Data["websocket"] = location.Websocket
	this.Data["handshakeTimeout"] = int(location.Websocket.HandshakeTimeoutDuration().Seconds())
	if hasWebsocket {
		this.Data["allowAllOrigins"] = location.Websocket.AllowAllOrigins
		this.Data["origins"] = location.Websocket.Origins
	} else {
		this.Data["allowAllOrigins"] = true
		this.Data["origins"] = []string{}
	}
	this.Data["modes"] = teaconfigs.AllWebsocketForwardModes()

	this.Show()
}

// 提交修改
func (this *UpdateAction) RunPost(params struct {
	ServerId         string
	LocationId       string
	On               bool
	HandshakeTimeout uint
	AllowAllOrigins  bool
	Origins          []string
	ForwardMode      string
}) {
	server, location := locationutils.SetCommonInfo(this, params.ServerId, params.LocationId, "websocket")

	if location.Websocket == nil {
		location.Websocket = teaconfigs.NewWebsocketConfig()
	}
	location.Websocket.On = params.On
	location.Websocket.HandshakeTimeout = fmt.Sprintf("%ds", params.HandshakeTimeout)
	location.Websocket.AllowAllOrigins = params.AllowAllOrigins
	location.Websocket.Origins = params.Origins
	location.Websocket.ForwardMode = params.ForwardMode

	err := server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()

	this.Success()
}
