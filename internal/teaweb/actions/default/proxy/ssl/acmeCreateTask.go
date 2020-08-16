package ssl

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type AcmeCreateTaskAction actions.Action

// 创建证书生成任务
func (this *AcmeCreateTaskAction) RunGet(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}
	this.Data["selectedTab"] = "https"
	this.Data["server"] = server

	users := teaconfigs.SharedACMELocalUserList().Users
	if len(users) > 0 {
		this.Data["users"] = users
	} else {
		this.Data["users"] = []*teaconfigs.ACMELocalUser{}
	}

	this.Show()
}

func (this *AcmeCreateTaskAction) RunPost(params struct {
	Must *actions.Must
}) {

}
