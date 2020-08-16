package certs

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/utils/time"
	"time"
)

type AcmeTaskAction actions.Action

// ACME任务详情
func (this *AcmeTaskAction) RunGet(params struct {
	TaskId string
}) {
	list := teaconfigs.SharedSSLCertList()
	task := list.FindTask(params.TaskId)

	if task == nil {
		this.Fail("找不到Task")
	}

	this.Data["task"] = task

	date := task.Request.CertDate()
	this.Data["dayFrom"] = date[0]
	this.Data["dayTo"] = date[1]
	this.Data["isExpired"] = len(date[1]) > 0 && timeutil.Format("Y-m-d") > date[1]
	this.Data["runTime"] = timeutil.Format("Y-m-d H:i:s", time.Unix(task.RunAt, 0))
	this.Data["runError"] = task.RunError

	this.Show()
}
