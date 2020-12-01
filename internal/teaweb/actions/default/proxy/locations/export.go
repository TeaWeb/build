package locations

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/rands"
	"gopkg.in/yaml.v3"
	"regexp"
)

type ExportAction actions.Action

// 导出路径规则
func (this *ExportAction) RunGet(params struct {
	ServerId   string
	LocationId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	location := server.FindLocation(params.LocationId)
	if location == nil {
		this.Fail("找不到Location")
	}

	location.Id = rands.HexString(16)

	data, err := yaml.Marshal(location)
	if err != nil {
		this.Fail(err.Error())
	}

	pattern := regexp.MustCompile(`[^\w]`).ReplaceAllLiteralString(location.PatternString(), "_")
	pattern = regexp.MustCompile(`_+`).ReplaceAllLiteralString(pattern, "_")

	this.AddHeader("Content-Disposition", "attachment; filename=location."+pattern+"_"+location.Id+".conf")
	this.Write(data)
}
