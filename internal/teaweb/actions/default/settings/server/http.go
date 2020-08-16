package server

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type HttpAction actions.Action

func (this *HttpAction) Run(params struct{}) {
	this.Data["error"] = ""

	server, err := teaconfigs.LoadWebConfig()
	if err != nil {
		this.Data["error"] = "读取配置错误：" + err.Error()
		this.Show()
		return
	}

	this.Data["server"] = server

	this.Show()
}
