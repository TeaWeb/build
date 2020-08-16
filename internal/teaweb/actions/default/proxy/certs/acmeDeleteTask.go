package certs

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type AcmeDeleteTaskAction actions.Action

// 删除ACME任务
func (this *AcmeDeleteTaskAction) RunPost(params struct {
	TaskId string
}) {
	list := teaconfigs.SharedSSLCertList()
	list.RemoveTask(params.TaskId)
	err := list.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
