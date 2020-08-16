package notices

import (
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/iwind/TeaGo/actions"
)

type ReceiverAction actions.Action

// 接收人信息
func (this *ReceiverAction) Run(params struct {
	ReceiverId string
}) {
	setting := notices.SharedNoticeSetting()
	level, receiver := setting.FindReceiver(params.ReceiverId)
	if receiver == nil {
		this.Fail("找不到Receiver")
	}

	media := setting.FindMedia(receiver.MediaId)

	this.Data["level"] = notices.FindNoticeLevel(level)
	this.Data["receiver"] = receiver
	this.Data["media"] = media
	this.Data["mediaType"] = notices.FindNoticeMediaType(media.Type)

	this.Show()
}
