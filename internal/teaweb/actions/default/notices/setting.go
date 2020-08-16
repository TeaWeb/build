package notices

import (
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

type SettingAction actions.Action

// 接收设置
func (this *SettingAction) Run(params struct{}) {
	setting := notices.SharedNoticeSetting()
	this.Data["levels"] = lists.Map(notices.AllNoticeLevels(), func(k int, v interface{}) interface{} {
		m := v.(maps.Map)

		receivers := lists.Map(setting.LevelConfig(types.Uint8(m["code"])).Receivers, func(k int, v interface{}) interface{} {
			r := v.(*notices.NoticeReceiver)
			media := setting.FindMedia(r.MediaId)
			return maps.Map{
				"name":  r.Name,
				"media": media.Name,
			}
		})
		m["receivers"] = receivers

		return m
	})

	this.Show()
}
