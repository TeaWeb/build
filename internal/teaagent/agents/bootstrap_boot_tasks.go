package teaagents

import "github.com/iwind/TeaGo/logs"

// 启动任务
func bootTasks() {
	logs.Println("booting ...")
	if !runningAgent.On {
		return
	}
	for _, app := range runningAgent.Apps {
		if !app.On {
			continue
		}
		for _, taskConfig := range app.Tasks {
			if !taskConfig.On {
				continue
			}
			task := NewTask(app.Id, taskConfig)
			if task.ShouldBoot() {
				err := task.RunLog()
				if err != nil {
					logs.Println(err.Error())
				}
			}
		}
	}
}
