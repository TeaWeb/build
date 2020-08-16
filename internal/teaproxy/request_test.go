package teaproxy

import (
	"bytes"
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/assert"
	"net/http"
	"net/url"
	"runtime"
	"testing"
	"time"
)

type testResponseWriter struct {
	a    *assert.Assertion
	data []byte
}

func testNewResponseWriter(a *assert.Assertion) *ResponseWriter {
	return NewResponseWriter(&testResponseWriter{
		a: a,
	})
}

func (this *testResponseWriter) Header() http.Header {
	return http.Header{}
}

func (this *testResponseWriter) Write(data []byte) (int, error) {
	this.data = append(this.data, data...)
	return len(data), nil
}

func (this *testResponseWriter) WriteHeader(statusCode int) {
}

func (this *testResponseWriter) Close() {
	this.a.Log(string(this.data))
}

func TestRequest_CallRoot(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()
	writer := testNewResponseWriter(a)

	req, _ := http.NewRequest(http.MethodGet, "http://teaos.cn/layout.css", nil)

	request := NewRequest(req)
	request.root = Tea.ViewsDir() + "/@default"
	request.uri = "/layout.css"
	err := request.call(writer)
	a.IsNil(err)

	a.Log("status:", writer.StatusCode())
	a.Log("requestTime:", request.requestCost)
	a.Log("bytes send:", writer.SentBodyBytes())
}

func TestRequest_CallBackend(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()
	writer := testNewResponseWriter(a)

	req, err := http.NewRequest("GET", "/index.php?__ACTION__=/@wx", nil)
	if err != nil {
		a.Fatal(err)
	}
	req.RemoteAddr = "127.0.0.1"
	request := NewRequest(req)
	request.scheme = "http"
	request.host = "wx.balefm.cn"
	request.backend = &teaconfigs.BackendConfig{
		Address: "127.0.0.1",
	}
	err = request.backend.Validate()
	if err != nil {
		t.Fatal(err)
	}
	err = request.call(writer)
	a.IsNil(err)

	a.Log("status:", writer.StatusCode())
	a.Log("requestTime:", request.requestCost)
	a.Log("bytes send:", writer.SentBodyBytes())
}

func TestRequest_CallProxy(t *testing.T) {
	if !teatesting.RequireHTTPServer() {
		return
	}

	a := assert.NewAssertion(t).Quiet()
	writer := testNewResponseWriter(a)

	req, err := http.NewRequest("GET", "/webhook", nil)
	if err != nil {
		a.Fatal(err)
	}
	req.RemoteAddr = "127.0.0.1"
	request := NewRequest(req)
	request.scheme = "http"
	request.host = "wx.balefm.cn"

	proxy := teaconfigs.NewServerConfig()
	proxy.AddBackend(&teaconfigs.BackendConfig{
		On:      true,
		Address: "127.0.0.1:9991",
	})
	//proxy.AddBackend(&teaconfigs.BackendConfig{
	//	On:      true,
	//	Address: "127.0.0.1:81",
	//})
	err = proxy.Validate()
	if err != nil {
		t.Fatal(err)
	}

	request.proxy = proxy

	err = request.call(writer)
	a.IsNil(err)

	a.Log("status:", writer.StatusCode())
	a.Log("requestTime:", request.requestCost)
	a.Log("bytes send:", writer.SentBodyBytes())
}

func TestRequest_CallFastcgi(t *testing.T) {
	if !teatesting.RequireFascgi() {
		return
	}

	a := assert.NewAssertion(t).Quiet()
	writer := testNewResponseWriter(a)

	req, err := http.NewRequest("GET", "/index.php?__ACTION__=/@wx/box/version", bytes.NewBuffer([]byte("hello=world")))
	//req, err := http.NewRequest("GET", "/index.php", bytes.NewBuffer([]byte("hello=world")))
	if err != nil {
		a.Fatal(err)
	}
	req.RemoteAddr = "127.0.0.1"
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	request := NewRequest(req)
	request.scheme = "http"
	request.host = "wx.balefm.cn"
	request.serverAddr = "127.0.0.1:80"

	request.fastcgi = &teaconfigs.FastcgiConfig{
		Params: map[string]string{
			"SCRIPT_FILENAME": "/Users/liuxiangchao/Documents/Projects/pp/apps/baleshop.ppk/index.php",
			//"DOCUMENT_ROOT":   "/Users/liuxiangchao/Documents/Projects/pp/apps/baleshop.ppk",
		},
		Pass: "127.0.0.1:9000",
	}
	err = request.fastcgi.Validate()
	if err != nil {
		t.Fatal(err)
	}
	err = request.call(writer)
	a.IsNil(err)

	a.Log("status:", writer.StatusCode())
	a.Log("requestTime:", request.requestCost)
	a.Log("bytes send:", writer.SentBodyBytes())
}

func TestRequest_CallFastcgiPerformance(t *testing.T) {
	if !teatesting.RequireFascgi() {
		return
	}

	a := assert.NewAssertion(t).Quiet()
	writer := testNewResponseWriter(a)

	req, err := http.NewRequest("GET", "/index.php?__ACTION__=/@wx/box/version", bytes.NewBuffer([]byte("hello=world")))
	//req, err := http.NewRequest("GET", "/index.php", bytes.NewBuffer([]byte("hello=world")))
	if err != nil {
		a.Fatal(err)
	}
	req.RemoteAddr = "127.0.0.1"
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	request := NewRequest(req)
	request.scheme = "http"
	request.host = "wx.balefm.cn"
	request.serverAddr = "127.0.0.1:80"

	request.fastcgi = &teaconfigs.FastcgiConfig{
		Params: map[string]string{
			"SCRIPT_FILENAME": "/Users/liuxiangchao/Documents/Projects/pp/apps/baleshop.ppk/index.php",
			//"DOCUMENT_ROOT":   "/Users/liuxiangchao/Documents/Projects/pp/apps/baleshop.ppk",
		},
		Pass: "127.0.0.1:9000",
	}
	err = request.fastcgi.Validate()
	if err != nil {
		t.Fatal(err)
	}
	err = request.call(writer)
	a.IsNil(err)

	a.Log("status:", writer.StatusCode())
	a.Log("requestTime:", request.requestCost)
	a.Log("bytes send:", writer.SentBodyBytes())
}

func TestRequest_Format(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()

	rawReq, err := http.NewRequest("GET", "http://www.example.com/hello/world?name=Lu&age=20", bytes.NewBuffer([]byte("hello=world")))
	if err != nil {
		t.Fatal(err)
	}
	rawReq.RemoteAddr = "127.0.0.1:1234"
	rawReq.Header.Add("Content-Type", "text/plain")

	req := NewRequest(rawReq)
	req.uri = "/hello/world?name=Lu&age=20"
	req.method = "GET"
	req.filePath = "hello.go"
	req.scheme = "http"
	req.host = "www.example.com"

	a.IsTrue(req.requestRemoteAddr() == "127.0.0.1")
	t.Log(req.requestRemotePort())
	a.IsTrue(req.requestRemotePort() == 1234)
	a.IsTrue(req.requestURI() == req.uri)
	a.IsTrue(req.requestPath() == "/hello/world")
	a.IsTrue(req.requestMethod() == "GET")
	a.IsTrue(req.requestLength() > 0)
	a.IsTrue(req.requestFilename() == req.filePath)
	a.IsTrue(req.requestProto() == "HTTP/1.1")
	a.IsTrue(req.requestQueryString() == "name=Lu&age=20")
	a.IsTrue(req.requestQueryParam("name") == "Lu")

	req.raw.Header["X-Real-IP"] = []string{"192.168.1.100"}
	a.IsTrue(req.requestRemoteAddr() == "192.168.1.100")

	delete(req.raw.Header, "X-Real-IP")
	req.raw.Header["X-Real-Ip"] = []string{"192.168.1.101"}
	a.IsTrue(req.requestRemoteAddr() == "192.168.1.101")

	delete(req.raw.Header, "X-Real-IP")
	delete(req.raw.Header, "X-Real-Ip")
	req.raw.Header["X-Forwarded-For"] = []string{"192.168.1.102, 192.168.1.103"}
	a.IsTrue(req.requestRemoteAddr() == "192.168.1.102")

	req.raw.Header["X-Forwarded-For"] = []string{"192.168.1.103"}
	a.IsTrue(req.requestRemoteAddr() == "192.168.1.103")

	t.Log(req.Format("hello ${teaVersion} remoteAddr:${remoteAddr} name:${arg.name} header:${header.Content-Type} test:${test}"))

	{
		req.host = "a.b.c.example.com"
		t.Log("===")
		t.Log(req.Format("host:${host} first:${host.first}, last:${host.last},  0:${host.0},  1:${host.1},  2:${host.2},  3:${host.3},  4:${host.4}"))
		t.Log(req.Format("-1:${host.-1} -2:${host.-2} -3:${host.-3} -4:${host.-4}"))
	}

	{
		req.host = "a.b.example.com"
		t.Log("===")
		t.Log(req.Format("host:${host} first:${host.first}, last:${host.last},  0:${host.0},  1:${host.1},  2:${host.2},  3:${host.3},  4:${host.4}"))
		t.Log(req.Format("-1:${host.-1} -2:${host.-2} -3:${host.-3} -4:${host.-4}"))
	}

	{
		req.host = "a.example.com"
		t.Log("===")
		t.Log(req.Format("host:${host} first:${host.first}, last:${host.last},  0:${host.0},  1:${host.1},  2:${host.2},  3:${host.3},  4:${host.4}"))
		t.Log(req.Format("-1:${host.-1} -2:${host.-2} -3:${host.-3} -4:${host.-4}"))
	}

	{
		req.host = "a.example.com"
		t.Log("===")
		t.Log(req.Format("host:${host} first:${host.first}, last:${host.last},  0:${host.0},  1:${host.1},  2:${host.2},  3:${host.3},  4:${host.4}"))
		t.Log(req.Format("-1:${host.-1} -2:${host.-2} -3:${host.-3} -4:${host.-4}"))
	}

	{
		req.host = "example.com"
		t.Log("===")
		t.Log(req.Format("host:${host} first:${host.first}, last:${host.last},  0:${host.0},  1:${host.1},  2:${host.2},  3:${host.3},  4:${host.4}"))
		t.Log(req.Format("-1:${host.-1} -2:${host.-2} -3:${host.-3} -4:${host.-4}"))
	}
}

func TestRequest_FormatPerformance(t *testing.T) {
	rawReq, err := http.NewRequest("GET", "http://www.example.com/hello/world?name=Lu&age=20", bytes.NewBuffer([]byte("hello=world")))
	if err != nil {
		t.Fatal(err)
	}
	rawReq.RemoteAddr = "127.0.0.1:1234"
	rawReq.Header.Add("Content-Type", "text/plain")

	req := NewRequest(rawReq)
	req.uri = "/hello/world?name=Lu&age=20"
	req.method = "GET"
	req.filePath = "hello.go"
	req.scheme = "http"

	count := 10000
	before := time.Now()
	result := ""
	for i := 0; i < count; i++ {
		for n := 0; n < 5; n++ {
			source := "hello ${teaVersion} remoteAddr:${remoteAddr} name:${arg.name} header:${header.Content-Type} test:${test} /hello " + fmt.Sprintf("%d", n)
			result = req.Format(source)
		}
	}

	cost := int(float64(count) / time.Since(before).Seconds())
	t.Log(cost)
	t.Log(result)
}

func TestRequest_Index(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()

	rawReq, err := http.NewRequest("GET", "http://www.example.com/hello/world?name=Lu&age=20", bytes.NewBuffer([]byte("hello=world")))
	if err != nil {
		t.Fatal(err)
	}

	req := NewRequest(rawReq)
	req.index = []string{}
	t.Log(req.findIndexFile(Tea.Root))

	req.index = []string{"main.go", "main2.go", "run.sh"}
	a.Equals(req.findIndexFile(Tea.Root), "main.go")

	req.index = []string{"main.*"}
	a.Equals(req.findIndexFile(Tea.Root), "main.go")
}

func TestRequest_LocationVariables(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()

	rawReq, err := http.NewRequest("GET", "http://www.example.com/hello/world?name=Lu&age=20", bytes.NewBuffer([]byte("hello=world")))
	if err != nil {
		t.Fatal(err)
	}

	server := teaconfigs.NewServerConfig()
	server.Root = "/home"

	{
		location := teaconfigs.NewLocation()
		location.On = true
		location.Pattern = "~ /hello/(\\w)(\\w+)"
		location.Root = "/hello/${1}/${host}"
		location.Index = []string{"hello_${1}${2}"}
		location.Charset = "${arg.charset}"
		location.AddResponseHeader(&shared.HeaderConfig{
			On: true, Name: "hello", Value: "${1}",
		})
		err := location.Validate()
		a.IsNil(err)

		server.AddLocation(location)

		matches, ok := location.Match("/hello/world", func(source string) string {
			return source
		})
		if ok {
			t.Log(matches)
		}
	}

	err = server.Validate()
	a.IsNil(err)

	req := NewRequest(rawReq)
	req.uri = "/hello/world?charset=utf-8"
	req.host = "www.example.com"

	err = req.configure(server, 0, false)
	if err != nil {
		t.Log(err.Error())
	}
	a.IsNil(err)

	t.Log("request uri:", req.requestURI())
	t.Log("root:", req.root)
	t.Log("index:", req.index)
	t.Log("charset:", req.charset)

	for _, header := range req.responseHeaders {
		t.Log("headers:", header.Name, ":", header.Value)
	}
}

func TestRequest_RewriteVariables(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()

	rawReq, err := http.NewRequest("GET", "http://www.example.com/hello/world?name=Lu&age=20", bytes.NewBuffer([]byte("hello=world")))
	if err != nil {
		t.Fatal(err)
	}

	server := teaconfigs.NewServerConfig()
	server.Root = "/home/${arg.charset}"
	server.Charset = "[${arg.charset}]"
	server.AddResponseHeader(&shared.HeaderConfig{
		Name:  "Charset",
		Value: "${arg.charset}",
	})

	{
		location := teaconfigs.NewLocation()
		location.On = true
		location.Pattern = "/"

		rewriteRule := teaconfigs.NewRewriteRule()
		rewriteRule.Pattern = "^/hello/(\\w+)$"
		rewriteRule.Replace = "/he/${1}${requestPath}?arg=${arg.charset}"
		location.AddRewriteRule(rewriteRule)

		err := location.Validate()
		a.IsNil(err)

		server.AddLocation(location)
	}

	err = server.Validate()
	a.IsNil(err)

	req := NewRequest(rawReq)
	req.uri = "/hello/world?charset=utf-8"
	req.host = "www.example.com"

	err = req.configure(server, 0, false)
	if err != nil {
		t.Log(err.Error())
	}
	a.IsNil(err)

	t.Log("request uri:", req.uri)
	t.Log("root:", req.root)
	t.Log("index:", req.index)
	t.Log("charset:", req.charset)

	for _, header := range req.responseHeaders {
		t.Log("headers:", header.Name, ":", header.Value)
	}
}

func BenchmarkPerformanceConfigure(b *testing.B) {
	runtime.GOMAXPROCS(1)

	rawReq, err := http.NewRequest("GET", "http://www.example.com/hello/world?name=Lu&age=20", bytes.NewBuffer([]byte("hello=world")))
	if err != nil {
		b.Fatal(err)
	}

	server := teaconfigs.NewServerConfig()

	backend := teaconfigs.NewBackendConfig()
	backend.Address = "127.0.0.1:1234"
	server.AddBackend(backend)

	{
		h := shared.NewHeaderConfig()
		h.Name = "TeaVersion"
		h.Value = "${teaVersion}"
		server.AddResponseHeader(h)
	}

	{
		h := shared.NewHeaderConfig()
		h.Name = "TeaPort"
		h.Value = "${remotePort}"
		server.AddResponseHeader(h)
	}

	{
		h := shared.NewHeaderConfig()
		h.Name = "TeaFile"
		h.Value = "${requestFilename}"
		server.AddResponseHeader(h)
	}

	{
		h := shared.NewHeaderConfig()
		h.Name = "Scheme"
		h.Value = "${scheme}"
		server.AddResponseHeader(h)
	}

	err = server.Validate()

	for i := 0; i < b.N; i++ {
		req := NewRequest(rawReq)
		req.uri = "/hello/world?charset=utf-8"
		req.host = "www.example.com"
		req.responseHeaders = []*shared.HeaderConfig{}

		err = req.configure(server, 0, false)
		if err != nil {
			b.Fatal(err.Error())
		}
	}
}

func TestPerformanceFormatHeaders(t *testing.T) {
	rawReq, err := http.NewRequest("GET", "http://www.example.com/hello/world?name=Lu&age=20", bytes.NewBuffer([]byte("hello=world")))
	if err != nil {
		t.Fatal(err)
	}

	server := teaconfigs.NewServerConfig()

	backend := teaconfigs.NewBackendConfig()
	backend.Address = "127.0.0.1:1234"
	server.AddBackend(backend)

	{
		h := shared.NewHeaderConfig()
		h.Name = "TeaVersion"
		h.Value = "${teaVersion}"
		server.AddResponseHeader(h)
	}

	{
		h := shared.NewHeaderConfig()
		h.Name = "TeaPort"
		h.Value = "${remotePort}"
		server.AddResponseHeader(h)
	}

	{
		h := shared.NewHeaderConfig()
		h.Name = "TeaFile"
		h.Value = "${requestFilename}"
		server.AddResponseHeader(h)
	}

	{
		h := shared.NewHeaderConfig()
		h.Name = "Scheme"
		h.Value = "${scheme}"
		server.AddResponseHeader(h)
	}

	err = server.Validate()

	count := 10000
	before := time.Now()

	for i := 0; i < count; i++ {
		req := NewRequest(rawReq)
		req.uri = "/hello/world?charset=utf-8"
		req.host = "www.example.com"
		req.responseHeaders = []*shared.HeaderConfig{}
	}

	cost := time.Since(before).Seconds()
	t.Log(float64(count)/cost, "qps")
}

func TestRequest_Format2(t *testing.T) {
	rawReq, err := http.NewRequest(http.MethodGet, "/hello?name=liu", nil)
	if err != nil {
		t.Fatal(err)
	}
	req := NewRequest(rawReq)
	req.uri = rawReq.URL.String()
	t.Log("arg.name:", req.Format("${arg.name}"))
}

func BenchmarkNewRequest(b *testing.B) {
	rawReq, err := http.NewRequest(http.MethodGet, "/hello?name=liu", nil)
	if err != nil {
		b.Fatal(err)
	}

	var req *Request

	for i := 0; i < b.N; i++ {
		req = NewRequest(rawReq)
		_ = req
	}
}

func BenchmarkParseURI(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = url.ParseRequestURI("http://teaos.cn/hello?name=liu")
	}
}
