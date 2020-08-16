package teawaf

import (
	"bytes"
	"github.com/iwind/TeaGo/assert"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func Test_Template(t *testing.T) {
	a := assert.NewAssertion(t)

	template := Template()
	err := template.Init()
	if err != nil {
		t.Fatal(err)
	}

	template.OnAction(func(action ActionString) (goNext bool) {
		return action != ActionBlock
	})

	testTemplate1001(a, t, template)
	testTemplate1002(a, t, template)
	testTemplate1003(a, t, template)
	testTemplate2001(a, t, template)
	testTemplate3001(a, t, template)
	testTemplate4001(a, t, template)
	testTemplate5001(a, t, template)
	testTemplate6001(a, t, template)
	testTemplate7001(a, t, template)
	testTemplate20001(a, t, template)
}

func Test_Template2(t *testing.T) {
	reader := bytes.NewReader([]byte(strings.Repeat("HELLO", 1024)))
	req, err := http.NewRequest(http.MethodGet, "http://example.com/index.php?id=123", reader)
	if err != nil {
		t.Fatal(err)
	}

	waf := Template()
	err = waf.Init()
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	goNext, _, set, err := waf.MatchRequest(req, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(time.Since(now).Seconds()*1000, "ms")

	if goNext {
		t.Log("ok")
		return
	}

	logs.PrintAsJSON(set, t)
}

func BenchmarkTemplate(b *testing.B) {
	waf := Template()
	err := waf.Init()
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader([]byte(strings.Repeat("Hello", 1024)))
		req, err := http.NewRequest(http.MethodGet, "http://example.com/index.php?id=123", reader)
		if err != nil {
			b.Fatal(err)
		}

		_, _, _, _ = waf.MatchRequest(req, nil)
	}
}

func testTemplate1001(a *assert.Assertion, t *testing.T, template *WAF) {
	req, err := http.NewRequest(http.MethodGet, "http://example.com/index.php?id=onmousedown%3D123", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, _, result, err := template.MatchRequest(req, nil)
	if err != nil {
		t.Fatal(err)
	}
	a.IsNotNil(result)
	if result != nil {
		a.IsTrue(result.Code == "1001")
	}
}

func testTemplate1002(a *assert.Assertion, t *testing.T, template *WAF) {
	req, err := http.NewRequest(http.MethodGet, "http://example.com/index.php?id=eval%28", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, _, result, err := template.MatchRequest(req, nil)
	if err != nil {
		t.Fatal(err)
	}
	a.IsNotNil(result)
	if result != nil {
		a.IsTrue(result.Code == "1002")
	}
}

func testTemplate1003(a *assert.Assertion, t *testing.T, template *WAF) {
	req, err := http.NewRequest(http.MethodGet, "http://example.com/index.php?id=<script src=\"123.js\">", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, _, result, err := template.MatchRequest(req, nil)
	if err != nil {
		t.Fatal(err)
	}
	a.IsNotNil(result)
	if result != nil {
		a.IsTrue(result.Code == "1003")
	}
}

func testTemplate2001(a *assert.Assertion, t *testing.T, template *WAF) {
	body := bytes.NewBuffer([]byte{})

	writer := multipart.NewWriter(body)

	{
		part, err := writer.CreateFormField("name")
		if err == nil {
			_, _ = part.Write([]byte("lu"))
		}
	}

	{
		part, err := writer.CreateFormField("age")
		if err == nil {
			_, _ = part.Write([]byte("20"))
		}
	}

	{
		part, err := writer.CreateFormFile("myFile", "hello.txt")
		if err == nil {
			_, _ = part.Write([]byte("Hello, World!"))
		}
	}

	{
		part, err := writer.CreateFormFile("myFile2", "hello.PHP")
		if err == nil {
			_, _ = part.Write([]byte("Hello, World, PHP!"))
		}
	}

	{
		part, err := writer.CreateFormFile("myFile3", "hello.asp")
		if err == nil {
			_, _ = part.Write([]byte("Hello, World, ASP Pages!"))
		}
	}

	{
		part, err := writer.CreateFormFile("myFile4", "hello.asp")
		if err == nil {
			_, _ = part.Write([]byte("Hello, World, ASP Pages!"))
		}
	}

	_ = writer.Close()

	req, err := http.NewRequest(http.MethodPost, "http://teaos.cn/", body)
	if err != nil {
		t.Fatal()
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	_, _, result, err := template.MatchRequest(req, nil)
	if err != nil {
		t.Fatal(err)
	}
	a.IsNotNil(result)
	if result != nil {
		a.IsTrue(result.Code == "2001")
	}
}

func testTemplate3001(a *assert.Assertion, t *testing.T, template *WAF) {
	req, err := http.NewRequest(http.MethodPost, "http://example.com/index.php?exec1+(", bytes.NewReader([]byte("exec('rm -rf /hello');")))
	if err != nil {
		t.Fatal(err)
	}
	_, _, result, err := template.MatchRequest(req, nil)
	if err != nil {
		t.Fatal(err)
	}
	a.IsNotNil(result)
	if result != nil {
		a.IsTrue(result.Code == "3001")
	}
}

func testTemplate4001(a *assert.Assertion, t *testing.T, template *WAF) {
	req, err := http.NewRequest(http.MethodPost, "http://example.com/index.php?whoami", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, _, result, err := template.MatchRequest(req, nil)
	if err != nil {
		t.Fatal(err)
	}
	a.IsNotNil(result)
	if result != nil {
		a.IsTrue(result.Code == "4001")
	}
}

func testTemplate5001(a *assert.Assertion, t *testing.T, template *WAF) {
	{
		req, err := http.NewRequest(http.MethodPost, "http://example.com/.././..", nil)
		if err != nil {
			t.Fatal(err)
		}
		_, _, result, err := template.MatchRequest(req, nil)
		if err != nil {
			t.Fatal(err)
		}
		a.IsNotNil(result)
		if result != nil {
			a.IsTrue(result.Code == "5001")
		}
	}

	{
		req, err := http.NewRequest(http.MethodPost, "http://example.com/..///./", nil)
		if err != nil {
			t.Fatal(err)
		}
		_, _, result, err := template.MatchRequest(req, nil)
		if err != nil {
			t.Fatal(err)
		}
		a.IsNotNil(result)
		if result != nil {
			a.IsTrue(result.Code == "5001")
		}
	}
}

func testTemplate6001(a *assert.Assertion, t *testing.T, template *WAF) {
	{
		req, err := http.NewRequest(http.MethodPost, "http://example.com/.svn/123.txt", nil)
		if err != nil {
			t.Fatal(err)
		}
		_, _, result, err := template.MatchRequest(req, nil)
		if err != nil {
			t.Fatal(err)
		}
		a.IsNotNil(result)
		if result != nil {
			a.IsTrue(result.Code == "6001")
		}
	}

	{
		req, err := http.NewRequest(http.MethodPost, "http://example.com/123.git", nil)
		if err != nil {
			t.Fatal(err)
		}
		_, _, result, err := template.MatchRequest(req, nil)
		if err != nil {
			t.Fatal(err)
		}
		a.IsNil(result)
	}
}

func testTemplate7001(a *assert.Assertion, t *testing.T, template *WAF) {
	for _, id := range []string{
		"union select",
		" and if(",
		"/*!",
		" and select ",
		" and id=123 ",
		"(case when a=1 then ",
		"updatexml (",
		"; delete from table",
	} {
		req, err := http.NewRequest(http.MethodPost, "http://example.com/?id="+url.QueryEscape(id), nil)
		if err != nil {
			t.Fatal(err)
		}
		_, _, result, err := template.MatchRequest(req, nil)
		if err != nil {
			t.Fatal(err)
		}
		a.IsNotNil(result)
		if result != nil {
			a.IsTrue(lists.ContainsAny([]string{"7001", "7002", "7003", "7004", "7005"}, result.Code))
		} else {
			t.Log("break:", id)
		}
	}
}

func testTemplate20001(a *assert.Assertion, t *testing.T, template *WAF) {
	// enable bot rule set
	for _, g := range template.Inbound {
		if g.Code == "bot" {
			g.On = true
			break
		}
	}

	for _, bot := range []string{
		"Googlebot",
		"AdsBot",
		"bingbot",
		"BingPreview",
		"facebookexternalhit",
		"Slurp",
		"Sogou",
		"Baiduspider http://baidu.com",
	} {
		req, err := http.NewRequest(http.MethodPost, "http://example.com/", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("User-Agent", bot)
		_, _, result, err := template.MatchRequest(req, nil)
		if err != nil {
			t.Fatal(err)
		}
		a.IsNotNil(result)
		if result != nil {
			a.IsTrue(lists.ContainsAny([]string{"20001"}, result.Code))
		} else {
			t.Log("break:", bot)
		}
	}
}
