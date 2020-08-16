package api

import (
	"bytes"
	"github.com/iwind/TeaGo/types"
	"io/ioutil"
	"net/http"
)

// Key Auth
type APIAuthKeyAuth struct {
	Key         string
	HeaderField string
	FormField   string
}

func NewAPIAuthKeyAuth(options map[string]interface{}) APIAuthInterface {
	auth := &APIAuthKeyAuth{}

	key, found := options["key"]
	if found {
		auth.Key = types.String(key)
	}

	headerField, found := options["headerField"]
	if found {
		auth.HeaderField = types.String(headerField)
	}

	formField, found := options["formField"]
	if found {
		auth.FormField = types.String(formField)
	}

	return auth
}

func (this *APIAuthKeyAuth) UniqueKey() string {
	return this.Key
}

func (this *APIAuthKeyAuth) KeyFromRequest(req *http.Request) string {
	return ""
}

func (this *APIAuthKeyAuth) MatchRequest(req *http.Request) bool {
	key := ""

	if len(this.HeaderField) > 0 {
		key = req.Header.Get(this.HeaderField)
		if key == this.Key {
			return true
		}
	}

	if len(this.FormField) > 0 {
		key = req.URL.Query().Get(this.FormField)
		if key == this.Key {
			return true
		}

		if req.Method == http.MethodPost {
			// TODO 需要提升性能
			data, err := ioutil.ReadAll(req.Body)
			if err == nil {
				newReq, err := http.NewRequest(req.Method, req.URL.String(), bytes.NewReader(data))
				if err == nil {
					newReq.Header = req.Header
					newReq.ParseForm()
					key = newReq.FormValue(this.FormField)
				}

				req.Body = ioutil.NopCloser(bytes.NewReader(data))
				if key == this.Key {
					return true
				}
			}
		}
	}

	return false
}
