package locations

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/rands"
	"gopkg.in/yaml.v3"
)

type DuplicateAction actions.Action

// 复制
func (this *DuplicateAction) RunPost(params struct {
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

	data, err := yaml.Marshal(location)
	if err != nil {
		this.Fail(err.Error())
	}

	newLocation := &teaconfigs.LocationConfig{}
	err = yaml.Unmarshal(data, newLocation)
	if err != nil {
		this.Fail(err.Error())
	}

	newLocation.Id = rands.HexString(16)
	if len(newLocation.Name) == 0 {
		newLocation.Name = "复制自" + location.PatternString()
	} else {
		newLocation.Name += "（复制自" + location.PatternString() + "）"
	}

	result := []*teaconfigs.LocationConfig{}
	for _, l := range server.Locations {
		result = append(result, l)
		if l.Id == location.Id {
			result = append(result, newLocation)
		}
	}
	server.Locations = result

	err = server.Save()
	if err != nil {
		this.Fail()
	}

	// 通知更新
	proxyutils.NotifyChange()

	this.Success()
}
