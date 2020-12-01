package server

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/settings"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/rands"
	"net"
	"strings"
)

type HttpsUpdateAction actions.Action

// 保存HTTPS设置
func (this *HttpsUpdateAction) Run(params struct {
	On           bool
	ListenValues []string

	CertType string
	CertFile *actions.File
	KeyFile  *actions.File
	CertId   string

	Must *actions.Must
}) {
	if len(params.ListenValues) == 0 {
		this.Fail("请输入绑定地址")
	}

	server, err := teaconfigs.LoadWebConfig()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	server.Https.On = params.On

	listen := []string{}
	for _, addr := range params.ListenValues {
		addr = teautils.FormatAddress(addr)
		if len(addr) == 0 {
			continue
		}
		if _, _, err := net.SplitHostPort(addr); err != nil {
			addr += ":443"
		}
		listen = append(listen, addr)
	}
	server.Https.Listen = listen

	if params.CertType == "shared" {
		if len(params.CertId) == 0 {
			this.Fail("请选择证书文件")
		}

		cert := teaconfigs.SharedSSLCertList().FindCert(params.CertId)
		if cert == nil {
			this.Fail("找不到选择的证书")
		}
		server.CertId = params.CertId

		if !strings.ContainsAny(cert.CertFile, "/\\") {
			server.Https.Cert = "configs/" + cert.CertFile
		} else {
			server.Https.Cert = cert.CertFile
		}

		if !strings.ContainsAny(cert.KeyFile, "/\\") {
			server.Https.Key = "configs/" + cert.KeyFile
		} else {
			server.Https.Key = cert.KeyFile
		}
	} else {
		if params.CertFile == nil {
			this.Fail("请上传证书文件")
		}
		if params.KeyFile == nil {
			this.Fail("请上传私钥文件")
		}

		// cert file
		if params.CertFile != nil {
			certFilename := "ssl." + rands.HexString(16) + params.CertFile.Ext
			_, err := params.CertFile.WriteToPath(Tea.ConfigFile(certFilename))
			if err != nil {
				this.Fail("证书文件上传失败，请检查configs/目录权限")
			}
			server.Https.Cert = "configs/" + certFilename
		}

		// key file
		if params.KeyFile != nil {
			keyFilename := "ssl." + rands.HexString(16) + params.KeyFile.Ext
			_, err := params.KeyFile.WriteToPath(Tea.ConfigFile(keyFilename))
			if err != nil {
				this.Fail("证书文件上传失败，请检查configs/目录权限")
			}
			server.Https.Key = "configs/" + keyFilename
		}
	}

	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	settings.NotifyServerChange()

	this.Next("/settings", nil).
		Success("保存成功，重启服务后生效")
}
