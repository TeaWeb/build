package api

import (
	"encoding/base64"
	"github.com/iwind/TeaGo/types"
	"net/http"
	"strings"
)

// Basic Auth
type APIAuthBasicAuth struct {
	Username string
	Password string
}

func NewAPIAuthBasicAuth(options map[string]interface{}) APIAuthInterface {
	auth := &APIAuthBasicAuth{}
	username, found := options["username"]
	if found {
		auth.Username = types.String(username)
	}

	password, found := options["password"]
	if found {
		auth.Password = types.String(password)
	}
	return auth
}

func (this *APIAuthBasicAuth) UniqueKey() string {
	return this.Username
}

func (this *APIAuthBasicAuth) KeyFromRequest(req *http.Request) string {
	return ""
}

func (this *APIAuthBasicAuth) MatchRequest(req *http.Request) bool {
	authorization := req.Header.Get("Authorization")
	if len(authorization) == 0 {
		return false
	}

	pieces := strings.SplitN(authorization, " ", 2)
	if len(pieces) == 1 {
		return false
	}

	authType := pieces[0]
	if authType != "Basic" {
		return false
	}
	data := strings.TrimSpace(pieces[1])
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return false
	}

	info := string(b)
	pieces = strings.SplitN(info, ":", 2)
	if len(pieces) == 1 {
		return false
	}
	return pieces[0] == this.Username && pieces[1] == this.Password
}
