package teaproxy

import (
	"errors"
	"fmt"
	"github.com/TeaWeb/build/internal/teaplugins"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/gofcgi"
	"io"
	"net"
	"net/url"
	"path/filepath"
	"strings"
)

// 调用Fastcgi
func (this *Request) callFastcgi(writer *ResponseWriter) error {
	env := this.fastcgi.FilterParams(this.raw)
	if len(this.root) > 0 {
		if !env.Has("DOCUMENT_ROOT") {
			env["DOCUMENT_ROOT"] = this.root
		}
	}
	if !env.Has("REMOTE_ADDR") {
		env["REMOTE_ADDR"] = this.raw.RemoteAddr
	}
	if !env.Has("QUERY_STRING") {
		u, err := url.ParseRequestURI(this.uri)
		if err == nil {
			env["QUERY_STRING"] = u.RawQuery
		} else {
			env["QUERY_STRING"] = this.raw.URL.RawQuery
		}
	}
	if !env.Has("SERVER_NAME") {
		env["SERVER_NAME"] = this.host
	}
	if !env.Has("REQUEST_URI") {
		env["REQUEST_URI"] = this.uri
	}
	if !env.Has("HOST") {
		env["HOST"] = this.host
	}

	if len(this.serverAddr) > 0 {
		if !env.Has("SERVER_ADDR") {
			env["SERVER_ADDR"] = this.serverAddr
		}
		if !env.Has("SERVER_PORT") {
			_, port, err := net.SplitHostPort(this.serverAddr)
			if err == nil {
				env["SERVER_PORT"] = port
			}
		}
	}

	// 连接池配置
	poolSize := this.fastcgi.PoolSize
	if poolSize <= 0 {
		poolSize = 16
	}

	client, err := gofcgi.SharedPool(this.fastcgi.Network(), this.fastcgi.Address(), uint(poolSize)).Client()
	if err != nil {
		this.serverError(writer)
		logs.Error(errors.New("fastcgi: " + err.Error()))
		this.addError(errors.New("fastcgi: " + err.Error()))
		return nil
	}

	// 请求相关
	if !env.Has("REQUEST_METHOD") {
		env["REQUEST_METHOD"] = this.method
	}
	if !env.Has("CONTENT_LENGTH") {
		env["CONTENT_LENGTH"] = fmt.Sprintf("%d", this.raw.ContentLength)
	}
	if !env.Has("CONTENT_TYPE") {
		env["CONTENT_TYPE"] = this.raw.Header.Get("Content-Type")
	}

	// 处理SCRIPT_FILENAME
	scriptPath := env.GetString("SCRIPT_FILENAME")
	if len(scriptPath) > 0 && (strings.Index(scriptPath, "/") < 0 && strings.Index(scriptPath, "\\") < 0) {
		env["SCRIPT_FILENAME"] = env.GetString("DOCUMENT_ROOT") + Tea.DS + scriptPath
	}
	scriptFilename := filepath.Base(this.raw.URL.Path)

	// PATH_INFO
	pathInfoReg := this.fastcgi.PathInfoRegexp()
	pathInfo := ""
	if pathInfoReg != nil {
		matches := pathInfoReg.FindStringSubmatch(this.raw.URL.Path)
		countMatches := len(matches)
		if countMatches == 1 {
			pathInfo = matches[0]
		} else if countMatches == 2 {
			pathInfo = matches[1]
		} else if countMatches > 2 {
			scriptFilename = matches[1]
			pathInfo = matches[2]
		}

		if !env.Has("PATH_INFO") {
			env["PATH_INFO"] = pathInfo
		}
	}

	this.addVarMapping(map[string]string{
		"fastcgi.documentRoot": env.GetString("DOCUMENT_ROOT"),
		"fastcgi.filename":     scriptFilename,
		"fastcgi.pathInfo":     pathInfo,
	})

	params := map[string]string{}
	for key, value := range env {
		params[key] = this.Format(types.String(value))
	}

	for k, v := range this.raw.Header {
		if k == "Connection" {
			continue
		}
		for _, subV := range v {
			params["HTTP_"+strings.ToUpper(strings.Replace(k, "-", "_", -1))] = subV
		}
	}

	// 自定义请求Header
	if len(this.requestHeaders) > 0 {
		for _, header := range this.requestHeaders {
			if !header.On {
				continue
			}
			v := header.Value
			if header.HasVariables() {
				v = this.Format(v)
			}
			params["HTTP_"+strings.ToUpper(strings.Replace(header.Name, "-", "_", -1))] = v
		}
	}

	host, found := params["HTTP_HOST"]
	if !found || len(host) == 0 {
		params["HTTP_HOST"] = this.host
	}

	fcgiReq := gofcgi.NewRequest()
	fcgiReq.SetTimeout(this.fastcgi.ReadTimeoutDuration())
	fcgiReq.SetParams(params)
	fcgiReq.SetBody(this.raw.Body, uint32(this.requestLength()))

	resp, stderr, err := client.Call(fcgiReq)
	if err != nil {
		this.serverError(writer)
		//if this.debug {
		logs.Error(err)
		this.addError(err)
		//}
		return nil
	}

	if len(stderr) > 0 {
		err := errors.New("Fastcgi Error: " + strings.TrimSpace(string(stderr)) + " script: " + maps.NewMap(params).GetString("SCRIPT_FILENAME"))
		logs.Error(err)
		this.addError(err)
	}

	defer resp.Body.Close()

	// 忽略的Header
	ignoreHeaders := this.convertIgnoreHeaders()
	hasIgnoreHeaders := ignoreHeaders.Len() > 0

	// 设置Header
	var hasCharset = len(this.charset) > 0
	for k, v := range resp.Header {
		if k == "Connection" {
			continue
		}
		if hasIgnoreHeaders && ignoreHeaders.Has(strings.ToUpper(k)) {
			continue
		}

		for _, subV := range v {
			// 字符集
			if hasCharset && k == "Content-Type" {
				if _, found := textMimeMap[subV]; found {
					if !strings.Contains(subV, "charset=") {
						subV += "; charset=" + this.charset
					}
				}
			}
			writer.Header().Add(k, subV)
		}
	}

	// 自定义Header
	this.WriteResponseHeaders(writer, resp.StatusCode)

	// 插件过滤
	if teaplugins.HasResponseFilters {
		resp.Header = writer.Header()
		resp = teaplugins.FilterResponse(resp)

		// reset headers
		oldHeaders := writer.Header()
		for key := range oldHeaders {
			oldHeaders.Del(key)
		}

		for key, value := range resp.Header {
			for _, v := range value {
				oldHeaders.Add(key, v)
			}
		}
	}

	// 准备
	writer.Prepare(resp.ContentLength)

	// 设置响应码
	writer.WriteHeader(resp.StatusCode)

	// 输出内容
	pool := this.bytePool(resp.ContentLength)
	buf := pool.Get()
	_, err = io.CopyBuffer(writer, resp.Body, buf)
	pool.Put(buf)
	if err != nil {
		logs.Error(err)
		this.addError(err)
		return nil
	}

	return nil
}
