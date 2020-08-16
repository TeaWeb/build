package notices

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/iwind/TeaGo/actions"
)

type DeleteNoticeReceiverAction actions.Action

// 删除接收人
func (this *DeleteNoticeReceiverAction) Run(params struct {
	ServerId   string
	Level      notices.NoticeLevel
	ReceiverId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	server.RemoveNoticeReceiver(params.Level, params.ReceiverId)
	err := server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
