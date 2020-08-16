package noticeutils

import (
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/iwind/TeaGo/logs"
)

// 发送任务
type PushTask struct {
	Level     notices.NoticeLevel
	Receivers []*notices.NoticeReceiver
	Subject   string
	Message   string
}

var taskCh = make(chan *PushTask, 1024)

// 添加任务
func AddTask(level notices.NoticeLevel, receivers []*notices.NoticeReceiver, subject string, message string) {
	taskCh <- &PushTask{
		Level:     level,
		Receivers: receivers,
		Subject:   subject,
		Message:   message,
	}
}

// 运行任务
func RunTasks() {
	go func() {
		for task := range taskCh {
			if task == nil || len(task.Receivers) == 0 {
				continue
			}
			logs.Println("[notice]push " + task.Subject)
			setting := notices.SharedNoticeSetting()
			setting.NotifyReceivers(task.Level, task.Receivers, task.Subject, task.Message, func(receiverId string, minutes int) int {
				count, err := teadb.NoticeDAO().CountReceivedNotices(receiverId, map[string]interface{}{}, minutes)
				if err != nil {
					logs.Error(err)
				}
				return count
			})
		}
	}()
}
