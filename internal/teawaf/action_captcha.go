package teawaf

import (
	"github.com/TeaWeb/build/internal/teawaf/requests"
	"github.com/iwind/TeaGo/rands"
	"github.com/iwind/TeaGo/types"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"net/http"
	"net/url"
	"time"
)

var captchaSalt = rands.HexString(16)

const (
	CaptchaSeconds = 600 // 10 minutes
)

type CaptchaAction struct {
}

func (this *CaptchaAction) Perform(waf *WAF, request *requests.Request, writer http.ResponseWriter) (allow bool) {
	// TEAWEB_CAPTCHA:
	cookie, err := request.Cookie("TEAWEB_WAF_CAPTCHA")
	if err == nil && cookie != nil && len(cookie.Value) > 32 {
		m := cookie.Value[:32]
		timestamp := cookie.Value[32:]
		if stringutil.Md5(captchaSalt+timestamp) == m && time.Now().Unix() < types.Int64(timestamp) { // verify md5
			return true
		}
	}

	refURL := request.URL.String()
	if len(request.Referer()) > 0 {
		refURL = request.Referer()
	}
	http.Redirect(writer, request.Raw(), "/WAFCAPTCHA?url="+url.QueryEscape(refURL), http.StatusTemporaryRedirect)

	return false
}
