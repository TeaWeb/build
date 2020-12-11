package accesslogs

import (
	"bytes"
	"encoding/json"
	"github.com/TeaWeb/build/internal/teadb/shared"
	"github.com/TeaWeb/build/internal/teageo"
	"github.com/TeaWeb/uaparser"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/assert"
	"github.com/mailru/easyjson"
	"log"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestLogParseMimeType(t *testing.T) {
	accessLog := &AccessLog{
		ContentType: "text/html;charset=utf-8",
	}
	accessLog.parseMime()
	t.Log(accessLog.Extend.File.MimeType)
	if accessLog.Extend.File.MimeType != "text/html" {
		t.Error("[ERROR]")
	}

	if accessLog.Extend.File.Charset != "UTF-8" {
		t.Error("[ERROR]", accessLog.Extend.File.Charset)
	}

	accessLog.ContentType = "text/html"
	accessLog.parseMime()
	if accessLog.Extend.File.MimeType != "text/html" {
		t.Error("[ERROR]", accessLog.Extend.File.MimeType)
	}

	if accessLog.Extend.File.Charset != "" {
		t.Error("[ERROR]", accessLog.Extend.File.Charset)
	}

	accessLog.ContentType = "text/html; CHARSET=UTF-8"
	accessLog.parseMime()
	if accessLog.Extend.File.Charset != "UTF-8" {
		t.Error("[ERROR]", accessLog.Extend.File.Charset)
	}
}

func TestLogParseExtension(t *testing.T) {
	accessLog := &AccessLog{
		RequestPath: "/hello",
	}

	accessLog.parseExtension()
	if accessLog.Extend.File.Extension != "" {
		t.Error("[ERROR]", accessLog.Extend.File.Extension)
	}

	accessLog.RequestPath = "/hello.js"
	accessLog.parseExtension()
	if accessLog.Extend.File.Extension != "js" {
		t.Error("[ERROR]", accessLog.Extend.File.Extension)
	}

	accessLog.RequestPath = "/hello.JS"
	accessLog.parseExtension()
	if accessLog.Extend.File.Extension != "js" {
		t.Error("[ERROR]", accessLog.Extend.File.Extension)
	}

	accessLog.RequestPath = "/hello.tar.gz"
	accessLog.parseExtension()
	if accessLog.Extend.File.Extension != "gz" {
		t.Error("[ERROR]", accessLog.Extend.File.Extension)
	}
}

func TestLogOSParser1(t *testing.T) {
	userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.59 Safari/537.36"
	parser, err := uaparser.NewParser(Tea.Root + "/web/resources" + Tea.DS + "regexes.yaml")
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	client, found := parser.Parse(userAgent)
	if !found {
		t.Log("not found")
		return
	}

	t.Log(client.Browser.Family) // "Amazon Silk"
	t.Log(client.Browser.Major)  // "1"
	t.Log(client.Browser.Minor)  // "1"
	t.Log(client.Browser.Patch)  // "0-80"
	t.Log(client.OS.Family)      // "Android"
	t.Log(client.OS.Major)       // ""
	t.Log(client.OS.Minor)       // ""
	t.Log(client.OS.Patch)       // ""
	t.Log(client.OS.PatchMinor)  // ""
	t.Log(client.Device.Family)  // "Kindle Fire"

	cost := time.Since(now).Seconds()
	t.Log("cost:", cost, "s")
	t.Log("QPS", 1/cost)
}

func TestLogOSParser2(t *testing.T) {
	userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.59 Safari/537.36"
	parser, err := uaparser.NewParser(Tea.Root + "/web/resources" + Tea.DS + "regexes.yaml")
	if err != nil {
		log.Fatal(err)
	}

	parser.Parse(userAgent) // 看看是否有缓存

	now := time.Now()
	agent, _ := parser.Parse(userAgent)
	cost := time.Since(now).Seconds()

	t.Logf("%#v", agent)
	t.Log("cost:", cost, "s")
	t.Log("QPS", 1/cost)
}

func TestLogOSParser3(t *testing.T) {
	userAgent := " Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1)"
	parser, err := uaparser.NewParser(Tea.Root + "/web/resources" + Tea.DS + "regexes.yaml")
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	agent, _ := parser.Parse(userAgent)
	t.Logf("%#v", agent)

	os, _ := parser.ParseOS(userAgent)
	t.Logf("%#v", os)

	cost := float64(time.Since(now).Nanoseconds()) / 1000000000
	t.Log("cost:", cost)
	t.Log("QPS", 1/cost)
}

func TestLogOSParser4(t *testing.T) {
	userAgent := "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36"
	parser, err := uaparser.NewParser(Tea.Root + "/web/resources" + Tea.DS + "regexes.yaml")
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	agent, _ := parser.ParseBrowser(userAgent)
	t.Logf("%#v", agent)

	os, _ := parser.ParseOS(userAgent)
	t.Logf("%#v", os)

	cost := float64(time.Since(now).Nanoseconds()) / 1000000000
	t.Log("cost:", cost)
	t.Log("QPS", 1/cost)
}

func TestLogParse5(t *testing.T) {
	parser, err := uaparser.NewParser(Tea.Root + Tea.DS + "web" + Tea.DS + "resources" + Tea.DS + "regexes.yaml")
	if err != nil {
		t.Fatal(err)
	}
	userAgentParser = parser

	userAgents := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.59 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:60.0) Gecko/20100101 Firefox/60.0",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1)",
		"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36",
		"",
		"hello",
	}

	for _, userAgent := range userAgents {
		accessLog := &AccessLog{
			UserAgent: userAgent,
		}
		accessLog.parseUserAgent()

		t.Log("=======")
		t.Log(userAgent)
		if accessLog.Extend != nil {
			t.Logf("%#v", accessLog.Extend.Client)
		}
	}
}

func TestAccessLog_nil(t *testing.T) {
	accessLog := &AccessLog{}
	t.Log(string(accessLog.jsonEncode(1)))
	t.Log(string(accessLog.jsonEncode("abc")))
	t.Log(string(accessLog.jsonEncode(nil)))

	{
		var h *AccessLogClient = nil
		t.Log(string(accessLog.jsonEncode(h)))
	}

	{
		var h = new(AccessLogClient)
		t.Log(string(accessLog.jsonEncode(h)))
	}

}

func BenchmarkAccessLog_ParseUserAgent(b *testing.B) {
	runtime.GOMAXPROCS(1)

	parser, err := uaparser.NewParser(Tea.Root + Tea.DS + "web" + Tea.DS + "resources" + Tea.DS + "regexes.yaml")
	if err != nil {
		b.Fatal(err)
	}
	userAgentParser = parser
	teageo.SetupDB()

	for i := 0; i < b.N; i++ {
		accessLog := &AccessLog{}
		accessLog.UserAgent = "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36"
		accessLog.RemoteAddr = "160.16.241." + strconv.Itoa(i%200)
		accessLog.ContentType = "text/html; charset=utf-8"
		accessLog.RequestPath = "/hello.php"
		//accessLog.SetWritingFields([]int{AccessLogFieldHeader})
		accessLog.Parse()
	}
}

func TestAccessLogger_DB(t *testing.T) {
	accessLog := &AccessLog{
		Id:   shared.NewObjectId(),
		Args: "a=b",
		Arg: map[string][]string{
			"name": {"liu", "lu"},
		},
		Cookie: map[string]string{
			"sid": "123456",
		},
		RemoteAddr:    "127.0.0.1",
		RemotePort:    80,
		TimeLocal:     "23/Jul/2018:22:23:35 +0800",
		TimeISO8601:   "2018-07-23T22:23:35+08:00",
		Status:        200,
		BodyBytesSent: 1048,
		Request:       "GET / HTTP/1.1",
	}
	data, err := easyjson.Marshal(accessLog)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}

func TestAccessLog_Format(t *testing.T) {
	accessLog := AccessLog{
		Args: "a=b",
		Arg: map[string][]string{
			"name": {"liu", "lu"},
		},
		Cookie: map[string]string{
			"sid": "123456",
		},
		RemoteAddr:    "127.0.0.1",
		RemotePort:    80,
		TimeLocal:     "23/Jul/2018:22:23:35 +0800",
		TimeISO8601:   "2018-07-23T22:23:35+08:00",
		Status:        200,
		BodyBytesSent: 1048,
		Request:       "GET / HTTP/1.1",
		Header: map[string][]string{
			"User-Agent": {
				"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.68 Safari/537.36",
			},
			"Referer": {
				"https://www.baidu.com/",
			},
		},
	}

	accessLog.Parse()

	format := "${args} ${arg.name} ${cookie.sid} ${remoteAddr} - [${timeLocal}] \"${request}\" ${status} ${bodyBytesSent} \"${http.Referer}\" \"${http.UserAgent}\""

	for i := 0; i < 2; i++ {
		before := time.Now()
		t.Log(accessLog.Format(format))
		t.Log(time.Since(before).Seconds()*1000, "ms")
	}

	format = "Extend:${extend.File} ${extend.Geo}"
	t.Log(accessLog.Format(format))
}

func TestAccessLog_FormatFloat(t *testing.T) {
	accessLog := &AccessLog{
		RequestTime: 7.964e-5,
	}
	t.Log(accessLog.Format("${requestTime}"))
}

func TestAccessLog_ParseGEO(t *testing.T) {
	teageo.SetupDB()

	{
		accessLog := &AccessLog{
			RemoteAddr: "114.240.210.253",
		}
		accessLog.parseGeoIP()
	}

	before := time.Now()
	accessLog := &AccessLog{
		RemoteAddr: "114.240.210.253",
	}
	accessLog.parseGeoIP()
	cost := time.Since(before).Seconds()

	if accessLog.Extend != nil {
		t.Logf("%#v", accessLog.Extend.Geo)
	}
	t.Log(1 / cost)
}

func TestAccessLog_CleanFields(t *testing.T) {
	a := assert.NewAssertion(t)

	accessLog := NewAccessLog()
	accessLog.UserAgent = "123"
	accessLog.Header = map[string][]string{}
	accessLog.Cookie = map[string]string{}
	accessLog.writingFields = []int{AccessLogFieldHeader, AccessLogFieldArg, AccessLogFieldCookie}
	accessLog.CleanFields()
	t.Log(accessLog.UserAgent)
	t.Log(accessLog.Header)
	t.Log(accessLog.Cookie)

	a.IsNil(accessLog.Extend)
}

func TestAccessLog_JSON(t *testing.T) {
	before := time.Now()
	accessLog := NewAccessLog()
	accessLog.Request = "GET /hello"
	data, err := easyjson.Marshal(accessLog)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(time.Since(before).Seconds()*1000, "ms")
	t.Log(string(data))
}

func TestAccessLog_Backend(t *testing.T) {
	accessLog := NewAccessLog()
	accessLog.BackendId = "1234"
	accessLog.BackendAddress = "127.0.0.1:8001"
	accessLog.BackendCode = "web001"
	accessLog.BackendScheme = "http"
	accessLog.Hostname = "Web001"
	t.Log(accessLog.Format("${backend.address} ${backend.host} ${backend.id} ${backend.code} ${backend.scheme} ${hostname}"))
	data, err := easyjson.Marshal(accessLog)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}

func TestAccessLog_JSON_compare(t *testing.T) {
	accessLog := &AccessLog{
		RequestPath:     "/hello",
		RequestURI:      "/hello?name=lu&age=20",
		RequestMethod:   "POST",
		RequestFilename: "/hello.html",
		Header: map[string][]string{
			"UserAgent": {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36"},
		},
	}
	data, _ := easyjson.Marshal(accessLog)
	t.Log(string(data))

	data1, _ := json.Marshal(accessLog)
	t.Log(string(data1))

	t.Log(bytes.Compare(data, data1))
}

func BenchmarkAccessLog_JSON_easyjson(b *testing.B) {
	runtime.GOMAXPROCS(1)
	for i := 0; i < b.N; i++ {
		accessLog := &AccessLog{
			RequestPath:     "/hello",
			RequestURI:      "/hello?name=lu&age=20",
			RequestMethod:   "POST",
			RequestFilename: "/hello.html",
			Header: map[string][]string{
				"UserAgent": {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36"},
			},
		}
		_, _ = easyjson.Marshal(accessLog)
	}
}

func BenchmarkAccessLog_JSON_raw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		accessLog := NewAccessLog()
		accessLog.Request = "GET /hello"
		accessLog.Host = "teaos.cn"
		_, _ = json.Marshal(accessLog)
	}
}
