package teawaf

import (
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/build/internal/teawaf/requests"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

// url client configure
var urlPrefixReg = regexp.MustCompile("^(?i)(http|https)://")
var httpClient = teautils.SharedHttpClient(5 * time.Second)

type BlockAction struct {
	StatusCode int    `yaml:"statusCode" json:"statusCode"`
	Body       string `yaml:"body" json:"body"` // supports HTML
	URL        string `yaml:"url" json:"url"`
}

func (this *BlockAction) Perform(waf *WAF, request *requests.Request, writer http.ResponseWriter) (allow bool) {
	if writer != nil {
		// if status code eq 444, we close the connection
		if this.StatusCode == 444 {
			hijack, ok := writer.(http.Hijacker)
			if ok {
				conn, _, _ := hijack.Hijack()
				if conn != nil {
					_ = conn.Close()
					return
				}
			}
		}

		// output response
		if this.StatusCode > 0 {
			writer.WriteHeader(this.StatusCode)
		} else {
			writer.WriteHeader(http.StatusForbidden)
		}
		if len(this.URL) > 0 {
			if urlPrefixReg.MatchString(this.URL) {
				req, err := http.NewRequest(http.MethodGet, this.URL, nil)
				if err != nil {
					logs.Error(err)
					return false
				}
				resp, err := httpClient.Do(req)
				if err != nil {
					logs.Error(err)
					return false
				}
				defer func() {
					_ = resp.Body.Close()
				}()

				for k, v := range resp.Header {
					for _, v1 := range v {
						writer.Header().Add(k, v1)
					}
				}

				buf := make([]byte, 1024)
				_, _ = io.CopyBuffer(writer, resp.Body, buf)
			} else {
				path := this.URL
				if !filepath.IsAbs(this.URL) {
					path = Tea.Root + string(os.PathSeparator) + path
				}

				data, err := ioutil.ReadFile(path)
				if err != nil {
					logs.Error(err)
					return false
				}
				_, _ = writer.Write(data)
			}
			return false
		}
		if len(this.Body) > 0 {
			_, _ = writer.Write([]byte(this.Body))
		} else {
			_, _ = writer.Write([]byte("The request is blocked by " + teaconst.TeaProductName))
		}
	}
	return false
}
