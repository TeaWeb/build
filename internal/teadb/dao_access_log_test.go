package teadb

import (
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teadb/shared"
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/iwind/TeaGo/logs"
	stringutil "github.com/iwind/TeaGo/utils/string"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"strings"
	"testing"
	"time"
)

func TestAccessLogDAO_InsertOne(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	accessLog := newAccessLog()
	id, _ := shared.ObjectIdFromHex("5cfbbecd79c023a965148da9")
	accessLog.Id = id
	err := AccessLogDAO().InsertOne(accessLog)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestAccessLogDAO_InsertAccessLogs(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	list := []interface{}{}
	for i := 0; i < 5; i++ {
		list = append(list, newAccessLog())
	}
	err := AccessLogDAO().InsertAccessLogs(list)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(err)
}

func TestAccessLogDAO_FindAccessLogCookie(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	dao := AccessLogDAO()
	accessLog, err := dao.FindAccessLogCookie(timeutil.Format("Ymd"), "5cfbbecd79c023a965148da9")
	if err != nil {
		t.Fatal(err)
	}
	if accessLog == nil {
		t.Log("not found")
		return
	}
	t.Log(stringutil.JSONEncodePretty(accessLog.Cookie))
}

func TestAccessLogDAO_FindRequestHeaderAndBody(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	dao := AccessLogDAO()
	accessLog, err := dao.FindRequestHeaderAndBody(timeutil.Format("Ymd"), "5cfbbecd79c023a965148da9")
	if err != nil {
		t.Fatal(err)
	}
	if accessLog == nil {
		t.Log("not found")
		return
	}
	t.Log(accessLog.Header)
	if len(accessLog.RequestData) == 0 {
		t.Log(accessLog.RequestData)
	} else {
		t.Log(string(accessLog.RequestData))
	}
	logs.PrintAsJSON(accessLog, t)
}

func TestAccessLogDAO_FindResponseHeaderAndBody(t *testing.T) {
	dao := AccessLogDAO()
	accessLog, err := dao.FindResponseHeaderAndBody(timeutil.Format("Ymd"), "5cfbbecd79c023a965148da9")
	if err != nil {
		t.Fatal(err)
	}
	if accessLog == nil {
		t.Log("not found")
		return
	}
	t.Log(accessLog.SentHeader)
	t.Log(string(accessLog.ResponseBodyData))
	logs.PrintAsJSON(accessLog, t)
}

func TestAccessLogDAO_ListAccessLogs(t *testing.T) {
	{
		dao := AccessLogDAO()
		accessLogs, err := dao.ListAccessLogs(timeutil.Format("Ymd"), "5W8NLAoMYo6iJ78V", "", false, "", 0, 5)
		if err != nil {
			t.Fatal(err)
		}

		for _, accessLog := range accessLogs {
			t.Log(accessLog.Id, accessLog.ServerId, accessLog.Errors, accessLog.RemoteAddr)
			//logs.PrintAsJSON(accessLog, t)
		}
	}

	t.Log("=== from last id ===")
	{
		dao := AccessLogDAO()
		accessLogs, err := dao.ListAccessLogs(timeutil.Format("Ymd"), "5W8NLAoMYo6iJ78V", "5d73341f837f90ad48b60d3c", true, "", 0, 5)
		if err != nil {
			t.Fatal(err)
		}

		for _, accessLog := range accessLogs {
			t.Log(accessLog.Id, accessLog.ServerId, accessLog.Errors, accessLog.RemoteAddr)
		}
	}
}

func TestAccessLogDAO_ListAccessLogs_PastDays(t *testing.T) {
	dao := AccessLogDAO()
	_, err := dao.ListAccessLogs("201901", "5W8NLAoMYo6iJ78V", "", false, "", 0, 5)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccessLogDAO_HasNextAccessLog(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	dao := AccessLogDAO()
	ones, err := dao.ListTopAccessLogs(timeutil.Format("Ymd"), 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(ones) == 0 {
		t.Log("has no next")
		return
	}

	{
		b, err := dao.HasNextAccessLog(timeutil.Format("Ymd"), "5W8NLAoMYo6iJ78V", ones[0].Id.Hex(), false, "")
		if err != nil {
			t.Fatal(err)
		}
		if b {
			t.Log("has next")
		} else {
			t.Log("has no next")
		}
	}

	{
		b, err := dao.HasNextAccessLog(timeutil.Format("Ymd"), "5W8NLAoMYo6iJ78V", strings.Repeat("0", 24),
			false, "")
		if err != nil {
			t.Fatal(err)
		}
		if b {
			t.Log("has next")
		} else {
			t.Log("has no next")
		}
	}
}

func TestAccessLogDAO_HasAccessLog(t *testing.T) {
	{
		b, err := AccessLogDAO().HasAccessLog(timeutil.Format("Ymd"), "5W8NLAoMYo6iJ78V")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(b)
	}

	{
		b, err := AccessLogDAO().HasAccessLog(timeutil.Format("Ymd"), "5W8NLAoMYo6iJ78Y")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(b)
	}
}

func TestAccessLogDAO_HasAccessLogWithWAF(t *testing.T) {
	{
		b, err := AccessLogDAO().HasAccessLogWithWAF(timeutil.Format("Ymd"), "pq6HzRfIjcGsUqNe")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(b)
	}

	{
		b, err := AccessLogDAO().HasAccessLogWithWAF(timeutil.Format("Ymd"), "pq6HzRfIjcGsU123")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(b)
	}
}

func TestAccessLogDAO_GroupWAFRuleGroups(t *testing.T) {
	ruleSets, err := AccessLogDAO().GroupWAFRuleGroups(timeutil.Format("Ymd"), "pq6HzRfIjcGsUqNe")
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(ruleSets, t)
}

func TestAccessLogDAO_ListLatestAccessLogs(t *testing.T) {
	dao := AccessLogDAO()
	{
		accessLogs, err := dao.ListLatestAccessLogs(timeutil.Format("Ymd"), "5W8NLAoMYo6iJ78V", "", false, 5)
		if err != nil {
			t.Fatal(err)
		}
		for _, accessLog := range accessLogs {
			t.Log(accessLog.Id, accessLog.ServerId, accessLog.Errors, accessLog.RemoteAddr)
		}
	}

	t.Log("===from id===")
	{
		accessLogs, err := dao.ListLatestAccessLogs(timeutil.Format("Ymd"), "5W8NLAoMYo6iJ78V", "5cfbc98141a7eae69097db95", true, 5)
		if err != nil {
			t.Fatal(err)
		}
		for _, accessLog := range accessLogs {
			t.Log(accessLog.Id, accessLog.ServerId, accessLog.Errors, accessLog.RemoteAddr)
		}
	}
}

func TestAccessLogDAO_ListTopAccessLogs(t *testing.T) {
	dao := AccessLogDAO()
	accessLogs, err := dao.ListTopAccessLogs(timeutil.Format("Ymd"), 10)
	if err != nil {
		t.Fatal(err)
	}
	for _, accessLog := range accessLogs {
		t.Log(accessLog.Id, accessLog.ServerId, accessLog.Errors, accessLog.RemoteAddr)
	}
}

func TestAccessLogDAO_QueryAccessLogs(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	dao := AccessLogDAO()

	query := NewQuery("")
	query.Offset(1)
	query.Limit(5)
	query.Debug()

	accessLogs, err := dao.QueryAccessLogs(timeutil.Format("Ymd"), "5W8NLAoMYo6iJ78V", query)
	if err != nil {
		t.Fatal(err)
	}
	for _, accessLog := range accessLogs {
		t.Log(accessLog.Id, accessLog.ServerId, accessLog.Errors, accessLog.RemoteAddr)
	}
}

func newAccessLog() *accesslogs.AccessLog {
	accessLog := accesslogs.NewAccessLog()
	accessLog.ServerId = "5W8NLAoMYo6iJ78V"
	accessLog.BackendId = "backend.123456"
	accessLog.LocationId = "location.123456"
	accessLog.FastcgiId = "fastcgi.123456"
	accessLog.RewriteId = "rewrite.123456"
	accessLog.TeaVersion = teaconst.TeaVersion
	accessLog.RemoteAddr = "192.168.1.100"
	accessLog.RemotePort = 8080
	accessLog.RemoteUser = "user"
	accessLog.RequestURI = "/hello?world=1"
	accessLog.RequestPath = "/hello"
	accessLog.RequestMethod = "POST"
	accessLog.RequestFilename = "hello.txt"
	accessLog.RequestLength = 1024
	accessLog.RequestTime = 0.031
	accessLog.Scheme = "http"
	accessLog.Proto = "HTTP/1.1"
	accessLog.BytesSent = 1024
	accessLog.BodyBytesSent = 2048
	accessLog.Status = 200
	accessLog.StatusMessage = "OK"
	accessLog.SentHeader = map[string][]string{
		"Content-Type": {"text/plain"},
		"Hello":        {"World"},
	}
	accessLog.TimeISO8601 = "2019-06-28T10:49:28.134+08:00"
	accessLog.TimeLocal = "28/Jun/2019:10:49:28 +0800"
	accessLog.Msec = float64(time.Now().Unix()) + 0.052
	accessLog.Timestamp = time.Now().Unix()
	accessLog.Host = "www.teaos.cn"
	accessLog.Referer = "http://www.teaos.cn/index"
	accessLog.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36"
	accessLog.Request = "GET /hello?world=1 HTTP/1.1"
	accessLog.ContentType = "text/plain"
	accessLog.Cookie = map[string]string{
		"LOGIN_SHOP": "1",
		"sid":        "qs512D67HhAhqQuArUOK66rs5MjLLMPJ",
	}
	accessLog.Arg = map[string][]string{
		"world": {"1"},
	}
	accessLog.Args = "world=1"
	accessLog.QueryString = "world=1"
	accessLog.Header = map[string][]string{
		"Connection":                {"keep-alive"},
		"Cookie":                    {"LOGIN_SHOP=1; sid=qs512D67HhAhqQuArUOK66rs5MjLLMPJ",},
		"Upgrade-Insecure-Requests": {"1"},
		"User-Agent":                {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36",},
		"X-Forwarded-Proto":         {"http"},
		"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3",},
		"Accept-Language":           {"zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7,zh-TW;q=0.6,de;q=0.5,ja;q=0.4",},
		"X-Forwarded-By":            {"127.0.0.1:8882",},
		"X-Forwarded-For":           {"127.0.0.1",},
		"X-Forwarded-Host":          {"127.0.0.1:8882",},
		"X-Real-IP":                 {"127.0.0.1",},
		"Accept-Encoding":           {"gzip, deflate, br",},
	}
	accessLog.ServerName = "www.teaos.cn"
	accessLog.ServerPort = 8080
	accessLog.ServerProtocol = "HTTP/1.1"
	accessLog.BackendAddress = "127.0.0.1:9991"
	accessLog.FastcgiAddress = "127.0.0.1:9000"
	accessLog.RequestData = []byte("request data bytes")
	accessLog.ResponseHeaderData = []byte("response header data bytes")
	accessLog.ResponseBodyData = []byte("response body data bytes")
	accessLog.HasErrors = true
	accessLog.Errors = []string{"error1", "error2"}
	accessLog.Extend = &accesslogs.AccessLogExtend{}
	accessLog.Attrs = map[string]string{
		"cache": "1",
		"a":     "b",
	}
	return accessLog
}
