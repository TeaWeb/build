package server

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/settings"
	"github.com/iwind/TeaGo/actions"
	"net"
)

type HttpUpdateAction actions.Action

// 保存HTTP设置
func (this *HttpUpdateAction) Run(params struct {
	On           bool
	ListenValues []string
	Must         *actions.Must
}) {
	if len(params.ListenValues) == 0 {
		this.Fail("请输入绑定地址")
	}

	server, err := teaconfigs.LoadWebConfig()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	server.Http.On = params.On

	listen := []string{}
	for _, addr := range params.ListenValues {
		addr = teautils.FormatAddress(addr)
		if len(addr) == 0 {
			continue
		}
		if _, _, err := net.SplitHostPort(addr); err != nil {
			addr += ":80"
		}
		listen = append(listen, addr)
	}
	server.Http.Listen = listen

	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	settings.NotifyServerChange()

	this.Next("/settings", nil).
		Success("保存成功，重启服务后生效")
}
