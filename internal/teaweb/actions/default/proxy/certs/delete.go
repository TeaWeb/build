package certs

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/certs/certutils"
	"github.com/iwind/TeaGo/actions"
)

type DeleteAction actions.Action

// 删除证书
func (this *DeleteAction) RunPost(params struct {
	CertId string
}) {
	list := teaconfigs.SharedSSLCertList()
	cert := list.FindCert(params.CertId)
	if cert == nil {
		this.Fail("找不到要删除的证书")
	}

	if len(certutils.FindAllServersUsingCert(params.CertId)) > 0 {
		this.Fail("此证书正在被使用不能删除")
	}

	list.RemoveCert(params.CertId)
	err := list.Save()
	if err != nil {
		this.Fail("删除失败：" + err.Error())
	}

	this.Success()
}
