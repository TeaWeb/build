package certs

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type AcmeCreateTaskAction actions.Action

// 创建证书生成任务
func (this *AcmeCreateTaskAction) RunGet(params struct{}) {
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
