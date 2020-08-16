package agent

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"testing"
)

func TestPushAction_FindNoticeReceivers(t *testing.T) {
	action := new(PushAction)

	{
		receivers := action.findNoticeReceivers(agents.NewAgentConfigFromId("local"), "system", nil, []notices.NoticeLevel{notices.NoticeLevelWarning})
		for _, receiver := range receivers {
			t.Log("warning:", receiver.Name, receiver.User)
		}
	}

	{
		receivers := action.findNoticeReceivers(agents.NewAgentConfigFromId("local"), "system", nil, []notices.NoticeLevel{notices.NoticeLevelError})
		for _, receiver := range receivers {
			t.Log("error:", receiver.Name, receiver.User)
		}
	}

	{
		receivers := action.findNoticeReceivers(agents.NewAgentConfigFromId("local"), "system", nil, []notices.NoticeLevel{notices.NoticeLevelSuccess})
		if len(receivers) == 0 {
			t.Log("success:", "not found")
		} else {
			for _, receiver := range receivers {
				t.Log("success:", receiver.Name, receiver.User)
			}
		}
	}

	{
		receivers := action.findNoticeReceivers(agents.NewAgentConfigFromId("local"), "system", nil, []notices.NoticeLevel{notices.NoticeLevelSuccess, notices.NoticeLevelWarning, notices.NoticeLevelError})
		if len(receivers) == 0 {
			t.Log("all:", "not found")
		} else {
			for _, receiver := range receivers {
				t.Log("all:", receiver.Name, receiver.User)
			}
		}
	}
}
