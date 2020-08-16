package certs

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/certs/certutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"time"
)

type IndexAction actions.Action

// 证书列表
func (this *IndexAction) RunGet(params struct{}) {
	certs := lists.Map(teaconfigs.SharedSSLCertList().Certs, func(k int, v interface{}) interface{} {
		cert := v.(*teaconfigs.SSLCertConfig)
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
		commonName := ""
		if len(errorString) == 0 {
			commonName = cert.Issuer().CommonName
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

		countRef := len(certutils.FindAllServersUsingCert(cert.Id))

		return maps.Map{
			"id":               cert.Id,
			"on":               cert.On,
			"description":      cert.Description,
			"error":            errorString,
			"dayBefore":        dayBefore,
			"dayAfter":         dayAfter,
			"dnsNames":         cert.DNSNames(),
			"countRef":         countRef,
			"isActive":         isActive,
			"isExpired":        isExpired,
			"isExpiring7Days":  isExpiring7Days,
			"isExpiring30Days": isExpiring30Days,
			"commonName":       commonName,
			"isACME":           len(cert.TaskId) > 0,
			"isCA":             cert.IsCA,
		}
	})
	this.Data["certs"] = certs

	this.Show()
}
