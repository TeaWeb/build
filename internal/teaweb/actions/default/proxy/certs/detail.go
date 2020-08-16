package certs

import (
	"crypto/x509"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/certs/certutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"time"
)

type DetailAction actions.Action

// 详情
func (this *DetailAction) RunGet(params struct {
	CertId string
}) {
	list := teaconfigs.SharedSSLCertList()
	cert := list.FindCert(params.CertId)
	if cert == nil {
		this.Fail("找不到要查看的证书")
	}

	errorString := ""

	err := cert.Validate()
	if err != nil {
		errorString = err.Error()
	}

	dayBefore := ""
	dayAfter := ""
	isActive := false
	isExpired := false
	isExpiring7Days := false
	isExpiring30Days := false
	commonNames := []string{}
	if len(errorString) == 0 {
		for _, data := range cert.CertObject().Certificate {
			c, err := x509.ParseCertificate(data)
			if err != nil {
				continue
			}
			commonNames = append(commonNames, c.Subject.CommonName)
		}

		dayBefore = timeutil.Format("Y-m-d", cert.TimeBefore())
		dayAfter = timeutil.Format("Y-m-d", cert.TimeAfter())

		if cert.TimeAfter().After(time.Now()) {
			isActive = true

			days := int(time.Since(cert.TimeAfter()).Hours() / 24)
			if days >= -7 {
				isExpiring7Days = true
				isExpiring30Days = true
			} else if days >= -30 {
				isExpiring30Days = true
			}
		} else {
			isExpired = true
		}
	}

	lists.Reverse(commonNames)

	references := []maps.Map{}
	for _, server := range certutils.FindAllServersUsingCert(cert.Id) {
		references = append(references, maps.Map{
			"name": server.Description,
			"link": "/proxy/ssl?serverId=" + server.Id,
		})
	}

	certBodyBytes, _ := cert.ReadCert()
	keyBodyBytes, _ := cert.ReadKey()

	this.Data["cert"] = maps.Map{
		"id":               cert.Id,
		"on":               cert.On,
		"description":      cert.Description,
		"error":            errorString,
		"dayBefore":        dayBefore,
		"dayAfter":         dayAfter,
		"dnsNames":         cert.DNSNames(),
		"references":       references,
		"isActive":         isActive,
		"isExpired":        isExpired,
		"isExpiring7Days":  isExpiring7Days,
		"isExpiring30Days": isExpiring30Days,
		"commonNames":      commonNames,
		"certText":         string(certBodyBytes),
		"keyText":          string(keyBodyBytes),
		"isLocal":          cert.IsLocal,
		"certFile":         cert.CertFile,
		"keyFile":          cert.KeyFile,
		"isCA":             cert.IsCA,
		"isACME":           len(cert.TaskId) > 0,
	}

	this.Show()
}
