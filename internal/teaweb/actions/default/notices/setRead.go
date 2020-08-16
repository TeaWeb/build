package notices

import (
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
)

type SetReadAction actions.Action

// 设置已读
func (this *SetReadAction) Run(params struct {
	Scope     string
	NoticeIds []string
}) {
	if params.Scope == "page" {
		if len(params.NoticeIds) == 0 {
			this.Success()
		}

		err := teadb.NoticeDAO().UpdateNoticesRead(params.NoticeIds)
		if err != nil {
			this.Fail("操作失败：" + err.Error())
		}
	} else {
		err := teadb.NoticeDAO().UpdateAllNoticesRead()
		if err != nil {
			logs.Error(err)
		}
	}

	this.Success()
}
