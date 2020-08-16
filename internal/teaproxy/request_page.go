package teaproxy

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"io"
	"net/http"
	"os"
	"regexp"
)

var urlPrefixRegexp = regexp.MustCompile("^(?i)(http|https|ftp)://")

func (this *Request) callPage(writer *ResponseWriter, status int) (shouldStop bool) {
	if len(this.pages) == 0 {
		return false
	}

	for _, page := range this.pages {
		if page.Match(status) {
			if urlPrefixRegexp.MatchString(page.URL) {
				err := this.callURL(writer, http.MethodGet, page.URL, "")
				if err != nil {
					logs.Error(err)
				}
				return true
			} else {
				file := Tea.Root + Tea.DS + page.URL
				fp, err := os.Open(file)
				if err != nil {
					logs.Error(err)
					msg := "404 page not found: '" + page.URL + "'"

					writer.WriteHeader(http.StatusNotFound)
					_, err := writer.Write([]byte(msg))
					if err != nil {
						logs.Error(err)
					}
					return true
				}

				// 修改状态码
				if page.NewStatus > 0 {
					// 自定义响应Headers
					this.WriteResponseHeaders(writer, page.NewStatus)
					writer.WriteHeader(page.NewStatus)
				} else {
					this.WriteResponseHeaders(writer, status)
					writer.WriteHeader(status)
				}
				buf := bytePool1k.Get()
				_, err = io.CopyBuffer(writer, fp, buf)
				bytePool1k.Put(buf)
				if err != nil {
					logs.Error(err)
				}
				err = fp.Close()
				if err != nil {
					logs.Error(err)
				}
			}

			return true
		}
	}
	return false
}
