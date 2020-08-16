package notices

import (
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/iwind/TeaGo/actions"
)

type SetReadAction actions.Action

// 设置已读
func (this *SetReadAction) Run(params struct {
	AgentId   string
	Scope     string
	NoticeIds []string
}) {
	if params.Scope == "page" {
		if len(params.NoticeIds) == 0 {
			this.Success()
		}

		err := teadb.NoticeDAO().UpdateAgentNoticesRead(params.AgentId, params.NoticeIds)
		if err != nil {
			this.Fail("操作失败：" + err.Error())
		}
	} else {
		err := teadb.NoticeDAO().UpdateAllAgentNoticesRead(params.AgentId)
		if err != nil {
			this.Fail("操作失败：" + err.Error())
		}
	}

	this.Success()
}
