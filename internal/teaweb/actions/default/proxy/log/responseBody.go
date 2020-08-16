package log

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/base64"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"io/ioutil"
	"strings"
)

type ResponseBodyAction actions.Action

// 响应Body
func (this *ResponseBodyAction) Run(params struct {
	LogId string
	Day   string
}) {
	if len(params.Day) == 0 {
		params.Day = timeutil.Format("Ymd")
	}
	accessLog, err := teadb.AccessLogDAO().FindResponseHeaderAndBody(params.Day, params.LogId)
	if err != nil {
		this.Fail(err.Error())
	}
	if accessLog != nil {
		this.Data["headers"] = accessLog.SentHeader
		this.Data["isImage"] = false
		this.Data["isText"] = false
		this.Data["contentType"] = ""
		this.Data["encoding"] = ""

		if len(accessLog.ResponseBodyData) == 0 {
			this.Data["body"] = ""
			this.Data["rawBody"] = ""
		} else {
			isText := false

			// content type
			isEncoded := false
			contentTypes, ok := accessLog.SentHeader["Content-Type"]
			if ok && len(contentTypes) > 0 {
				contentType := contentTypes[0]
				semiIndex := strings.Index(contentType, ";")
				if semiIndex > -1 {
					contentType = contentType[:semiIndex]
				}
				if len(contentType) > 0 {
					this.Data["contentType"] = contentType

					if (strings.HasPrefix(contentType, "application/")) || strings.HasPrefix(contentType, "text/") {
						this.Data["isText"] = true
						isText = true
					} else if strings.HasPrefix(contentType, "image/") {
						isEncoded = true
						this.Data["isImage"] = true
						this.Data["body"] = "data:" + contentType + ";base64," + base64.StdEncoding.EncodeToString(accessLog.ResponseBodyData)

						if len(accessLog.ResponseBodyData) > 1024 {
							this.Data["rawBody"] = "[图片文件只能预览部分内容]\n" + string(accessLog.ResponseBodyData[:1024])
						} else {
							this.Data["rawBody"] = string(accessLog.ResponseBodyData)
						}
					}
				}
			}

			if !isEncoded {
				isCompressed := false
				if accessLog.SentHeader != nil {
					encodings, ok := accessLog.SentHeader["Content-Encoding"]
					if ok && len(encodings) > 0 {
						encoding := encodings[0]

						isCompressed = true
						if len(accessLog.ResponseBodyData) > 1024 {
							this.Data["rawBody"] = "[通过" + encoding + "算法压缩的内容，只能预览部分内容]\n" + string(accessLog.ResponseBodyData[:1024])
						} else {
							this.Data["rawBody"] = string(accessLog.ResponseBodyData)
						}

						this.Data["body"] = this.Data["rawBody"]
						this.Data["encoding"] = encoding

						// decode
						if encoding == "gzip" {
							reader, err := gzip.NewReader(bytes.NewReader(accessLog.ResponseBodyData))
							if err != nil {
								logs.Error(err)
							} else {
								data, err := ioutil.ReadAll(reader)
								if err != nil {
									logs.Error(err)
								} else {
									this.Data["body"] = string(data)
								}
								err = reader.Close()
								if err != nil {
									logs.Error(err)
								}
							}
						} else if encoding == "deflate" {
							reader := flate.NewReader(bytes.NewReader(accessLog.ResponseBodyData))
							data, err := ioutil.ReadAll(reader)
							if err != nil {
								logs.Error(err)
							} else {
								this.Data["body"] = string(data)
								err = reader.Close()
								if err != nil {
									logs.Error(err)
								}
							}
						} else if encoding == "br" {
							// 参考 https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Encoding
							// brotli compress, we do nothing now
							this.Data["body"] = "此内容暂时不能预览"
						} else if encoding == "compress" {
							// 参考 https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Encoding
							// do nothing now
							this.Data["body"] = "此内容暂时不能预览"
						} else {
							this.Data["body"] = "此内容暂时不能预览"
						}
					}
				}

				if !isCompressed {
					if !isText {
						if len(accessLog.ResponseBodyData) > 1024 {
							this.Data["body"] = string(accessLog.ResponseBodyData[:1024]) + "\n..."
							this.Data["rawBody"] = string(accessLog.ResponseBodyData[:1024]) + "\n..."
						} else {
							this.Data["body"] = string(accessLog.ResponseBodyData)
							this.Data["rawBody"] = string(accessLog.ResponseBodyData)
						}
					} else {
						this.Data["body"] = string(accessLog.ResponseBodyData)
						this.Data["rawBody"] = string(accessLog.ResponseBodyData)
					}
				}
			}
		}
	} else {
		this.Data["headers"] = map[string][]string{}
		this.Data["body"] = ""
		this.Data["rawBody"] = ""
	}

	this.Success()
}
