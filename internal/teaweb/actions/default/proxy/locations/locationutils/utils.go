package locationutils

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

// 设置公用信息
func SetCommonInfo(action actions.ActionWrapper, serverId string, locationId string, subTab string) (server *teaconfigs.ServerConfig, location *teaconfigs.LocationConfig) {
	obj := action.Object()

	server = teaconfigs.NewServerConfigFromId(serverId)
	if server == nil {
		obj.Fail("找不到Server")
	}

	location = server.FindLocation(locationId)
	if location == nil {
		obj.Fail("找不到要修改的Location")
	}

	obj.Data["location"] = maps.Map{
		"id":          location.Id,
		"pattern":     location.PatternString(),
		"fastcgi":     location.Fastcgi,
		"rewrite":     location.Rewrite,
		"headers":     location.Headers,
		"cachePolicy": location.CachePolicy,
		"websocket":   location.Websocket,
		"backends":    location.Backends,
		"wafOn":       location.WAFOn,
		"wafId":       location.WafId,
	}

	obj.Data["selectedTab"] = "location"
	obj.Data["selectedSubTab"] = subTab
	obj.Data["server"] = server

	return server, location
}
