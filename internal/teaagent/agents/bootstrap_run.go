package teaagents

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/iwind/TeaGo/logs"
	"os"
)

// 运行任务或者监控项
func runTaskOrItem() {
	if len(os.Args) <= 2 {
		logs.Println("no task to run")
		return
	}

	taskId := os.Args[2]
	if len(taskId) == 0 {
		logs.Println("no task to run")
		return
	}

	agent := agents.NewAgentConfigFromId(connectConfig.Id)
	if agent == nil {
		logs.Println("agent not found")
		return
	}
	appConfig, taskConfig := agent.FindTask(taskId)
	if taskConfig == nil {
		// 查找Item
		appConfig, itemConfig := agent.FindItem(taskId)
		if itemConfig == nil {
			logs.Println("task or item not found")
		} else {
			err := itemConfig.Validate()
			if err != nil {
				logs.Println("error:" + err.Error())
			} else {
				item := NewItem(appConfig.Id, itemConfig)
				v, err := item.Run()
				if err != nil {
					logs.Println("error:" + err.Error())
				} else {
					logs.Println("value:", v)
				}
			}
		}
		return
	}

	task := NewTask(appConfig.Id, taskConfig)
	_, stdout, stderr, err := task.Run()
	if len(stdout) > 0 {
		logs.Println("stdout:", stdout)
	}
	if len(stderr) > 0 {
		logs.Println("stderr:", stderr)
	}
	if err != nil {
		logs.Println(err.Error())
	}
}
