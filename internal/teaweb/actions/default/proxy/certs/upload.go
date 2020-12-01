package certs

import (
	"bytes"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/rands"
	"io/ioutil"
	"strings"
)

type UploadAction actions.Action

// 上传证书
func (this *UploadAction) RunGet(params struct{}) {
	this.Show()
}

func (this *UploadAction) RunPost(params struct {
	Description string

	IsLocal bool

	CertFile *actions.File
	KeyFile  *actions.File

	CertPath string
	KeyPath  string

	CertType string

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

// 上传证书+私钥对
func (this *UploadAction) RunPostPair(params struct {
	Description string

	IsLocal bool

	CertFile *actions.File
	KeyFile  *actions.File

	CertPath string
	KeyPath  string

	CertType string

	On bool

	Must *actions.Must
}) {
	params.Must.
		Field("description", params.Description).
		Require("请输入证书说明")

	var certFilename string
	var keyFilename string

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

		certFilename = params.CertPath
		keyFilename = params.KeyPath
	} else {
		if params.CertFile == nil {
			this.FailField("certFile", "请选择要上传的证书文件")
		}

		if params.KeyFile == nil {
			this.FailField("keyFile", "请选择要上传的私钥文件")
		}

		certData, err := params.CertFile.Read()
		if err != nil {
			this.Fail("读取证书失败：" + err.Error())
		}

		if bytes.Contains(certData, []byte("PRIVATE KEY--")) {
			this.FailField("certFile", "证书文件不能包含密钥内容")
		}

		keyData, err := params.KeyFile.Read()
		if err != nil {
			this.Fail("读取私钥失败：" + err.Error())
		}

		if bytes.Contains(keyData, []byte("CERTIFICATE--")) {
			this.FailField("keyFile", "私钥文件不能包含证书内容")
		}

		certFilename = "ssl." + rands.HexString(16) + strings.ToLower(params.CertFile.Ext)
		err = ioutil.WriteFile(Tea.ConfigFile(certFilename), certData, 0777)
		if err != nil {
			this.Fail("保存证书失败：" + err.Error())
		}

		keyFilename = "ssl." + rands.HexString(16) + strings.ToLower(params.KeyFile.Ext)
		err = ioutil.WriteFile(Tea.ConfigFile(keyFilename), keyData, 0777)
		if err != nil {
			this.Fail("保存私钥失败：" + err.Error())
		}
	}

	cert := teaconfigs.NewSSLCertConfig(certFilename, keyFilename)
	cert.IsShared = true
	cert.Description = params.Description
	cert.IsLocal = params.IsLocal
	cert.On = params.On
	list := teaconfigs.SharedSSLCertList()
	list.AddCert(cert)
	err := list.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}

// 上传CA证书
func (this *UploadAction) RunPostCA(params struct {
	Description string

	IsLocal bool

	CertFile *actions.File
	KeyFile  *actions.File

	CertPath string
	KeyPath  string

	CertType string

	On bool

	Must *actions.Must
}) {
	params.Must.
		Field("description", params.Description).
		Require("请输入证书说明")

	var certFilename string

	if params.IsLocal {
		params.Must.
			Field("certPath", params.CertPath).
			Require("请输入证书文件路径")

		if !files.NewFile(params.CertPath).Exists() {
			this.FailField("certPath", "证书文件路径不存在")
		}

		certFilename = params.CertPath
	} else {
		if params.CertFile == nil {
			this.FailField("certFile", "请选择要上传的证书文件")
		}

		certData, err := params.CertFile.Read()
		if err != nil {
			this.Fail("读取证书失败：" + err.Error())
		}

		if bytes.Contains(certData, []byte("PRIVATE KEY--")) {
			this.FailField("certFile", "证书文件不能包含密钥内容")
		}

		certFilename = "ssl." + rands.HexString(16) + strings.ToLower(params.CertFile.Ext)
		err = ioutil.WriteFile(Tea.ConfigFile(certFilename), certData, 0777)
		if err != nil {
			this.Fail("保存证书失败：" + err.Error())
		}
	}

	cert := teaconfigs.NewSSLCertConfig(certFilename, "")
	cert.IsShared = true
	cert.IsCA = true
	cert.Description = params.Description
	cert.IsLocal = params.IsLocal
	cert.On = params.On

	list := teaconfigs.SharedSSLCertList()
	list.AddCert(cert)
	err := list.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
