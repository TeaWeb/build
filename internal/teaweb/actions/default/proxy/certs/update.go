package certs

import (
	"bytes"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"io/ioutil"
	"strings"
)

type UpdateAction actions.Action

// 修改
func (this *UpdateAction) RunGet(params struct {
	CertId string
}) {
	list := teaconfigs.SharedSSLCertList()
	cert := list.FindCert(params.CertId)
	if cert == nil {
		this.Fail("找不到要修改的证书")
	}

	certTime := ""
	keyTime := ""

	if len(cert.CertFile) > 0 {
		stat, _ := files.NewFile(Tea.ConfigFile(cert.CertFile)).Stat()
		if stat != nil {
			certTime = timeutil.Format("Y-m-d H:i:s", stat.ModTime)
		}
	}

	if len(cert.KeyFile) > 0 {
		stat, _ := files.NewFile(Tea.ConfigFile(cert.KeyFile)).Stat()
		if stat != nil {
			keyTime = timeutil.Format("Y-m-d H:i:s", stat.ModTime)
		}
	}

	this.Data["cert"] = maps.Map{
		"id":          cert.Id,
		"on":          cert.On,
		"description": cert.Description,
		"certFile":    cert.CertFile,
		"keyFile":     cert.KeyFile,
		"certTime":    certTime,
		"keyTime":     keyTime,
		"isLocal":     cert.IsLocal,
		"isCA":        cert.IsCA,
	}

	this.Show()
}

// 提交修改
func (this *UpdateAction) RunPost(params struct {
	CertId      string
	Description string

	IsLocal  bool
	CertType string

	CertFile *actions.File
	KeyFile  *actions.File

	CertPath string
	KeyPath  string

	On bool

	Must *actions.Must
}) {
	if params.CertType == "pair" {
		this.RunPostPair(params)
	} else if params.CertType == "ca" {
		this.RunPostCA(params)
	} else {
		this.Fail("请选择正确的证书类型")
	}
}

// 证书+私钥
func (this *UpdateAction) RunPostPair(params struct {
	CertId      string
	Description string

	IsLocal  bool
	CertType string

	CertFile *actions.File
	KeyFile  *actions.File

	CertPath string
	KeyPath  string

	On bool

	Must *actions.Must
}) {
	params.Must.
		Field("description", params.Description).
		Require("请输入证书说明")

	list := teaconfigs.SharedSSLCertList()
	cert := list.FindCert(params.CertId)
	if cert == nil {
		this.Fail("找不到要修改的证书")
	}

	cert.IsLocal = params.IsLocal
	cert.IsCA = false

	if params.IsLocal {
		params.Must.
			Field("certPath", params.CertPath).
			Require("请输入证书文件路径").
			Field("keyPath", params.KeyPath).
			Require("请输入私钥文件路径")

		if !files.NewFile(params.CertPath).Exists() && !files.NewFile(Tea.ConfigFile(params.CertPath)).Exists() {
			this.FailField("certPath", "证书文件路径不存在")
		}

		if !files.NewFile(params.KeyPath).Exists() && !files.NewFile(Tea.ConfigFile(params.KeyPath)).Exists() {
			this.FailField("keyPath", "私钥文件路径不存在")
		}

		cert.CertFile = params.CertPath
		cert.KeyFile = params.KeyPath
	} else {
		if params.CertFile != nil {
			certData, err := params.CertFile.Read()
			if err != nil {
				this.Fail("读取证书失败：" + err.Error())
			}

			if bytes.Contains(certData, []byte("PRIVATE KEY--")) {
				this.FailField("certFile", "证书文件不能包含密钥内容")
			}

			cert.CertFile = "ssl." + rands.HexString(16) + strings.ToLower(params.CertFile.Ext)
			err = ioutil.WriteFile(Tea.ConfigFile(cert.CertFile), certData, 0777)
			if err != nil {
				this.Fail("保存证书失败：" + err.Error())
			}
		}

		if params.KeyFile != nil {
			keyData, err := params.KeyFile.Read()
			if err != nil {
				this.Fail("读取私钥失败：" + err.Error())
			}

			if bytes.Contains(keyData, []byte("CERTIFICATE--")) {
				this.FailField("keyFile", "私钥文件不能包含证书内容")
			}

			cert.KeyFile = "ssl." + rands.HexString(16) + strings.ToLower(params.KeyFile.Ext)
			err = ioutil.WriteFile(Tea.ConfigFile(cert.KeyFile), keyData, 0777)
			if err != nil {
				this.Fail("保存私钥失败：" + err.Error())
			}
		}
	}

	cert.IsShared = true
	cert.Description = params.Description
	cert.On = params.On

	err := list.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 通知更新
	proxyutils.NotifyChange()

	this.Success()
}

// 证书+私钥
func (this *UpdateAction) RunPostCA(params struct {
	CertId      string
	Description string

	IsLocal  bool
	CertType string

	CertFile *actions.File
	KeyFile  *actions.File

	CertPath string
	KeyPath  string

	On bool

	Must *actions.Must
}) {
	params.Must.
		Field("description", params.Description).
		Require("请输入证书说明")

	list := teaconfigs.SharedSSLCertList()
	cert := list.FindCert(params.CertId)
	if cert == nil {
		this.Fail("找不到要修改的证书")
	}

	cert.IsLocal = params.IsLocal
	cert.IsCA = true

	if params.IsLocal {
		params.Must.
			Field("certPath", params.CertPath).
			Require("请输入证书文件路径")

		if !files.NewFile(params.CertPath).Exists() {
			this.FailField("certPath", "证书文件路径不存在")
		}

		cert.CertFile = params.CertPath
	} else {
		if params.CertFile != nil {
			certData, err := params.CertFile.Read()
			if err != nil {
				this.Fail("读取证书失败：" + err.Error())
			}

			if bytes.Contains(certData, []byte("PRIVATE KEY--")) {
				this.FailField("certFile", "证书文件不能包含密钥内容")
			}

			cert.CertFile = "ssl." + rands.HexString(16) + strings.ToLower(params.CertFile.Ext)
			err = ioutil.WriteFile(Tea.ConfigFile(cert.CertFile), certData, 0777)
			if err != nil {
				this.Fail("保存证书失败：" + err.Error())
			}
		}
	}

	cert.IsShared = true
	cert.Description = params.Description
	cert.On = params.On

	err := list.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 通知更新
	proxyutils.NotifyChange()

	this.Success()
}
