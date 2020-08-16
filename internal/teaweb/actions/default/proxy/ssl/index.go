package ssl

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/time"
	"strings"
)

type IndexAction actions.Action

// SSL设置
func (this *IndexAction) Run(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	this.Data["minVersion"] = "TLS 1.0"
	if server.SSL != nil && len(server.SSL.MinVersion) > 0 {
		this.Data["minVersion"] = server.SSL.MinVersion
	}

	this.Data["selectedTab"] = "https"
	this.Data["server"] = server
	this.Data["isHTTP"] = server.IsHTTP()
	this.Data["isTCP"] = server.IsTCP()

	this.Data["errs"] = teaproxy.SharedManager.FindServerErrors(params.ServerId)

	errorMessages := []string{}
	warningMessages := []string{}
	certs := []maps.Map{}

	notMatchedDomains := []string{}
	globalDNSNames := []string{}

	if server.SSL != nil {
		err := server.SSL.Validate()
		if err != nil {
			errorMessages = append(errorMessages, err.Error())
		}
		for index, certConfig := range server.SSL.Certs {
			info := []maps.Map{}

			// 共享的证书
			if certConfig.IsShared {
				certConfig = certConfig.FindShared()
				if certConfig == nil {
					continue
				}
			}

			// 证书是否为空
			if len(certConfig.FullCertPath()) == 0 {
				if server.SSL.On {
					errorMessages = append(errorMessages, fmt.Sprintf("证书#%d：", index+1)+"证书文件不能为空")
				}
				certs = append(certs, maps.Map{
					"config": certConfig,
					"info":   info,
				})
				continue
			}

			// 密钥是否为空
			if len(certConfig.FullKeyPath()) == 0 {
				if server.SSL.On {
					errorMessages = append(errorMessages, fmt.Sprintf("证书#%d：", index+1)+"密钥文件不能为空")
				}
				certs = append(certs, maps.Map{
					"config": certConfig,
					"info":   info,
				})
				continue
			}

			// 证书和密钥是否匹配
			cert, err := tls.LoadX509KeyPair(certConfig.FullCertPath(), certConfig.FullKeyPath())
			if err != nil {
				if server.SSL.On {
					errorMessages = append(errorMessages, fmt.Sprintf("证书#%d：", index+1)+err.Error())
				}
				certs = append(certs, maps.Map{
					"config": certConfig,
					"info":   info,
				})
				continue
			}

			allDNSNames := []string{}
			for _, data := range cert.Certificate {
				c, err := x509.ParseCertificate(data)
				if err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("证书#%d：", index+1)+err.Error())
					certs = append(certs, maps.Map{
						"config": certConfig,
						"info":   info,
					})
					continue
				}
				dnsNames := ""
				if len(c.DNSNames) > 0 {
					dnsNames = "[" + strings.Join(c.DNSNames, ", ") + "]"
					allDNSNames = append(allDNSNames, c.DNSNames...)
				}
				info = append(info, maps.Map{
					"subject":  c.Subject.CommonName + " " + dnsNames,
					"issuer":   c.Issuer.CommonName,
					"before":   timeutil.Format("Y-m-d", c.NotBefore),
					"after":    timeutil.Format("Y-m-d", c.NotAfter),
					"dnsNames": dnsNames,
				})
			}

			lists.Reverse(info)
			certs = append(certs, maps.Map{
				"config": certConfig,
				"info":   info,
			})

			if len(allDNSNames) > 0 {
				globalDNSNames = append(globalDNSNames, allDNSNames...)
			}
		}
	}

	if server.SSL != nil && server.SSL.On {
		if len(server.Name) == 0 {
			warningMessages = append(warningMessages, "当前代理服务没有设置域名，可能会导致用户访问时的域名无法和证书正确匹配。<a href=\"/proxy/update?serverId="+server.Id+"\">设置域名 &raquo;</a>")
		} else {
			// 检查domain
			for _, domain := range server.Name {
				if !teautils.MatchDomains(globalDNSNames, domain) {
					if !lists.ContainsString(notMatchedDomains, domain) {
						notMatchedDomains = append(notMatchedDomains, domain)
					}
				}
			}

			if len(notMatchedDomains) > 0 {
				message := "当前代理服务的已设置的部分域名和证书不匹配，访问以下这些域名时将不会使用证书："
				for _, domain := range notMatchedDomains {
					message += `<span class="ui label tiny">` + domain + `</span>`
				}
				message += "。<a href=\"/proxy/update?serverId=" + server.Id + "\">设置域名 &raquo;</a>"
				warningMessages = append(warningMessages, message)
			}
		}
	}

	this.Data["errorMessages"] = errorMessages
	this.Data["warningMessages"] = warningMessages
	this.Data["certs"] = certs

	// CA证书
	if server.SSL != nil {
		this.Data["clientAuthTypeName"] = teaconfigs.FindSSLClientAuthTypeName(server.SSL.ClientAuthType)

		certList := teaconfigs.SharedSSLCertList()
		caCerts := []maps.Map{}
		for _, certId := range server.SSL.ClientCACertIds {
			cert := certList.FindCert(certId)
			if cert == nil {
				continue
			}
			caCerts = append(caCerts, maps.Map{
				"description": cert.Description,
				"id":          cert.Id,
			})
		}
		this.Data["clientCACerts"] = caCerts
	} else {
		this.Data["clientAuthTypeName"] = teaconfigs.FindSSLClientAuthTypeName(teaconfigs.SSLClientAuthTypeNoClientCert)
		this.Data["clientCACerts"] = []maps.Map{}
	}

	this.Show()
}
