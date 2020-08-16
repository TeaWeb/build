package teaagents

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"regexp"
)

// 定时任务
func scheduleTasks() error {
	// 生成脚本
	taskIds := []string{}

	for _, app := range runningAgent.Apps {
		if !app.On {
			continue
		}
		for _, taskConfig := range app.Tasks {
			if !taskConfig.On {
				continue
			}
			taskIds = append(taskIds, taskConfig.Id)

			// 是否正在运行
			runningTask, found := runningTasks[taskConfig.Id]
			isChanged := true
			if found {
				// 如果有修改，则需要重启
				if runningTask.config.Version != taskConfig.Version {
					logs.Println("stop schedule task", taskConfig.Id, taskConfig.Name)
					runningTask.Stop()

					if taskConfig.On && len(taskConfig.Schedule) > 0 {
						logs.Println("restart schedule task", taskConfig.Id, taskConfig.Name)
						runningTask.config = taskConfig
						runningTask.Schedule()
					}
				} else {
					isChanged = false
				}
			} else if taskConfig.On && len(taskConfig.Schedule) > 0 { // 新任务，则启动
				logs.Println("schedule task", taskConfig.Id, taskConfig.Name)
				task := NewTask(app.Id, taskConfig)
				task.Schedule()

				runningTasksLocker.Lock()
				runningTasks[taskConfig.Id] = task
				runningTasksLocker.Unlock()
			}

			// 生成脚本
			if isChanged {
				_, err := taskConfig.GenerateAgain()
				if err != nil {
					return err
				}
			}
		}
	}

	// 停止运行
	for taskId, runningTask := range runningTasks {
		if !lists.Contains(taskIds, taskId) {
			runningTasksLocker.Lock()
			delete(runningTasks, taskId)
			runningTasksLocker.Unlock()
			err := runningTask.Stop()
			if err != nil {
				logs.Error(err)
			}
		}
	}

	// 删除不存在的任务脚本
	files.NewFile(Tea.ConfigFile("agents/")).Range(func(file *files.File) {
		filename := file.Name()

		for _, ext := range []string{"script", "bat"} {
			if regexp.MustCompile("^task\\.\\w+\\." + ext + "$").MatchString(filename) {
				taskId := filename[len("task:") : len(filename)-len("."+ext)]
				if !lists.Contains(taskIds, taskId) {
					err := file.Delete()
					if err != nil {
						logs.Error(err)
					}
				}
			}
		}
	})

	return nil
}

// 监控数据采集
func scheduleItems() error {
	logs.Println("schedule items")
	itemIds := []string{}

	for _, app := range runningAgent.Apps {
		if !app.On {
			continue
		}
		for _, itemConfig := range app.Items {
			if !itemConfig.On {
				continue
			}
			runningItemsLocker.Lock()
			itemIds = append(itemIds, itemConfig.Id)
			runningItem, found := runningItems[itemConfig.Id]
			if found {
				runningItem.Stop()
			}

			item := NewItem(app.Id, itemConfig)
			item.Schedule()
			runningItems[itemConfig.Id] = item
			logs.Println("add item", item.config.Name)
			runningItemsLocker.Unlock()
		}
	}

	// 删除不运行的
	for itemId, item := range runningItems {
		if !lists.Contains(itemIds, itemId) {
			item.Stop()
			runningItemsLocker.Lock()
			delete(runningItems, itemId)
			logs.Println("delete item", item.config.Name)
			runningItemsLocker.Unlock()
		}
	}

	return nil
}
