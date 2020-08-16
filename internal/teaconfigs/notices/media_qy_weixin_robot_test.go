package notices

import (
	"github.com/TeaWeb/build/internal/teatesting"
	"testing"
)

func TestNoticeQyWeixinRobotMedia_Send(t *testing.T) {
	if teatesting.IsGlobal() {
		return
	}
	media := NewNoticeQyWeixinRobotMedia()
	media.WebhookURL = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=123456" //需要换成你自己的webhook
	media.TextFormat = FormatText
	resp, err := media.Send("", "这是标题", "*这是内容*")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("resp:", string(resp))
}
