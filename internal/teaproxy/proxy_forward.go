package teaproxy

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/TeaWeb/build/internal/teaproxy/mitm"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

type ForwardProxy struct {
	req    *Request
	writer *ResponseWriter
}

func (this *ForwardProxy) forwardHTTP() error {
	defer this.req.log()

	// watch
	if this.req.isWatching {
		// 判断如果Content-Length过长，则截断
		reqData, err := httputil.DumpRequest(this.req.raw, true)
		if err == nil {
			if len(reqData) > 100240 {
				reqData = reqData[:100240]
			}
			this.req.requestData = reqData
		}

		this.writer.SetBodyCopying(true)
	} else {
		max := 512 * 1024 // 512K
		if this.req.accessLog != nil && lists.ContainsInt(this.req.accessLog.Fields, accesslogs.AccessLogFieldRequestBody) {
			body, err := ioutil.ReadAll(this.req.raw.Body)
			if err == nil {
				if len(body) > max {
					this.req.requestData = body[:max]
				} else {
					this.req.requestData = body
				}
			}
			this.req.raw.Body = ioutil.NopCloser(bytes.NewReader(body))
		}
		if this.req.accessLog != nil && lists.ContainsInt(this.req.accessLog.Fields, accesslogs.AccessLogFieldResponseBody) {
			this.writer.SetBodyCopying(true)
		}
	}

	this.req.raw.RequestURI = ""

	// 删除代理相关Header
	for n, _ := range this.req.raw.Header {
		if lists.ContainsString([]string{
			"Connection",
			"Accept-Encoding",
			"Proxy-Connection",
			"Proxy-Authenticate",
			"Proxy-Authorization",
		}, n) {
			this.req.raw.Header.Del(n)
		}
	}

	client := teautils.SharedHttpClient(30 * time.Second)
	resp, err := client.Do(this.req.raw)
	if err != nil {
		this.req.serverError(this.writer)
		this.req.addError(err)
		return nil
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	for k, v := range resp.Header {
		if k == "Connection" {
			continue
		}
		for _, subV := range v {
			this.writer.Header().Add(k, subV)
		}
	}

	this.writer.Prepare(resp.ContentLength)
	this.writer.WriteHeader(resp.StatusCode)

	_, _ = io.Copy(this.writer, resp.Body)
	return nil
}

func (this *ForwardProxy) forwardConnect() error {
	defer this.req.log()

	hijacker, ok := this.writer.writer.(http.Hijacker)
	if !ok {
		this.req.serverError(this.writer)
		return nil
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		this.req.serverError(this.writer)
		this.req.addError(err)
		return nil
	}

	_, _ = clientConn.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))

	hostConn, err := net.DialTimeout("tcp", this.req.host, 30*time.Second)
	if err != nil {
		this.req.serverError(this.writer)
		this.req.addError(err)
		return nil
	}

	go func() {
		_, _ = io.Copy(clientConn, hostConn)
		_ = clientConn.Close()
		_ = hostConn.Close()
	}()
	go func() {
		_, _ = io.Copy(hostConn, clientConn)
		_ = clientConn.Close()
		_ = hostConn.Close()
	}()
	return nil
}

func (this *ForwardProxy) forwardMitm() error {
	hijacker, ok := this.writer.writer.(http.Hijacker)
	this.req.log()
	if !ok {
		this.req.serverError(this.writer)
		return nil
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		this.req.serverError(this.writer)
		this.req.addError(err)
		return nil
	}

	_, _ = clientConn.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))

	hostName, _, _ := net.SplitHostPort(this.req.host)
	if len(hostName) == 0 {
		hostName = this.req.host
	}

	mitmLocker.Lock()
	clientConfig, ok := mitmCache[hostName]
	mitmLocker.Unlock()

	if !ok {
		certData, err := ioutil.ReadFile(Tea.Root + "/web/certs/teaweb.proxy.pem")
		if err != nil {
			logs.Error(err)
			return nil
		}

		certBlock, _ := pem.Decode(certData)
		if err != nil {
			logs.Error(err)
			return nil
		}

		rootCert, err := x509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			logs.Error(err)
			return nil
		}

		keyData, err := ioutil.ReadFile(Tea.Root + "/web/certs/teaweb.proxy.key")
		if err != nil {
			logs.Error(err)
			return nil
		}

		block, _ := pem.Decode(keyData)
		if err != nil {
			logs.Error(err)
			return nil
		}

		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			logs.Error(err)
			return nil
		}

		config, err := mitm.NewConfig(rootCert, privateKey)
		if err != nil {
			this.req.serverError(this.writer)
			this.req.addError(err)
			return nil
		}

		mitmLocker.Lock()
		clientConfig = config.TLSForHost(hostName)
		clientConfig.ServerName = hostName
		mitmCache[hostName] = clientConfig

		// 清理
		if len(mitmCache) >= 100000 {
			mitmCache = map[string]*tls.Config{}
		}
		mitmLocker.Unlock()
	}

	hostConn, err := tls.Dial("tcp", this.req.host, nil)
	if err != nil {
		this.req.serverError(this.writer)
		this.req.addError(err)
		return nil
	}

	client := tls.Server(clientConn, clientConfig)

	var closer = func() {
		_ = hostConn.Close()
		_ = clientConn.Close()
	}
	defer closer()

	clientReader := bufio.NewReader(client)
	hostReader := bufio.NewReader(hostConn)

	accessLogOn := this.req.accessLog != nil && this.req.accessLog.On

	for {
		rawReq, err := http.ReadRequest(clientReader)
		if err != nil {
			closer()
			break
		}
		this.setProxyHeaders(rawReq)

		var req *Request = nil
		logResponse := false
		if accessLogOn {
			req = NewRequest(rawReq)

			req.responseWriter = NewResponseWriter(NewEmptyResponseWriter())
			req.accessLog = this.req.accessLog
			req.enableStat = this.req.enableStat
			req.host = this.req.host
			req.method = rawReq.Method
			req.uri = rawReq.URL.RequestURI()
			req.rawScheme = "https"
			req.scheme = "https" // 转发后的scheme
			req.serverName = this.req.serverName
			req.serverAddr = this.req.serverAddr

			_ = req.configure(this.req.server, 0, false)

			// watch
			if this.req.isWatching {
				// 判断如果Content-Length过长，则截断
				reqData, err := httputil.DumpRequest(req.raw, true)
				if err == nil {
					if len(reqData) > 100240 {
						reqData = reqData[:100240]
					}
					req.requestData = reqData
				}

				logResponse = true
			} else {
				max := 512 * 1024 // 512K
				if req.accessLog != nil && lists.ContainsInt(req.accessLog.Fields, accesslogs.AccessLogFieldRequestBody) {
					body, err := ioutil.ReadAll(req.raw.Body)
					if err == nil {
						if len(body) > max {
							req.requestData = body[:max]
						} else {
							req.requestData = body
						}
					}
					req.raw.Body = ioutil.NopCloser(bytes.NewReader(body))
				}
				if req.accessLog != nil && lists.ContainsInt(this.req.accessLog.Fields, accesslogs.AccessLogFieldResponseBody) {
					logResponse = true
				}
			}
		}

		err = rawReq.Write(hostConn)
		if err != nil {
			if accessLogOn {
				req.addError(err)
			}
			closer()
			break
		}

		resp, err := http.ReadResponse(hostReader, nil)
		if err != nil {
			if accessLogOn {
				req.addError(err)
				req.log()
			}

			closer()
			break
		}

		if accessLogOn {
			req.responseWriter.Prepare(resp.ContentLength)
			req.responseWriter.WriteHeader(resp.StatusCode)
			req.responseWriter.AddHeaders(resp.Header)

			if logResponse {
				req.responseWriter.SetBodyCopying(true)

				bodyData, _ := ioutil.ReadAll(resp.Body)
				req.responseWriter.body = bodyData
				resp.Body = ioutil.NopCloser(bytes.NewReader(bodyData))
			}
		}

		err = resp.Write(client)
		if err != nil {
			if accessLogOn {
				req.addError(err)
				req.log()
			}
			closer()
			break
		}

		req.log()
	}

	return nil
}

// 设置代理相关头部信息
// 参考：https://tools.ietf.org/html/rfc7239
func (this *ForwardProxy) setProxyHeaders(rawRequest *http.Request) {
	remoteAddr := this.req.raw.RemoteAddr
	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil {
		remoteAddr = host
	}

	// x-real-ip
	{
		_, ok1 := rawRequest.Header["X-Real-IP"]
		_, ok2 := rawRequest.Header["X-Real-Ip"]
		if !ok1 && !ok2 {
			rawRequest.Header["X-Real-IP"] = []string{remoteAddr}
		}
	}

	// X-Forwarded-For
	{
		forwardedFor, ok := rawRequest.Header["X-Forwarded-For"]
		if ok {
			rawRequest.Header["X-Forwarded-For"] = []string{strings.Join(forwardedFor, ", ") + ", " + remoteAddr}
		} else {
			rawRequest.Header["X-Forwarded-For"] = []string{remoteAddr}
		}
	}

	// others
	rawRequest.Header.Set("X-Forwarded-By", this.req.serverAddr)

	if _, ok := rawRequest.Header["X-Forwarded-Host"]; !ok {
		rawRequest.Header.Set("X-Forwarded-Host", this.req.host)
	}

	if _, ok := rawRequest.Header["X-Forwarded-Proto"]; !ok {
		rawRequest.Header.Set("X-Forwarded-Proto", this.req.rawScheme)
	}
}
