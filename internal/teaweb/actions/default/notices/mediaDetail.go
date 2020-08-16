package notices

import (
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/iwind/TeaGo/actions"
)

type MediaDetailAction actions.Action

// 媒介详情
func (this *MediaDetailAction) Run(params struct {
	MediaId string
}) {
	setting := notices.SharedNoticeSetting()
	media := setting.FindMedia(params.MediaId)
	if media == nil {
		this.Fail("找不到Media")
	}

	this.Data["media"] = media

	// 媒介类型
	this.Data["mediaType"] = notices.FindNoticeMediaType(media.Type)

	this.Show()
}
