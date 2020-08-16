package notices

import (
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/iwind/TeaGo/actions"
)

type DeleteReceiverAction actions.Action

// 删除接收人
func (this *DeleteReceiverAction) Run(params struct {
	Level      notices.NoticeLevel
	ReceiverId string
}) {
	setting := notices.SharedNoticeSetting()
	level := setting.LevelConfig(params.Level)
	level.RemoveReceiver(params.ReceiverId)
	err := setting.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
