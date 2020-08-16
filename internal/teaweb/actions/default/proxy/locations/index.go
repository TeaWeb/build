package locations

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

// 路径规则列表
func (this *IndexAction) Run(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	this.Data["selectedTab"] = "location"
	this.Data["server"] = server

	locations := []maps.Map{}
	for _, location := range server.Locations {
		err := location.Validate()
		if err != nil {
			logs.Error(err)
		}
		locations = append(locations, maps.Map{
			"on":                location.On,
			"id":                location.Id,
			"name":              location.Name,
			"type":              location.PatternType(),
			"pattern":           location.PatternString(),
			"patternTypeName":   teaconfigs.FindLocationPatternTypeName(location.PatternType()),
			"isBreak":           location.IsBreak,
			"isCaseInsensitive": location.IsCaseInsensitive(),
			"isReverse":         location.IsReverse(),
			"rewrite":           location.Rewrite,
			"headers":           location.Headers,
			"fastcgi":           location.Fastcgi,
			"root":              location.Root,
			"index":             location.Index,
			//"gzipLevel":         location.GzipLevel,
			"cachePolicy": location.CachePolicyObject(),
			"websocket":   location.Websocket != nil && location.Websocket.On,
			"backends":    location.Backends,
			"hasWAF":      len(location.WafId) > 0 && location.WAFOn,
		})
	}

	this.Data["locations"] = locations

	this.Show()
}
