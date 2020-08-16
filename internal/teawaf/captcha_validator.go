package teawaf

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/TeaWeb/build/internal/teawaf/requests"
	"github.com/dchest/captcha"
	"github.com/iwind/TeaGo/logs"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"net/http"
	"time"
)

var captchaValidator = &CaptchaValidator{}

type CaptchaValidator struct {
}

func (this *CaptchaValidator) Run(request *requests.Request, writer http.ResponseWriter) {
	if request.Method == http.MethodPost && len(request.FormValue("TEAWEB_WAF_CAPTCHA_ID")) > 0 {
		this.validate(request, writer)
	} else {
		this.show(request, writer)
	}
}

func (this *CaptchaValidator) show(request *requests.Request, writer http.ResponseWriter) {
	// show captcha
	captchaId := captcha.NewLen(6)
	buf := bytes.NewBuffer([]byte{})
	err := captcha.WriteImage(buf, captchaId, 200, 100)
	if err != nil {
		logs.Error(err)
		return
	}

	_, _ = writer.Write([]byte(`<!DOCTYPE html>
<html>
<head>
	<title>Verify Yourself</title>
</head>
<body>
<form method="POST">
	<input type="hidden" name="TEAWEB_WAF_CAPTCHA_ID" value="` + captchaId + `"/>
	<img src="data:image/png;base64, ` + base64.StdEncoding.EncodeToString(buf.Bytes()) + `"/>` + `
	<div>
		<p>Input verify code above:</p>
		<input type="text" name="TEAWEB_WAF_CAPTCHA_CODE" maxlength="6" size="18" autocomplete="off" z-index="1" style="font-size:16px;line-height:24px; letter-spacing: 15px; padding-left: 4px"/>
	</div>
	<div>
		<button type="submit" onclick="window.location = '/webhook'" style="line-height:24px;margin-top:10px">Verify Yourself</button>
	</div>
</form>
</body>
</html>`))
}

func (this *CaptchaValidator) validate(request *requests.Request, writer http.ResponseWriter) (allow bool) {
	captchaId := request.FormValue("TEAWEB_WAF_CAPTCHA_ID")
	if len(captchaId) > 0 {
		captchaCode := request.FormValue("TEAWEB_WAF_CAPTCHA_CODE")
		if captcha.VerifyString(captchaId, captchaCode) {
			// set cookie
			timestamp := fmt.Sprintf("%d", time.Now().Unix()+CaptchaSeconds)
			m := stringutil.Md5(captchaSalt + timestamp)
			http.SetCookie(writer, &http.Cookie{
				Name:   "TEAWEB_WAF_CAPTCHA",
				Value:  m + timestamp,
				MaxAge: CaptchaSeconds,
				Path:   "/", // all of dirs
			})

			rawURL := request.URL.Query().Get("url")
			http.Redirect(writer, request.Raw(), rawURL, http.StatusSeeOther)

			return false
		} else {
			http.Redirect(writer, request.Raw(), request.URL.String(), http.StatusSeeOther)
		}
	}

	return true
}
