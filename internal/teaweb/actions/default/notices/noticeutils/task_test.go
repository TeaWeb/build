package noticeutils

import (
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"testing"
)

func TestAddTask(t *testing.T) {
	AddTask(notices.NoticeLevelWarning, []*notices.NoticeReceiver{}, "subject", "message")
}
