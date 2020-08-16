package notices

import (
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
)

type BadgeAction actions.Action

// 计算未读数量
func (this *BadgeAction) Run(params struct{}) {
	if teadb.SharedDB().IsAvailable() {
		count, err := teadb.NoticeDAO().CountAllUnreadNotices()
		if err != nil {
			logs.Error(err)
		}
		this.Data["count"] = count
	} else {
		this.Data["count"] = 0
	}
	this.Data["soundOn"] = notices.SharedNoticeSetting().SoundOn

	this.Success()
}
