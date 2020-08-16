package ssl

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type AcmeUsersAction actions.Action

// ACME用户列表
func (this *AcmeUsersAction) RunGet(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}
	this.Data["selectedTab"] = "https"
	this.Data["server"] = server

	userList := teaconfigs.SharedACMELocalUserList()
	if len(userList.Users) > 0 {
		this.Data["users"] = userList.Users
	} else {
		this.Data["users"] = []*teaconfigs.ACMELocalUser{}
	}

	this.Show()
}
