package ssl

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type AcmeDeleteTaskAction actions.Action

// 删除ACME任务
func (this *AcmeDeleteTaskAction) RunPost(params struct {
	ServerId string
	TaskId   string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	if server.SSL == nil {
		this.Success()
	}

	server.SSL.RemoveCertTask(params.TaskId)
	err := server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
