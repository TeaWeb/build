package teawaf

import (
	"github.com/TeaWeb/build/internal/teawaf/requests"
	"net/http"
)

type ActionString = string

const (
	ActionLog     = "log"     // allow and log
	ActionBlock   = "block"   // block
	ActionCaptcha = "captcha" // block and show captcha
	ActionAllow   = "allow"   // allow
	ActionGoGroup = "go_group" // go to next rule group
	ActionGoSet   = "go_set"   // go to next rule set
)

type ActionInterface interface {
	Perform(waf *WAF, request *requests.Request, writer http.ResponseWriter) (allow bool)
}
