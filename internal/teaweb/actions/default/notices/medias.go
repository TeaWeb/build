package notices

import (
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type MediasAction actions.Action

// 媒介列表
func (this *MediasAction) Run(params struct{}) {
	setting := notices.SharedNoticeSetting()
	this.Data["medias"] = lists.Map(setting.Medias, func(k int, v interface{}) interface{} {
		media := v.(*notices.NoticeMediaConfig)
		return maps.Map{
			"on":       media.On,
			"id":       media.Id,
			"name":     media.Name,
			"typeName": notices.FindNoticeMediaTypeName(media.Type),
		}
	})

	this.Show()
}
