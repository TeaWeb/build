package notices

import (
	"github.com/iwind/TeaGo/maps"
	"net/http"
	"testing"
)

func TestNoticeMediaConfig_Raw(t *testing.T) {
	{
		c := NewNoticeMediaConfig()
		c.Type = NoticeMediaTypeWebhook
		c.Options = maps.Map{
			"url":    "http://example.com/json",
			"method": http.MethodGet,
		}
		raw, err := c.Raw()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%#v", raw)
	}

	{
		c := NewNoticeMediaConfig()
		c.Type = NoticeMediaTypeScript
		c.Options = maps.Map{
			"path": "/opt/www/notify.sh",
		}
		raw, err := c.Raw()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%#v", raw)
	}

	{
		c := NewNoticeMediaConfig()
		c.Type = NoticeMediaTypeEmail
		c.Options = maps.Map{
			"smtp": "smtp.qq.com",
		}
		raw, err := c.Raw()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%#v", raw)
	}
}
