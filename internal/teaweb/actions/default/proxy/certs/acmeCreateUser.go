package certs

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type AcmeCreateUserAction actions.Action

// 创建用户
func (this *AcmeCreateUserAction) RunGet(params struct{}) {
	this.Show()
}

// 提交注册
func (this *AcmeCreateUserAction) RunPost(params struct {
	Email string
	Name  string
	Must  *actions.Must
}) {
	params.Must.
		Field("email", params.Email).
		Email("请输入正确的用户邮箱")

	user := teaconfigs.NewACMELocalUser()
	user.Email = params.Email
	user.Name = params.Name

	req := teaconfigs.NewACMERequest()
	req.User = user
	_, err := req.Client()
	if err != nil {
		this.Fail("注册失败：" + err.Error())
	}

	userList := teaconfigs.SharedACMELocalUserList()
	userList.AddUser(user)
	err = userList.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
