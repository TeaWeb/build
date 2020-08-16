package api

import "net/http"

// 默认认证（无认证）
type APIAuthNoneAuth struct {
}

func NewAPIAuthNoneAuth(options map[string]interface{}) APIAuthInterface {
	return &APIAuthNoneAuth{}
}

func (this *APIAuthNoneAuth) UniqueKey() string {
	return "none"
}

func (this *APIAuthNoneAuth) KeyFromRequest(req *http.Request) string {
	return ""
}

func (this *APIAuthNoneAuth) MatchRequest(req *http.Request) bool {
	return true
}
