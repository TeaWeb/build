package notices

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
	"time"
)

func TestNoticeSetting_Notify(t *testing.T) {
	setting := SharedNoticeSetting()
	receiverIds := setting.Notify(NoticeLevelInfo, "subject", "消息内容第1行\n消息内容第2行", func(receiverId string, minutes int) int {
		return 0
	})
	t.Log(receiverIds)
	time.Sleep(3 * time.Second)
}

func TestNoticeSetting_FindAllNoticeReceivers(t *testing.T) {
	setting := &NoticeSetting{}
	setting.Levels = map[NoticeLevel]*NoticeLevelConfig{
		NoticeLevelWarning: {
			Receivers: []*NoticeReceiver{
				{
					User:    "zhang",
					MediaId: "1",
				},
				{
					User:    "zhang",
					MediaId: "1",
				},
				{
					User:    "wang",
					MediaId: "1",
				},
			},
		},
		NoticeLevelError: {
			Receivers: []*NoticeReceiver{
				{
					User:    "liu",
					MediaId: "1",
				},
				{
					User:    "liu",
					MediaId: "2",
				},
			},
		},
	}
	logs.PrintAsJSON(setting.FindAllNoticeReceivers(NoticeLevelWarning, NoticeLevelError), t)
}
