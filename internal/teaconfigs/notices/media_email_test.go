package notices

import (
	"github.com/TeaWeb/build/internal/teatesting"
	"testing"
)

func TestNoticeEmailMedia_Send(t *testing.T) {
	if teatesting.IsGlobal() {
		return
	}
	media := NewNoticeEmailMedia()
	media.SMTP = "smtp.qq.com:587"
	media.Username = "19644627@qq.com"
	media.Password = "123456" // 换成你的邮件密码或者授权码
	media.From = "19644627@qq.com"
	_, err := media.Send("iwind.liu@gmail.com", "This is test subject", "This is a test body <strong>粗体哦</strong><br/>换行哦")
	if err != nil {
		t.Fatal(err)
	}
}

func TestNoticeEmailMedia_Send_163(t *testing.T) {
	if teatesting.IsGlobal() {
		return
	}
	media := NewNoticeEmailMedia()
	media.SMTP = "smtp.163.com:465"
	media.Username = "iwind_php@163.com"
	media.Password = "123456" // 换成你的邮件密码或者授权码
	media.From = "iwind_php@163.com"
	_, err := media.Send("iwind_php@163.com", "This is test subject", "This is a test body <strong>粗体哦</strong><br/>换行哦")
	if err != nil {
		t.Fatal(err)
	}
}
