package certs

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type AcmeUsersAction actions.Action

// ACME用户列表
func (this *AcmeUsersAction) RunGet(params struct {}) {
	userList := teaconfigs.SharedACMELocalUserList()
	if len(userList.Users) > 0 {
		this.Data["users"] = userList.Users
	} else {
		this.Data["users"] = []*teaconfigs.ACMELocalUser{}
	}

	this.Show()
}
