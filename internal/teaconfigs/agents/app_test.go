package agents

import (
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestAppConfig_FindAllNoticeReceivers(t *testing.T) {
	app := NewAppConfig()
	app.NoticeSetting[notices.NoticeLevelWarning] = []*notices.NoticeReceiver{
		{
			MediaId: "1",
			User:    "zhang",
		},
		{
			MediaId: "1",
			User:    "wang",
		},
		{
			MediaId: "1",
			User:    "liu",
		},
	}
	app.NoticeSetting[notices.NoticeLevelError] = []*notices.NoticeReceiver{
		{
			MediaId: "1",
			User:    "zhang2",
		},
		{
			MediaId: "2",
			User:    "liu",
		},
		{
			MediaId: "1",
			User:    "zhang",
		},
	}
	t.Log(app.FindAllNoticeReceivers(notices.NoticeLevelError))
	logs.PrintAsJSON(app.FindAllNoticeReceivers(notices.NoticeLevelWarning, notices.NoticeLevelError), t)
}
