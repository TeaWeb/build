package teaproxy

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"io"
	"net/http"
	"os"
)

// 调用临时关闭页面
func (this *Request) callShutdown(writer *ResponseWriter) error {
	shutdown := this.shutdown
	if shutdown == nil {
		return nil
	}

	if urlPrefixRegexp.MatchString(shutdown.URL) {
		return this.callURL(writer, http.MethodGet, shutdown.URL, "")
	} else {
		file := Tea.Root + Tea.DS + shutdown.URL
		fp, err := os.Open(file)
		if err != nil {
			logs.Error(err)
			msg := "404 page not found: '" + shutdown.URL + "'"

			writer.WriteHeader(http.StatusNotFound)
			_, err = writer.Write([]byte(msg))
			if err != nil {
				logs.Error(err)
			}
			return err
		}

		// 自定义响应Headers
		if shutdown.Status > 0 {
			this.WriteResponseHeaders(writer, shutdown.Status)
			writer.WriteHeader(shutdown.Status)
		} else {
			this.WriteResponseHeaders(writer, http.StatusOK)
			writer.WriteHeader(http.StatusOK)
		}
		buf := bytePool1k.Get()
		_, err = io.CopyBuffer(writer, fp, buf)
		bytePool1k.Put(buf)
		err = fp.Close()
		if err != nil {
			logs.Error(err)
		}

		return err
	}
}
