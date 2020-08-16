package groups

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/iwind/TeaGo/actions"
)

type DeleteNoticeReceiversAction actions.Action

// 删除接收人
func (this *DeleteNoticeReceiversAction) Run(params struct {
	GroupId    string
	Level      notices.NoticeLevel
	ReceiverId string
}) {
	config := agents.SharedGroupList()
	group := config.FindGroup(params.GroupId)
	if group == nil {
		this.Fail("要删除的组不存在")
	}

	group.RemoveNoticeReceiver(params.Level, params.ReceiverId)
	err := config.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
