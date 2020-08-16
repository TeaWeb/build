package notices

import (
	"github.com/TeaWeb/build/internal/teatesting"
	"testing"
)

func TestNewNoticeQyWeixinMedia(t *testing.T) {
	if teatesting.IsGlobal() {
		return
	}
	m := NewNoticeQyWeixinMedia()
	m.CorporateId = "xxx"
	m.AppSecret = "xxx"
	m.AgentId = "1000003"
	resp, err := m.Send("", "标题：报警标题", "内容：报警内容/全员都有")
	if err != nil {
		t.Log(string(resp))
		t.Fatal(err)
	}
	t.Log(string(resp))
}
