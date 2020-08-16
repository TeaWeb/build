package server

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/certs/certutils"
	"github.com/iwind/TeaGo/actions"
)

type HttpsAction actions.Action

func (this *HttpsAction) Run(params struct{}) {
	this.Data["error"] = ""

	server, err := teaconfigs.LoadWebConfig()
	if err != nil {
		this.Data["error"] = "读取配置错误：" + err.Error()
		this.Show()
		return
	}

	// 公共可以使用的证书
	this.Data["sharedCerts"] = certutils.ListPairCertsMap()

	this.Data["server"] = server
	this.Show()
}
