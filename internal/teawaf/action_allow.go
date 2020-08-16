package teawaf

import (
	"github.com/TeaWeb/build/internal/teawaf/requests"
	"net/http"
)

type AllowAction struct {
}

func (this *AllowAction) Perform(waf *WAF, request *requests.Request, writer http.ResponseWriter) (allow bool) {
	// do nothing
	return true
}
