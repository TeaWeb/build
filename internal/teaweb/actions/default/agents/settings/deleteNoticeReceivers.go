package settings

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/iwind/TeaGo/actions"
)

type DeleteNoticeReceiversAction actions.Action

// 删除接收人
func (this *DeleteNoticeReceiversAction) Run(params struct {
	AgentId    string
	Level      notices.NoticeLevel
	ReceiverId string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	agent.RemoveNoticeReceiver(params.Level, params.ReceiverId)
	err := agent.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
