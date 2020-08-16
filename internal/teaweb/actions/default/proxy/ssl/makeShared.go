package ssl

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type MakeSharedAction actions.Action

// 让证书共享化，即加入公共的证书管理中
func (this *MakeSharedAction) RunPost(params struct {
	ServerId string
	CertId   string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到要操作的代理服务")
	}

	if server.SSL == nil {
		this.Fail("还没有配置SSL")
	}

	cert := server.SSL.FindCert(params.CertId)
	if cert == nil {
		this.Fail("找不到要操作的证书")
	}

	list := teaconfigs.SharedSSLCertList()
	if list.FindCert(params.CertId) != nil {
		this.Fail("已经在证书管理中，无需重复添加")
	}

	cert.IsShared = true
	list.AddCert(cert)
	err := list.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
