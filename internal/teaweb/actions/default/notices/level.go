package notices

import (
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

type LevelAction actions.Action

// 级别
func (this *LevelAction) Run(params struct {
	Level uint8
}) {
	level := notices.FindNoticeLevel(params.Level)
	if level == nil {
		this.Fail("找不到Level信息")
	}

	this.Data["level"] = level

	setting := notices.SharedNoticeSetting()
	config := setting.LevelConfig(types.Uint8(params.Level))

	receivers := []maps.Map{}
	for _, receiver := range config.Receivers {
		media := setting.FindMedia(receiver.MediaId)
		if media == nil {
			continue
		}

		receivers = append(receivers, maps.Map{
			"on":            receiver.On,
			"id":            receiver.Id,
			"name":          receiver.Name,
			"user":          receiver.User,
			"mediaName":     media.Name,
			"mediaType":     media.Type,
			"mediaTypeName": notices.FindNoticeMediaTypeName(media.Type),
		})
	}

	this.Data["receivers"] = receivers

	this.Show()
}
