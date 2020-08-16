package teaconfigs

import (
	"crypto/tls"
	"github.com/iwind/TeaGo/maps"
)

// 认证类型
type SSLClientAuthType = int

const (
	SSLClientAuthTypeNoClientCert               SSLClientAuthType = 0
	SSLClientAuthTypeRequestClientCert          SSLClientAuthType = 1
	SSLClientAuthTypeRequireAnyClientCert       SSLClientAuthType = 2
	SSLClientAuthTypeVerifyClientCertIfGiven    SSLClientAuthType = 3
	SSLClientAuthTypeRequireAndVerifyClientCert SSLClientAuthType = 4
)

// 所有的客户端认证类型
func AllSSLClientAuthTypes() []maps.Map {
	return []maps.Map{
		{
			"name":      "不需要客户端证书",
			"type":      SSLClientAuthTypeNoClientCert,
			"requireCA": false,
		},
		{
			"name":      "请求客户端证书",
			"type":      SSLClientAuthTypeRequestClientCert,
			"requireCA": true,
		},
		{
			"name":      "需要客户端证书，但不校验",
			"type":      SSLClientAuthTypeRequireAnyClientCert,
			"requireCA": true,
		},
		{
			"name":      "有客户端证书的时候才校验",
			"type":      SSLClientAuthTypeVerifyClientCertIfGiven,
			"requireCA": true,
		},
		{
			"name":      "校验客户端提供的证书",
			"type":      SSLClientAuthTypeRequireAndVerifyClientCert,
			"requireCA": true,
		},
	}
}

// 查找单个认证方式的名称
func FindSSLClientAuthTypeName(authType SSLClientAuthType) string {
	for _, m := range AllSSLClientAuthTypes() {
		if m.GetInt("type") == authType {
			return m.GetString("name")
		}
	}
	return ""
}

// 认证类型和tls包内类型的映射
func GoSSLClientAuthType(authType SSLClientAuthType) tls.ClientAuthType {
	switch authType {
	case SSLClientAuthTypeNoClientCert:
		return tls.NoClientCert
	case SSLClientAuthTypeRequestClientCert:
		return tls.RequestClientCert
	case SSLClientAuthTypeRequireAnyClientCert:
		return tls.RequireAnyClientCert
	case SSLClientAuthTypeVerifyClientCertIfGiven:
		return tls.VerifyClientCertIfGiven
	case SSLClientAuthTypeRequireAndVerifyClientCert:
		return tls.RequireAndVerifyClientCert
	}
	return tls.NoClientCert
}
