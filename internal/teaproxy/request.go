package teaproxy

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/tealogs"
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/build/internal/teawaf"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// 文本mime-type列表
var textMimeMap = map[string]bool{
	"application/atom+xml":                true,
	"application/javascript":              true,
	"application/x-javascript":            true,
	"application/json":                    true,
	"application/rss+xml":                 true,
	"application/x-web-app-manifest+json": true,
	"application/xhtml+xml":               true,
	"application/xml":                     true,
	"image/svg+xml":                       true,
	"text/css":                            true,
	"text/plain":                          true,
	"text/javascript":                     true,
	"text/xml":                            true,
	"text/html":                           true,
	"text/xhtml":                          true,
	"text/sgml":                           true,
}

// byte pool
var bytePool256b = teautils.NewBytePool(20480, 256)
var bytePool1k = teautils.NewBytePool(20480, 1024)
var bytePool32k = teautils.NewBytePool(20480, 32*1024)
var bytePool128k = teautils.NewBytePool(20480, 128*1024)

// 环境变量
var HOSTNAME, _ = os.Hostname()

// 请求定义
// HTTP HEADER RFC: https://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html
type Request struct {
	isNew  bool
	raw    *http.Request
	server *teaconfigs.ServerConfig

	attrs map[string]string // 附加参数

	scheme                 string
	rawScheme              string // 原始的scheme
	uri                    string
	rawURI                 string // 跳转之前的uri
	host                   string
	method                 string
	serverName             string // @TODO
	serverAddr             string
	charset                string
	requestHeaders         []*shared.HeaderConfig // 自定义请求Header
	responseHeaders        []*shared.HeaderConfig // 自定义响应Header
	uppercaseIgnoreHeaders []string               // 忽略的响应Header
	varMapping             map[string]string      // 自定义变量

	root      string   // 资源根目录
	urlPrefix string   // URL前缀
	index     []string // 目录下默认访问的文件

	backend     *teaconfigs.BackendConfig
	backendCall *shared.RequestCall

	fastcgi      *teaconfigs.FastcgiConfig
	proxy        *teaconfigs.ServerConfig
	location     *teaconfigs.LocationConfig
	accessPolicy *shared.AccessPolicy

	cachePolicy  *shared.CachePolicy
	cacheEnabled bool

	waf *teawaf.WAF

	pages    []*teaconfigs.PageConfig
	shutdown *teaconfigs.ShutdownConfig

	rewriteId           string // 匹配的rewrite id
	rewriteReplace      string // 经过rewrite之后的URL
	rewriteRedirectMode string // 跳转方式
	rewriteIsExternal   bool   // 是否为外部URL
	rewriteIsPermanent  bool   // 是否Permanent跳转
	rewriteProxyHost    string // 重写主机名

	redirectToHttps bool

	websocket *teaconfigs.WebsocketConfig

	tunnel *teaconfigs.TunnelConfig

	// 执行请求
	filePath string

	responseWriter   *ResponseWriter
	responseCallback func(http.ResponseWriter)

	requestFromTime time.Time // 请求开始时间
	requestCost     float64   // 请求耗时
	requestMaxSize  int64

	isWatching  bool     // 是否在监控
	requestData []byte   // 导出的request，在监控请求的时候有用
	errors      []string // 错误信息

	enableStat bool
	accessLog  *teaconfigs.AccessLogConfig

	gzip  *teaconfigs.GzipConfig
	debug bool

	locationContext *teaconfigs.LocationConfig // 当前变量的上下文 *Location ...

	hasForwardHeader bool
	isDenied         bool
}

// 获取新的请求
func NewRequest(rawRequest *http.Request) *Request {
	req := &Request{
		varMapping:      map[string]string{},
		raw:             rawRequest,
		rawURI:          rawRequest.URL.RequestURI(),
		requestFromTime: time.Now(),
		enableStat:      true,
		attrs:           map[string]string{},
	}

	req.backendCall = shared.NewRequestCall()
	req.backendCall.Reset()
	req.backendCall.Request = rawRequest
	req.backendCall.Formatter = req.Format

	_, req.hasForwardHeader = rawRequest.Header["X-Forwarded-For"]

	return req
}

// 初始化
func (this *Request) init(rawRequest *http.Request) {
	this.varMapping = map[string]string{}
	this.raw = rawRequest
	this.rawURI = rawRequest.URL.RequestURI()
	this.requestFromTime = time.Now()
	this.enableStat = true
	this.attrs = map[string]string{}

	if this.backendCall != nil {
		this.backendCall.Reset()
	} else {
		this.backendCall = shared.NewRequestCall()
	}
	this.backendCall.Request = rawRequest
	this.backendCall.Formatter = this.Format

	_, this.hasForwardHeader = rawRequest.Header["X-Forwarded-For"]
}

// 重置
func (this *Request) reset(rawRequest *http.Request) {
	this.server = nil

	this.requestHeaders = nil
	this.responseHeaders = nil
	this.uppercaseIgnoreHeaders = nil

	this.urlPrefix = ""

	this.backend = nil
	this.fastcgi = nil
	this.proxy = nil
	this.location = nil
	this.accessPolicy = nil
	this.cachePolicy = nil
	this.cacheEnabled = false
	this.waf = nil
	this.pages = nil
	this.shutdown = nil

	this.rewriteId = ""
	this.rewriteReplace = ""
	this.rewriteRedirectMode = ""
	this.rewriteIsExternal = false
	this.rewriteIsPermanent = false
	this.rewriteProxyHost = ""

	this.redirectToHttps = false

	this.websocket = nil
	this.tunnel = nil

	this.filePath = ""

	this.responseCallback = nil

	this.requestCost = 0
	this.requestMaxSize = 0

	this.isWatching = false
	this.requestData = nil
	this.errors = nil

	this.accessLog = nil

	this.locationContext = nil
	this.isDenied = false

	this.gzip = nil
	this.debug = false

	this.init(rawRequest)
}

func (this *Request) configure(server *teaconfigs.ServerConfig, redirects int, breakRewrite bool) error {
	isChanged := this.server != server
	this.server = server

	if redirects > 8 {
		return errors.New("too many redirects")
	}
	redirects++

	rawPath := ""
	rawQuery := ""
	qIndex := strings.Index(this.uri, "?") // question mark index
	if qIndex > -1 {
		rawPath = this.uri[:qIndex]
		rawQuery = this.uri[qIndex+1:]
	} else {
		rawPath = this.uri
	}

	// 是否切换了内部的代理服务
	if isChanged {
		// root
		this.root = server.Root
		if len(this.root) > 0 {
			this.root = this.Format(this.root)
		}

		// 字符集
		if len(server.Charset) > 0 {
			this.charset = this.Format(server.Charset)
		}

		// Header
		if server.HasRequestHeaders() {
			this.requestHeaders = append(this.requestHeaders, server.RequestHeaders...)
		}

		if server.HasResponseHeaders() {
			this.responseHeaders = append(this.responseHeaders, server.Headers...)
		}

		if server.HasIgnoreHeaders() {
			this.uppercaseIgnoreHeaders = append(this.uppercaseIgnoreHeaders, server.UppercaseIgnoreHeaders()...)
		}

		// cache
		if server.CacheOn {
			cachePolicy := server.CachePolicyObject()
			if cachePolicy != nil && cachePolicy.On {
				this.cachePolicy = cachePolicy
			}
		} else {
			this.cachePolicy = nil
		}

		// waf
		if server.WAFOn {
			waf := server.WAF()
			if waf != nil && waf.On && waf.MatchConds(this.Format) {
				this.waf = waf
			}
		} else {
			this.waf = nil
		}

		// tunnel
		if server.Tunnel != nil && server.Tunnel.On {
			this.tunnel = server.Tunnel
			return nil
		} else {
			this.tunnel = nil
		}

		// other
		if server.MaxBodyBytes() > 0 {
			this.requestMaxSize = server.MaxBodyBytes()
		}
		if len(server.AccessLog) > 0 {
			this.accessLog = server.AccessLog[0]
		}
		if server.DisableStat {
			this.enableStat = false
		}
		if len(server.Pages) > 0 {
			this.pages = append(this.pages, server.Pages...)
		}
		if server.Shutdown != nil && server.Shutdown.On {
			this.shutdown = server.Shutdown
		}
		if server.Gzip != nil {
			this.gzip = server.Gzip
		}

		if server.RedirectToHttps && this.rawScheme == "http" {
			this.redirectToHttps = true
			return nil
		}
	}

	// 如果是正向代理，则直接返回
	if server.ForwardHTTP != nil {
		return nil
	}

	if !breakRewrite {
		// location的相关配置
		var locationConfigured = false
		for _, location := range server.Locations {
			if !location.On {
				continue
			}
			this.locationContext = location
			if locationMatches, ok := location.Match(rawPath, this.Format); ok {
				if location.IsDenied(this.Format) {
					this.isDenied = true
					return nil
				}

				this.addVarMapping(locationMatches)

				if len(location.Root) > 0 {
					this.root = this.Format(location.Root)
					this.urlPrefix = location.URLPrefix
					locationConfigured = true
				}
				if len(location.Charset) > 0 {
					this.charset = this.Format(location.Charset)
				}
				if len(location.Index) > 0 {
					this.index = this.formatAll(location.Index)
				}
				if location.MaxBodyBytes() > 0 {
					this.requestMaxSize = location.MaxBodyBytes()
				}
				if len(location.AccessLog) > 0 {
					this.accessLog = location.AccessLog[0]
				}
				this.enableStat = !location.DisableStat
				if location.Gzip != nil {
					this.gzip = location.Gzip
				}
				if len(location.Pages) > 0 {
					this.pages = append(this.pages, location.Pages...)
				}
				if location.Shutdown != nil && location.Shutdown.On {
					this.shutdown = location.Shutdown
				}
				if location.RedirectToHttps && this.rawScheme == "http" {
					this.redirectToHttps = true
					this.locationContext = nil
					return nil
				}

				if location.CacheOn {
					cachePolicy := location.CachePolicyObject()
					if cachePolicy != nil && cachePolicy.On {
						this.cachePolicy = cachePolicy
					}
				} else {
					this.cachePolicy = nil
				}

				if location.WAFOn {
					waf := location.WAF()
					if waf != nil && waf.On && waf.MatchConds(this.Format) {
						this.waf = waf
					}
				} else {
					this.waf = nil
				}

				if location.HasRequestHeaders() {
					this.requestHeaders = append(this.requestHeaders, location.RequestHeaders...)
				}

				if location.HasResponseHeaders() {
					this.responseHeaders = append(this.responseHeaders, location.Headers...)
				}

				if location.HasIgnoreHeaders() {
					this.uppercaseIgnoreHeaders = append(this.uppercaseIgnoreHeaders, location.UppercaseIgnoreHeaders()...)
				}

				if location.AccessPolicy != nil {
					this.accessPolicy = location.AccessPolicy
				}

				this.location = location

				// rewrite相关配置
				if len(location.Rewrite) > 0 {
					for _, rule := range location.Rewrite {
						if !rule.On {
							continue
						}

						if replace, varMapping, ok := rule.Match(rawPath, this.Format); ok {
							this.addVarMapping(varMapping)
							this.rewriteId = rule.Id
							this.rewriteIsPermanent = rule.IsPermanent

							if rule.HasResponseHeaders() {
								this.responseHeaders = append(this.responseHeaders, rule.Headers...)
							}

							if rule.HasIgnoreHeaders() {
								this.uppercaseIgnoreHeaders = append(this.uppercaseIgnoreHeaders, rule.UppercaseIgnoreHeaders()...)
							}

							// 外部URL
							if rule.IsExternalURL(replace) {
								this.rewriteReplace = replace
								this.rewriteIsExternal = true
								this.rewriteRedirectMode = rule.RedirectMode()
								this.rewriteProxyHost = rule.ProxyHost
								this.locationContext = nil
								return nil
							}

							// 内部URL
							if rule.RedirectMode() == teaconfigs.RewriteFlagRedirect {
								this.rewriteReplace = replace
								this.rewriteIsExternal = false
								this.rewriteRedirectMode = teaconfigs.RewriteFlagRedirect
								this.locationContext = nil
								return nil
							}

							newURI, err := url.ParseRequestURI(replace)
							if err != nil {
								this.uri = replace
								this.locationContext = nil
								return nil
							}
							if len(newURI.RawQuery) > 0 {
								this.uri = newURI.Path + "?" + newURI.RawQuery
								if len(rawQuery) > 0 {
									this.uri += "&" + rawQuery
								}
							} else {
								this.uri = newURI.Path
								if len(rawQuery) > 0 {
									this.uri += "?" + rawQuery
								}
							}

							switch rule.TargetType() {
							case teaconfigs.RewriteTargetURL:
								this.locationContext = nil
								return this.configure(server, redirects, rule.IsBreak)
							case teaconfigs.RewriteTargetProxy:
								proxyId := rule.TargetProxy()
								server := SharedManager.FindServer(proxyId)
								if server == nil {
									this.locationContext = nil
									return errors.New("server with '" + proxyId + "' not found")
								}
								if !server.On {
									this.locationContext = nil
									return errors.New("server with '" + proxyId + "' not available now")
								}
								this.locationContext = nil
								return this.configure(server, redirects, rule.IsBreak)
							}
							this.locationContext = nil
							return nil
						}
					}
				}

				// fastcgi
				fastcgi := location.NextFastcgi()
				if fastcgi != nil {
					this.fastcgi = fastcgi
					this.backend = nil // 防止冲突
					locationConfigured = true

					if fastcgi.HasResponseHeaders() {
						this.responseHeaders = append(this.responseHeaders, fastcgi.Headers...)
					}

					if fastcgi.HasIgnoreHeaders() {
						this.uppercaseIgnoreHeaders = append(this.uppercaseIgnoreHeaders, fastcgi.UppercaseIgnoreHeaders()...)
					}

					// break
					if location.IsBreak {
						break
					}

					continue
				}

				// proxy
				if len(location.Proxy) > 0 {
					server := SharedManager.FindServer(location.Proxy)
					if server == nil {
						this.locationContext = nil
						return errors.New("server with '" + location.Proxy + "' not found")
					}
					if !server.On {
						this.locationContext = nil
						return errors.New("server with '" + location.Proxy + "' not available now")
					}
					this.locationContext = nil
					return this.configure(server, redirects, breakRewrite)
				}

				// backends
				if len(location.Backends) > 0 {
					backend := location.NextBackend(this.backendCall)
					if backend == nil {
						this.locationContext = nil
						return errors.New("no backends available")
					}
					if len(this.backendCall.ResponseCallbacks) > 0 {
						this.responseCallback = this.backendCall.CallResponseCallbacks
					}
					this.backend = backend
					locationConfigured = true

					if backend.HasRequestHeaders() {
						this.requestHeaders = append(this.requestHeaders, backend.RequestHeaders...)
					}

					if backend.HasResponseHeaders() {
						this.responseHeaders = append(this.responseHeaders, backend.Headers...)
					}

					if backend.HasIgnoreHeaders() {
						this.uppercaseIgnoreHeaders = append(this.uppercaseIgnoreHeaders, backend.UppercaseIgnoreHeaders()...)
					}

					// break
					if location.IsBreak {
						break
					}

					continue
				}

				// websocket
				if location.Websocket != nil && location.Websocket.On {
					this.backend = location.Websocket.NextBackend(this.backendCall)
					this.websocket = location.Websocket
					this.locationContext = nil
					return nil
				}

				// break
				if location.IsBreak {
					break
				}
			}
			this.locationContext = nil
		}

		// 如果经过location找到了相关配置，就终止
		if locationConfigured {
			return nil
		}
	}

	// server的相关配置
	if !breakRewrite && len(server.Rewrite) > 0 {
		for _, rule := range server.Rewrite {
			if !rule.On {
				continue
			}
			if replace, varMapping, ok := rule.Match(rawPath, func(source string) string {
				return this.Format(source)
			}); ok {
				this.addVarMapping(varMapping)
				this.rewriteId = rule.Id
				this.rewriteIsPermanent = rule.IsPermanent

				if rule.HasRequestHeaders() {
					this.requestHeaders = append(this.requestHeaders, rule.RequestHeaders...)
				}

				if rule.HasResponseHeaders() {
					this.responseHeaders = append(this.responseHeaders, rule.Headers...)
				}

				if rule.HasIgnoreHeaders() {
					this.uppercaseIgnoreHeaders = append(this.uppercaseIgnoreHeaders, rule.UppercaseIgnoreHeaders()...)
				}

				// 外部URL
				if rule.IsExternalURL(replace) {
					this.rewriteReplace = replace
					this.rewriteIsExternal = true
					this.rewriteRedirectMode = rule.RedirectMode()
					this.rewriteProxyHost = rule.ProxyHost
					return nil
				}

				// 内部URL
				if rule.RedirectMode() == teaconfigs.RewriteFlagRedirect {
					this.rewriteReplace = replace
					this.rewriteIsExternal = false
					this.rewriteRedirectMode = teaconfigs.RewriteFlagRedirect
					return nil
				}

				newURI, err := url.ParseRequestURI(replace)
				if err != nil {
					this.uri = replace
					return nil
				}
				if len(newURI.RawQuery) > 0 {
					this.uri = newURI.Path + "?" + newURI.RawQuery
					if len(rawQuery) > 0 {
						this.uri += "&" + rawQuery
					}
				} else {
					if len(rawQuery) > 0 {
						this.uri = newURI.Path + "?" + rawQuery
					}
				}

				switch rule.TargetType() {
				case teaconfigs.RewriteTargetURL:
					return this.configure(server, redirects, rule.IsBreak)
				case teaconfigs.RewriteTargetProxy:
					proxyId := rule.TargetProxy()
					server := SharedManager.FindServer(proxyId)
					if server == nil {
						return errors.New("server with '" + proxyId + "' not found")
					}
					if !server.On {
						return errors.New("server with '" + proxyId + "' not available now")
					}
					return this.configure(server, redirects, rule.IsBreak)
				}
				return nil
			}
		}
	}

	// fastcgi
	fastcgi := server.NextFastcgi()
	if fastcgi != nil {
		this.fastcgi = fastcgi
		this.backend = nil // 防止冲突

		if fastcgi.HasRequestHeaders() {
			this.requestHeaders = append(this.requestHeaders, fastcgi.RequestHeaders...)
		}

		if fastcgi.HasResponseHeaders() {
			this.responseHeaders = append(this.responseHeaders, fastcgi.Headers...)
		}

		if fastcgi.HasIgnoreHeaders() {
			this.uppercaseIgnoreHeaders = append(this.uppercaseIgnoreHeaders, fastcgi.UppercaseIgnoreHeaders()...)
		}

		return nil
	}

	// proxy
	if len(server.Proxy) > 0 {
		server := SharedManager.FindServer(server.Proxy)
		if server == nil {
			return errors.New("server with '" + server.Proxy + "' not found")
		}
		if !server.On {
			return errors.New("server with '" + server.Proxy + "' not available now")
		}
		return this.configure(server, redirects, breakRewrite)
	}

	// 转发到后端
	backend := server.NextBackend(this.backendCall)
	if backend == nil {
		if len(this.root) == 0 {
			return errors.New("no backends available")
		}
		this.backend = nil
		return nil
	}
	if len(this.backendCall.ResponseCallbacks) > 0 {
		this.responseCallback = this.backendCall.CallResponseCallbacks
	}
	this.backend = backend

	if backend != nil {
		if backend.HasRequestHeaders() {
			this.requestHeaders = append(this.requestHeaders, backend.RequestHeaders...)
		}

		if backend.HasResponseHeaders() {
			this.responseHeaders = append(this.responseHeaders, backend.Headers...)
		}

		if backend.HasIgnoreHeaders() {
			this.uppercaseIgnoreHeaders = append(this.uppercaseIgnoreHeaders, backend.UppercaseIgnoreHeaders()...)
		}
	}

	return nil
}

func (this *Request) call(writer *ResponseWriter) error {
	if this.requestMaxSize > 0 {
		this.raw.Body = http.MaxBytesReader(writer, this.raw.Body, this.requestMaxSize)
	}

	this.responseWriter = writer

	// WAF
	if this.waf != nil {
		if this.callWAFRequest(writer) {
			this.callEnd(writer)
			return nil
		}
	}

	// 跳转到https
	if this.redirectToHttps {
		this.callRedirectToHttps(writer)
		return nil
	}

	// hook
	b := CallRequestBeforeHook(this, writer)
	if !b {
		this.callEnd(writer)
		return nil
	}

	// gzip压缩
	hasGzip := false
	if this.gzip != nil && this.gzip.Level > 0 && this.acceptGzipEncoding() {
		hasGzip = true
		writer.Gzip(this.gzip)
	}

	err := this.callBegin(writer)

	// 在结束之前关闭gzip以便能够获取完整的body
	if hasGzip {
		writer.Close()
	}

	this.callEnd(writer)
	return err
}

func (this *Request) callBegin(writer *ResponseWriter) error {
	// watch
	if this.isWatching {
		// 判断如果Content-Length过长，则截断
		reqData, err := httputil.DumpRequest(this.raw, true)
		if err == nil {
			if len(reqData) > 100240 {
				reqData = reqData[:100240]
			}
			this.requestData = reqData
		}

		writer.SetBodyCopying(true)
	} else {
		max := 512 * 1024 // 512K
		if this.accessLog != nil && lists.ContainsInt(this.accessLog.Fields, accesslogs.AccessLogFieldRequestBody) {
			body, err := ioutil.ReadAll(this.raw.Body)
			if err == nil {
				if len(body) > max {
					this.requestData = body[:max]
				} else {
					this.requestData = body
				}
			}
			this.raw.Body = ioutil.NopCloser(bytes.NewReader(body))
		}
		if this.accessLog != nil && lists.ContainsInt(this.accessLog.Fields, accesslogs.AccessLogFieldResponseBody) {
			writer.SetBodyCopying(true)
		}
	}

	// access policy
	if this.accessPolicy != nil {
		if !this.accessPolicy.AllowAccess(this.requestRemoteAddr()) {
			if !this.callPage(writer, http.StatusForbidden) {
				writer.WriteHeader(http.StatusForbidden)
				_, _ = writer.Write([]byte("Request Forbidden"))
			}
			return nil
		}

		reason, allowed := this.accessPolicy.AllowTraffic()
		if !allowed {
			if !this.callPage(writer, http.StatusTooManyRequests) {
				writer.WriteHeader(http.StatusTooManyRequests)
				_, _ = writer.Write([]byte("[" + reason + "]Request Quota Exceeded"))
			}
			return nil
		}
	}

	// 禁止访问页面
	if this.isDenied {
		if !this.callPage(writer, http.StatusForbidden) {
			writer.WriteHeader(http.StatusForbidden)
			_, _ = writer.Write([]byte("Request Forbidden"))
		}
		return nil
	}

	// 临时关闭页面
	if this.shutdown != nil {
		return this.callShutdown(writer)
	}

	if this.tunnel != nil {
		return this.callTunnel(writer)
	}
	if this.websocket != nil {
		return this.callWebsocket(writer)
	}
	if this.backend != nil {
		return this.callBackend(writer)
	}
	if this.proxy != nil {
		return this.callProxy(writer)
	}
	if this.fastcgi != nil {
		return this.callFastcgi(writer)
	}
	if len(this.rewriteId) > 0 && (this.rewriteIsExternal || this.rewriteRedirectMode == teaconfigs.RewriteFlagRedirect) {
		return this.callRewrite(writer)
	}
	if len(this.root) > 0 {
		return this.callRoot(writer)
	}

	return errors.New("unable to handle the request")
}

func (this *Request) callEnd(writer *ResponseWriter) {
	// log
	this.log()

	// call hook
	CallRequestAfterHook(this, writer)
}

func (this *Request) notFoundError(writer *ResponseWriter) {
	if this.callPage(writer, http.StatusNotFound) {
		return
	}

	msg := "404 page not found: '" + this.requestURI() + "'"

	writer.WriteHeader(http.StatusNotFound)
	_, _ = writer.Write([]byte(msg))
}

func (this *Request) serverError(writer *ResponseWriter) {
	if this.callPage(writer, http.StatusInternalServerError) {
		return
	}

	statusCode := http.StatusInternalServerError

	// 忽略的Header
	ignoreHeaders := this.convertIgnoreHeaders()
	hasIgnoreHeaders := ignoreHeaders.Len() > 0

	// 自定义Header
	for _, header := range this.responseHeaders {
		if header.Match(statusCode) {
			if hasIgnoreHeaders && ignoreHeaders.Has(strings.ToUpper(header.Name)) {
				continue
			}
			writer.Header().Set(header.Name, this.Format(header.Value))
		}
	}

	writer.WriteHeader(statusCode)
	_, _ = writer.Write([]byte(http.StatusText(statusCode)))
}

func (this *Request) requestRemoteAddr() string {
	// X-Forwarded-For
	forwardedFor := this.raw.Header.Get("X-Forwarded-For")
	if len(forwardedFor) > 0 {
		commaIndex := strings.Index(forwardedFor, ",")
		if commaIndex > 0 {
			return forwardedFor[:commaIndex]
		}
		return forwardedFor
	}

	// Real-IP
	{
		realIP, ok := this.raw.Header["X-Real-IP"]
		if ok && len(realIP) > 0 {
			return realIP[0]
		}
	}

	// Real-Ip
	{
		realIP, ok := this.raw.Header["X-Real-Ip"]
		if ok && len(realIP) > 0 {
			return realIP[0]
		}
	}

	// Remote-Addr
	remoteAddr := this.raw.RemoteAddr
	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil {
		return host
	} else {
		return remoteAddr
	}
}

func (this *Request) requestRemotePort() int {
	_, port, err := net.SplitHostPort(this.raw.RemoteAddr)
	if err == nil {
		return types.Int(port)
	}
	return 0
}

func (this *Request) requestRemoteUser() string {
	username, _, ok := this.raw.BasicAuth()
	if !ok {
		return ""
	}
	return username
}

func (this *Request) requestURI() string {
	return this.rawURI
}

func (this *Request) requestPath() string {
	uri, err := url.ParseRequestURI(this.requestURI())
	if err != nil {
		return ""
	}
	return uri.Path
}

func (this *Request) requestLength() int64 {
	return this.raw.ContentLength
}

func (this *Request) requestMethod() string {
	return this.method
}

func (this *Request) requestFilename() string {
	return this.filePath
}

func (this *Request) requestProto() string {
	return this.raw.Proto
}

func (this *Request) requestReferer() string {
	return this.raw.Referer()
}

func (this *Request) requestUserAgent() string {
	return this.raw.UserAgent()
}

func (this *Request) requestContentType() string {
	return this.raw.Header.Get("Content-Type")
}

func (this *Request) requestString() string {
	return this.method + " " + this.requestURI() + " " + this.requestProto()
}

func (this *Request) requestCookiesString() string {
	var cookies = []string{}
	for _, cookie := range this.raw.Cookies() {
		cookies = append(cookies, url.QueryEscape(cookie.Name)+"="+url.QueryEscape(cookie.Value))
	}
	return strings.Join(cookies, "&")
}

func (this *Request) requestCookie(name string) string {
	cookie, err := this.raw.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func (this *Request) requestQueryString() string {
	uri, err := url.ParseRequestURI(this.uri)
	if err != nil {
		return ""
	}
	return uri.RawQuery
}

func (this *Request) requestQueryParam(name string) string {
	uri, err := url.ParseRequestURI(this.uri)
	if err != nil {
		return ""
	}

	v, found := uri.Query()[name]
	if !found {
		return ""
	}
	return strings.Join(v, "&")
}

func (this *Request) requestServerPort() int {
	_, port, err := net.SplitHostPort(this.serverAddr)
	if err == nil {
		return types.Int(port)
	}
	return 0
}

func (this *Request) requestHeadersString() string {
	var headers = []string{}
	for k, v := range this.raw.Header {
		for _, subV := range v {
			headers = append(headers, k+": "+subV)
		}
	}
	return strings.Join(headers, ";")
}

func (this *Request) requestHeader(key string) string {
	v, found := this.raw.Header[key]
	if !found {
		return ""
	}
	return strings.Join(v, ";")
}

func (this *Request) acceptGzipEncoding() bool {
	encodingList := this.raw.Header.Get("Accept-Encoding")
	if len(encodingList) == 0 {
		return false
	}
	encodings := strings.Split(encodingList, ",")
	for _, encoding := range encodings {
		if encoding == "gzip" {
			return true
		}
	}
	return false
}

func (this *Request) CachePolicy() *shared.CachePolicy {
	return this.cachePolicy
}

func (this *Request) SetCachePolicy(config *shared.CachePolicy) {
	this.cachePolicy = config
}

func (this *Request) SetCacheEnabled() {
	this.cacheEnabled = true
}

// 判断缓存策略是否有效
func (this *Request) IsCacheEnabled() bool {
	return this.cacheEnabled
}

// 设置监控状态
func (this *Request) SetIsWatching(isWatching bool) {
	this.isWatching = isWatching
}

// 判断是否在监控
func (this *Request) IsWatching() bool {
	return this.isWatching
}

// 设置URI
func (this *Request) SetURI(uri string) {
	this.uri = uri
}

// 设置Host
func (this *Request) SetHost(host string) {
	this.host = host
}

// 设置原始的scheme
func (this *Request) SetRawScheme(scheme string) {
	this.rawScheme = scheme
}

// 获取原始的请求
func (this *Request) Raw() *http.Request {
	return this.raw
}

// 输出自定义Response Header
func (this *Request) WriteResponseHeaders(writer *ResponseWriter, statusCode int) {
	ignoreHeaders := this.convertIgnoreHeaders()
	hasIgnoreHeaders := ignoreHeaders.Len() > 0

	responseHeader := writer.Header()

	for _, header := range this.responseHeaders {
		if !header.On {
			continue
		}
		if header.Match(statusCode) {
			if hasIgnoreHeaders && ignoreHeaders.Has(strings.ToUpper(header.Name)) {
				continue
			}
			if header.HasVariables() {
				responseHeader.Set(header.Name, this.Format(header.Value))
			} else {
				responseHeader.Set(header.Name, header.Value)
			}
		}
	}

	// hsts
	if this.rawScheme == "https" &&
		this.server.SSL != nil &&
		this.server.SSL.On &&
		this.server.SSL.HSTS != nil &&
		this.server.SSL.HSTS.On &&
		this.server.SSL.HSTS.Match(this.host) {
		responseHeader.Set(this.server.SSL.HSTS.HeaderKey(), this.server.SSL.HSTS.HeaderValue())
	}
}

// 利用请求参数格式化字符串
func (this *Request) Format(source string) string {
	if len(source) == 0 {
		return ""
	}

	var hasVarMapping = len(this.varMapping) > 0

	return teautils.ParseVariables(source, func(varName string) string {
		// 自定义变量
		if hasVarMapping {
			value, found := this.varMapping[varName]
			if found {
				return value
			}
		}

		// 请求变量
		switch varName {
		case "teaVersion":
			return teaconst.TeaVersion
		case "remoteAddr":
			return this.requestRemoteAddr()
		case "rawRemoteAddr":
			addr := this.raw.RemoteAddr
			host, _, err := net.SplitHostPort(addr)
			if err == nil {
				addr = host
			}
			return addr
		case "remotePort":
			return strconv.Itoa(this.requestRemotePort())
		case "remoteUser":
			return this.requestRemoteUser()
		case "requestURI", "requestUri":
			return this.requestURI()
		case "requestPath":
			return this.requestPath()
		case "requestLength":
			return fmt.Sprintf("%d", this.requestLength())
		case "requestTime":
			return fmt.Sprintf("%.6f", this.requestCost)
		case "requestMethod":
			return this.requestMethod()
		case "requestFilename":
			filename := this.requestFilename()
			if len(filename) > 0 {
				return filename
			}

			if this.locationContext != nil && len(this.locationContext.Root) > 0 {
				return filepath.Clean(this.locationContext.Root + this.requestPath())
			}

			if len(this.root) > 0 {
				return filepath.Clean(this.root + this.requestPath())
			}

			return ""
		case "scheme":
			return this.rawScheme
		case "serverProtocol", "proto":
			return this.requestProto()
		case "bytesSent":
			return fmt.Sprintf("%d", this.responseWriter.SentBodyBytes()) // TODO 加上Header长度
		case "bodyBytesSent":
			return fmt.Sprintf("%d", this.responseWriter.SentBodyBytes())
		case "status":
			return strconv.Itoa(this.responseWriter.StatusCode())
		case "statusMessage":
			return http.StatusText(this.responseWriter.StatusCode())
		case "timeISO8601":
			return this.requestFromTime.Format("2006-01-02T15:04:05.000Z07:00")
		case "timeLocal":
			return this.requestFromTime.Format("2/Jan/2006:15:04:05 -0700")
		case "msec":
			return fmt.Sprintf("%.6f", float64(this.requestFromTime.Unix())+float64(this.requestFromTime.Nanosecond())/1000000000)
		case "timestamp":
			return fmt.Sprintf("%d", this.requestFromTime.Unix())
		case "host":
			return this.host
		case "referer":
			return this.requestReferer()
		case "userAgent":
			return this.requestUserAgent()
		case "contentType":
			return this.requestContentType()
		case "request":
			return this.requestString()
		case "cookies":
			return this.requestCookiesString()
		case "args", "queryString":
			return this.requestQueryString()
		case "headers":
			return this.requestHeadersString()
		case "serverName":
			return this.serverName
		case "serverPort":
			return strconv.Itoa(this.requestServerPort())
		case "hostname":
			return HOSTNAME
		case "documentRoot":
			if this.locationContext != nil && len(this.locationContext.Root) > 0 {
				return this.locationContext.Root
			}
			return this.root
		}

		dotIndex := strings.Index(varName, ".")
		if dotIndex < 0 {
			return "${" + varName + "}"
		}
		prefix := varName[:dotIndex]
		suffix := varName[dotIndex+1:]

		// cookie.
		if prefix == "cookie" {
			return this.requestCookie(suffix)
		}

		// arg.
		if prefix == "arg" {
			return this.requestQueryParam(suffix)
		}

		// header.
		if prefix == "header" || prefix == "http" {
			return this.requestHeader(suffix)
		}

		// backend.
		if prefix == "backend" {
			if this.backend != nil {
				switch suffix {
				case "address":
					return this.backend.Address
				case "host":
					index := strings.Index(this.backend.Address, ":")
					if index > -1 {
						return this.backend.Address[:index]
					} else {
						return ""
					}
				case "id":
					return this.backend.Id
				case "scheme":
					return this.backend.Scheme
				case "code":
					return this.backend.Code
				}
			}
			return ""
		}

		// node
		if prefix == "node" {
			node := teaconfigs.SharedNodeConfig()
			if node != nil {
				switch suffix {
				case "id":
					return node.Id
				case "name":
					return node.Name
				case "role":
					return node.Role
				}
			}
		}

		// host
		if prefix == "host" {
			pieces := strings.Split(this.host, ".")
			switch suffix {
			case "first":
				if len(pieces) > 0 {
					return pieces[0]
				}
				return ""
			case "last":
				if len(pieces) > 0 {
					return pieces[len(pieces)-1]
				}
				return ""
			case "0":
				if len(pieces) > 0 {
					return pieces[0]
				}
				return ""
			case "1":
				if len(pieces) > 1 {
					return pieces[1]
				}
				return ""
			case "2":
				if len(pieces) > 2 {
					return pieces[2]
				}
				return ""
			case "3":
				if len(pieces) > 3 {
					return pieces[3]
				}
				return ""
			case "4":
				if len(pieces) > 4 {
					return pieces[4]
				}
				return ""
			case "-1":
				if len(pieces) > 0 {
					return pieces[len(pieces)-1]
				}
				return ""
			case "-2":
				if len(pieces) > 1 {
					return pieces[len(pieces)-2]
				}
				return ""
			case "-3":
				if len(pieces) > 2 {
					return pieces[len(pieces)-3]
				}
				return ""
			case "-4":
				if len(pieces) > 3 {
					return pieces[len(pieces)-4]
				}
				return ""
			case "-5":
				if len(pieces) > 4 {
					return pieces[len(pieces)-5]
				}
				return ""
			}
		}

		return "${" + varName + "}"
	})
}

// 设置属性
func (this *Request) SetAttr(key string, value string) {
	// 需要处理key中的点（.）符号，因为很多数据库不支持在key中含有点
	key = strings.Replace(key, ".", "_", -1)
	this.attrs[key] = value
}

// 格式化一组字符串
func (this *Request) formatAll(sources []string) []string {
	result := []string{}
	for _, s := range sources {
		result = append(result, this.Format(s))
	}
	return result
}

// 记录日志
func (this *Request) log() {
	// 计算请求时间
	this.requestCost = time.Since(this.requestFromTime).Seconds()

	if (this.accessLog == nil || !this.accessLog.On) && !this.enableStat {
		return
	}

	cookies := map[string]string{}
	for _, cookie := range this.raw.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}

	addr := this.raw.RemoteAddr
	host, _, err := net.SplitHostPort(addr)
	if err == nil {
		addr = host
	}

	accessLog := &accesslogs.AccessLog{
		TeaVersion:      teaconst.TeaVersion,
		RemoteAddr:      this.requestRemoteAddr(),
		RawRemoteAddr:   addr,
		RemotePort:      this.requestRemotePort(),
		RemoteUser:      this.requestRemoteUser(),
		RequestURI:      this.requestURI(),
		RequestPath:     this.requestPath(),
		RequestLength:   this.requestLength(),
		RequestTime:     this.requestCost,
		RequestMethod:   this.requestMethod(),
		RequestFilename: this.requestFilename(),
		Scheme:          this.rawScheme,
		Proto:           this.requestProto(),
		BytesSent:       this.responseWriter.SentBodyBytes(), // TODO 加上Header Size
		BodyBytesSent:   this.responseWriter.SentBodyBytes(),
		Status:          this.responseWriter.StatusCode(),
		StatusMessage:   "",
		TimeISO8601:     this.requestFromTime.Format("2006-01-02T15:04:05.000Z07:00"),
		TimeLocal:       this.requestFromTime.Format("2/Jan/2006:15:04:05 -0700"),
		Msec:            float64(this.requestFromTime.Unix()) + float64(this.requestFromTime.Nanosecond())/1000000000,
		Timestamp:       this.requestFromTime.Unix(),
		Host:            this.host,
		Referer:         this.requestReferer(),
		UserAgent:       this.requestUserAgent(),
		Request:         this.requestString(),
		ContentType:     this.requestContentType(),
		Cookie:          cookies,
		Args:            this.requestQueryString(),
		QueryString:     this.requestQueryString(),
		Header:          this.raw.Header,
		ServerName:      this.serverName,
		ServerPort:      this.requestServerPort(),
		ServerProtocol:  this.requestProto(),
		Errors:          this.errors,
		HasErrors:       len(this.errors) > 0,
		Extend:          &accesslogs.AccessLogExtend{},
		Attrs:           this.attrs,
		Hostname:        HOSTNAME,
	}

	// 日志和统计
	accessLog.SetShouldWrite(this.accessLog != nil && this.accessLog.On && this.accessLog.Match(this.responseWriter.statusCode))
	accessLog.SetShouldStat(this.enableStat)
	if this.accessLog != nil {
		accessLog.SetWritingFields(this.accessLog.Fields)
		accessLog.StorageOnly = this.accessLog.StorageOnly

		// 筛选策略
		policyIds := this.accessLog.StoragePolicies
		if len(policyIds) > 0 {
			resultPolicyIds := []string{}
			for _, policyId := range policyIds {
				policy := tealogs.FindPolicy(policyId)
				if policy == nil || !policy.On || !policy.MatchConds(this.Format) {
					continue
				}
				resultPolicyIds = append(resultPolicyIds, policyId)
			}
			policyIds = resultPolicyIds
		}
		accessLog.StoragePolicyIds = policyIds
	}

	if this.server != nil {
		accessLog.ServerId = this.server.Id
	}

	if this.backend != nil {
		accessLog.BackendAddress = this.backend.Address
		accessLog.BackendId = this.backend.Id
		accessLog.BackendScheme = this.backend.Scheme
		accessLog.BackendCode = this.backend.Code
	}

	if this.fastcgi != nil {
		accessLog.FastcgiAddress = this.fastcgi.Pass
		accessLog.FastcgiId = this.fastcgi.Id
	}

	accessLog.RewriteId = this.rewriteId

	if this.location != nil {
		accessLog.LocationId = this.location.Id
	}

	accessLog.SentHeader = this.responseWriter.Header()

	if len(this.requestData) > 0 {
		accessLog.RequestData = this.requestData
	}

	if this.responseWriter.BodyIsCopying() {
		accessLog.ResponseHeaderData = this.responseWriter.HeaderData()
		accessLog.ResponseBodyData = this.responseWriter.Body()
	}

	logger := tealogs.SharedLogger()
	if logger != nil {
		logger.Push(accessLog)
	}
}

func (this *Request) findIndexFile(dir string) string {
	if len(this.index) == 0 {
		return ""
	}
	for _, index := range this.index {
		if len(index) == 0 {
			continue
		}

		// 模糊查找
		if strings.Contains(index, "*") {
			indexFiles, err := filepath.Glob(dir + Tea.DS + index)
			if err != nil {
				logs.Error(err)
				this.addError(err)
				continue
			}
			if len(indexFiles) > 0 {
				return filepath.Base(indexFiles[0])
			}
			continue
		}

		// 精确查找
		filePath := dir + Tea.DS + index
		stat, err := os.Stat(filePath)
		if err != nil || !stat.Mode().IsRegular() {
			continue
		}
		return index
	}
	return ""
}

func (this *Request) convertIgnoreHeaders() maps.Map {
	if len(this.uppercaseIgnoreHeaders) == 0 {
		return nil
	}

	m := maps.Map{}
	for _, h := range this.uppercaseIgnoreHeaders {
		m[h] = true
	}
	return m
}

func (this *Request) addVarMapping(varMapping map[string]string) {
	for k, v := range varMapping {
		this.varMapping[k] = v
	}
}

// 添加自定义变量
func (this *Request) SetVarMapping(varName string, varValue string) {
	this.varMapping[varName] = varValue
}

func (this *Request) addError(err error) {
	if err == nil {
		return
	}
	this.errors = append(this.errors, err.Error())
}

// 设置代理相关头部信息
// 参考：https://tools.ietf.org/html/rfc7239
func (this *Request) setProxyHeaders(header http.Header) {
	delete(header, "Connection")

	remoteAddr := this.raw.RemoteAddr
	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil {
		remoteAddr = host
	}

	// x-real-ip
	{
		_, ok1 := header["X-Real-IP"]
		_, ok2 := header["X-Real-Ip"]
		if !ok1 && !ok2 {
			header["X-Real-IP"] = []string{remoteAddr}
		}
	}

	// X-Forwarded-For
	{
		forwardedFor, ok := header["X-Forwarded-For"]
		if ok {
			if this.hasForwardHeader {
				header["X-Forwarded-For"] = []string{strings.Join(forwardedFor, ", ") + ", " + remoteAddr}
			}
		} else {
			header["X-Forwarded-For"] = []string{remoteAddr}
		}
	}

	// Forwarded
	/**{
		forwarded, ok := header["Forwarded"]
		if ok {
			header["Forwarded"] = []string{strings.Join(forwarded, ", ") + ", by=" + this.serverAddr + "; for=" + remoteAddr + "; host=" + this.host + "; proto=" + this.rawScheme}
		} else {
			header["Forwarded"] = []string{"by=" + this.serverAddr + "; for=" + remoteAddr + "; host=" + this.host + "; proto=" + this.rawScheme}
		}
	}**/

	// others
	this.raw.Header.Set("X-Forwarded-By", this.serverAddr)

	if _, ok := header["X-Forwarded-Host"]; !ok {
		this.raw.Header.Set("X-Forwarded-Host", this.host)
	}

	if _, ok := header["X-Forwarded-Proto"]; !ok {
		this.raw.Header.Set("X-Forwarded-Proto", this.rawScheme)
	}
}

// 计算合适的buffer size
func (this *Request) bytePool(contentLength int64) *teautils.BytePool {
	if contentLength <= 0 {
		return bytePool1k
	}
	if contentLength < 1024 { // 1K
		return bytePool256b
	}
	if contentLength < 32768 { // 32K
		return bytePool1k
	}
	if contentLength < 1048576 { // 1M
		return bytePool32k
	}
	return bytePool128k
}
