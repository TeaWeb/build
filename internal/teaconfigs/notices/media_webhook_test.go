package notices

import (
	"github.com/TeaWeb/build/internal/teatesting"
	"testing"
)

func TestNoticeMediaWebhook_Send(t *testing.T) {
	if !teatesting.RequireHTTPServer() {
		return
	}

	media := NewNoticeWebhookMedia()
	media.URL = "http://127.0.0.1:9991/webhook?subject=${NoticeSubject}&body=${NoticeBody}"
	resp, err := media.Send("zhangsan", "this is subject", "this is body")
	t.Log(string(resp), err)
}
