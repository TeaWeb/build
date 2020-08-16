package ssl

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type AcmeDownloadAction actions.Action

// 下载ACME证书
func (this *AcmeDownloadAction) RunGet(params struct {
	ServerId string
	TaskId   string
	Type     string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	this.Data["server"] = server
	this.Data["selectedTab"] = "https"

	var task *teaconfigs.SSLCertTask = nil
	if server.SSL != nil {
		task = server.SSL.FindCertTask(params.TaskId)
	}

	if task == nil {
		this.Fail("找不到Task")
	}

	if params.Type == "cert" {
		this.AddHeader("Content-Disposition", "attachment; filename=\"acme-"+task.Id+".pem\";")
		this.WriteString(task.Request.Cert)
	} else if params.Type == "key" {
		this.AddHeader("Content-Disposition", "attachment; filename=\"acme-"+task.Id+".key\";")
		this.WriteString(task.Request.Key)
	} else if params.Type == "viewCert" {
		this.WriteString(task.Request.Cert)
	} else if params.Type == "viewKey" {
		this.WriteString(task.Request.Key)
	} else {
		this.WriteString("unknown cert type")
	}
}
