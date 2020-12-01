package ssl

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/certs/certutils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type UpdateAction actions.Action

// 修改
func (this *UpdateAction) Run(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	this.Data["certs"] = []maps.Map{}
	if server.SSL != nil {
		err := server.SSL.Validate()
		if err != nil {
			logs.Error(err)
		}
		if len(server.SSL.Certs) > 0 {
			certs := []maps.Map{}
			for _, cert := range server.SSL.Certs {
				certs = append(certs, maps.Map{
					"id":          cert.Id,
					"on":          cert.On,
					"certFile":    cert.CertFile,
					"keyFile":     cert.KeyFile,
					"description": cert.Description,
					"isLocal":     cert.IsLocal,
					"isShared":    cert.IsShared,
				})
			}
			this.Data["certs"] = certs
		}
	}

	this.Data["selectedTab"] = "https"
	this.Data["server"] = server
	this.Data["isHTTP"] = server.IsHTTP()
	this.Data["isTCP"] = server.IsTCP()

	this.Data["versions"] = teaconfigs.AllTlsVersions
	if server.SSL != nil && server.SSL.HSTS != nil {
		this.Data["hsts"] = server.SSL.HSTS
	} else {
		this.Data["hsts"] = &teaconfigs.HSTSConfig{
			On:                false,
			MaxAge:            31536000,
			IncludeSubDomains: true,
			Preload:           false,
		}
	}

	this.Data["minVersion"] = "TLS 1.0"
	if server.SSL != nil && len(server.SSL.MinVersion) > 0 {
		this.Data["minVersion"] = server.SSL.MinVersion
	}

	// 加密算法套件
	this.Data["cipherSuites"] = teaconfigs.AllTLSCipherSuites
	this.Data["modernCipherSuites"] = teaconfigs.TLSModernCipherSuites
	this.Data["intermediateCipherSuites"] = teaconfigs.TLSIntermediateCipherSuites

	// 公共可以使用的证书
	this.Data["sharedCerts"] = certutils.ListPairCertsMap()
	this.Data["caCerts"] = certutils.ListCACertsMap()

	// 客户端认证
	this.Data["clientAuthTypes"] = teaconfigs.AllSSLClientAuthTypes()
	if server.SSL != nil {
		this.Data["clientAuthType"] = server.SSL.ClientAuthType
		if len(server.SSL.ClientCACertIds) == 0 {
			server.SSL.ClientCACertIds = []string{}
		}
		this.Data["clientCACertIds"] = server.SSL.ClientCACertIds
	} else {
		this.Data["clientAuthType"] = 0
		this.Data["clientCACertIds"] = []string{}
	}

	this.Show()
}

// 提交保存
func (this *UpdateAction) RunPost(params struct {
	ServerId     string
	HttpsOn      bool
	Http2Enabled bool

	Listen           []string
	CertIds          []string
	CertDescriptions []string

	CertIsLocal    []bool
	CertIsShared   []bool
	CertFilesPaths []string
	KeyFilesPaths  []string

	MinVersion     string
	CipherSuitesOn bool
	CipherSuites   []string

	HstsOn                bool
	HstsMaxAge            int
	HstsIncludeSubDomains bool
	HstsPreload           bool
	HstsDomains           []string

	ClientAuthType int
	CACertIds      []string `alias:"caCertIds"`
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	if server.SSL == nil {
		server.SSL = teaconfigs.NewSSLConfig()
	}
	server.SSL.On = params.HttpsOn
	server.SSL.HTTP2Disabled = !params.Http2Enabled
	server.SSL.Listen = teautils.FormatAddressList(params.Listen)

	if lists.ContainsString(teaconfigs.AllTlsVersions, params.MinVersion) {
		server.SSL.MinVersion = params.MinVersion
	}

	server.SSL.HSTS = &teaconfigs.HSTSConfig{
		On:                params.HstsOn,
		MaxAge:            params.HstsMaxAge,
		Domains:           params.HstsDomains,
		IncludeSubDomains: params.HstsIncludeSubDomains,
		Preload:           params.HstsPreload,
	}

	server.SSL.CipherSuites = []string{}
	if params.CipherSuitesOn {
		for _, cipherSuite := range params.CipherSuites {
			if lists.ContainsString(teaconfigs.AllTLSCipherSuites, cipherSuite) {
				server.SSL.CipherSuites = append(server.SSL.CipherSuites, cipherSuite)
			}
		}
	}

	fileBytes := map[string][]byte{} // field => []byte
	fileExts := map[string]string{}  // field => .ext
	if this.Request.MultipartForm != nil {
		for field, headers := range this.Request.MultipartForm.File {
			for _, header := range headers {
				fp, err := header.Open()
				if err != nil {
					continue
				}
				data, err := ioutil.ReadAll(fp)
				if err != nil {
					_ = fp.Close()
					continue
				}
				fileBytes[field] = data
				fileExts[field] = strings.ToLower(filepath.Ext(header.Filename))
				_ = fp.Close()

				break
			}
		}
	}

	// 证书
	certs := []*teaconfigs.SSLCertConfig{}
	for index, description := range params.CertDescriptions {
		if index >= len(params.CertIds) ||
			index >= len(params.CertIsLocal) ||
			index >= len(params.CertFilesPaths) ||
			index >= len(params.KeyFilesPaths) ||
			index >= len(params.CertIsShared) {
			continue
		}

		cert := teaconfigs.NewSSLCertConfig("", "")
		cert.Description = description
		cert.IsLocal = params.CertIsLocal[index]
		cert.IsShared = params.CertIsShared[index]

		if cert.IsShared {
			sharedCertId := this.ParamString("sharedCertIds" + fmt.Sprintf("%d", index))
			if len(sharedCertId) == 0 {
				this.Fail("请选择一个公用证书")
			}
			cert.Id = sharedCertId
		} else if cert.IsLocal {
			cert.CertFile = params.CertFilesPaths[index]
			cert.KeyFile = params.KeyFilesPaths[index]

			if len(cert.CertFile) == 0 {
				this.Fail("请输入证书#" + fmt.Sprintf("%d", index+1) + "文件路径")
			}

			if len(cert.KeyFile) == 0 {
				this.Fail("请输入证书#" + fmt.Sprintf("%d", index+1) + "私钥文件路径")
			}

			// 保留属性
			oldCert := server.SSL.FindCert(params.CertIds[index])
			if oldCert != nil {
				cert.TaskId = oldCert.TaskId
			}
		} else {
			// 兼容以前的版本（v0.1.4）
			if params.CertIds[index] == "old_version_cert" {
				cert.CertFile = server.SSL.Certificate
				cert.KeyFile = server.SSL.CertificateKey
			} else {
				// 保留先前上传的文件
				oldCert := server.SSL.FindCert(params.CertIds[index])
				if oldCert != nil {
					cert.CertFile = oldCert.CertFile
					cert.KeyFile = oldCert.KeyFile
					cert.TaskId = oldCert.TaskId
				}
			}

			{
				field := fmt.Sprintf("certFiles%d", index)
				data, ok := fileBytes[field]
				if ok {
					filename := "ssl." + rands.HexString(16) + fileExts[field]
					configFile := files.NewFile(Tea.ConfigFile(filename))
					err := configFile.Write(data)
					if err != nil {
						this.Fail(err.Error())
					}
					cert.CertFile = filename
				}
			}

			{
				field := fmt.Sprintf("keyFiles%d", index)
				data, ok := fileBytes[field]
				if ok {
					filename := "ssl." + rands.HexString(16) + fileExts[field]
					configFile := files.NewFile(Tea.ConfigFile(filename))
					err := configFile.Write(data)
					if err != nil {
						this.Fail(err.Error())
					}
					cert.KeyFile = filename
				}
			}
		}

		certs = append(certs, cert)
	}
	server.SSL.Certs = certs

	// 清除以前的版本（v0.1.4）
	server.SSL.Certificate = ""
	server.SSL.CertificateKey = ""

	// 客户端认证
	server.SSL.ClientAuthType = params.ClientAuthType
	server.SSL.ClientCACertIds = params.CACertIds

	if server.SSL.ClientAuthType != teaconfigs.SSLClientAuthTypeNoClientCert && len(server.SSL.ClientCACertIds) == 0 {
		this.Fail("已选择的客户端认证方式需要上传CA证书")
	}

	err := server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()

	this.Success()
}
