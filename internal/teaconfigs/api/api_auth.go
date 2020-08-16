package api

import (
	"github.com/iwind/TeaGo/maps"
	"net/http"
)

// 认证类型
const (
	APIAuthTypeNone      = "none"
	APIAuthTypeBasicAuth = "basicAuth"
	APIAuthTypeKeyAuth   = "keyAuth"
)

// 认证接口
type APIAuthInterface interface {
	// 唯一Key
	UniqueKey() string

	// 从Request中读取Key
	KeyFromRequest(req *http.Request) string

	// 匹配Request
	MatchRequest(req *http.Request) bool
}

// 新对象
func NewAPIAuth(authType string, options map[string]interface{}) APIAuthInterface {
	if authType == APIAuthTypeNone {
		return NewAPIAuthNoneAuth(options)
	}
	if authType == APIAuthTypeBasicAuth {
		return NewAPIAuthBasicAuth(options)
	}
	if authType == APIAuthTypeKeyAuth {
		return NewAPIAuthKeyAuth(options)
	}

	return nil
}

// 所有认证类型
func AllAuthTypes() []maps.Map {
	return []maps.Map{
		{
			"name": "BasicAuth",
			"code": APIAuthTypeBasicAuth,
		},
		{
			"name": "KeyAuth",
			"code": APIAuthTypeKeyAuth,
		},
	}
}

// 检查认证类型是否存在
func ContainsAuthType(authType string) bool {
	if authType == APIAuthTypeNone {
		return true
	}
	for _, auth := range AllAuthTypes() {
		if auth.GetString("code") == authType {
			return true
		}
	}
	return false
}

// 获取认证类型名
func FindAuthTypeName(authType string) string {
	if authType == "none" || len(authType) == 0 {
		return "暂无认证"
	}
	for _, auth := range AllAuthTypes() {
		if auth.GetString("code") == authType {
			return auth.GetString("name")
		}
	}
	return ""
}
