package certs

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type AcmeUserDeleteAction actions.Action

// 删除用户
func (this *AcmeUserDeleteAction) RunPost(params struct {
	UserId string
}) {
	userList := teaconfigs.SharedACMELocalUserList()

	if userList.FindUser(params.UserId) == nil {
		this.Fail("找不到要删除的用户")
	}

	userList.RemoveUser(params.UserId)
	err := userList.Save()
	if err != nil {
		this.Fail("删除失败：" + err.Error())
	}

	this.Success()
}
