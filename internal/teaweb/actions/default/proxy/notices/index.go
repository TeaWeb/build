package notices

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

func (this *IndexAction) RunGet(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	this.Data["selectedTab"] = "notices"
	this.Data["server"] = server

	this.Data["levels"] = lists.Map(notices.AllNoticeLevels(), func(k int, v interface{}) interface{} {
		level := v.(maps.Map)
		code := level["code"].(notices.NoticeLevel)
		receivers, found := server.NoticeSetting[code]

		// 当前Agent的设置
		if found && len(receivers) > 0 {
			level["receivers"] = proxyutils.ConvertReceiversToMaps(receivers)
		} else {
			level["receivers"] = []interface{}{}
		}

		return level
	})

	this.Show()
}
