package teaagents

import (
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/timers"
	"time"
)

// 监控项定义
type Item struct {
	appId     string
	config    *agents.Item
	lastTimer *time.Ticker
	oldValue  interface{}
}

// 获取新任务
func NewItem(appId string, config *agents.Item) *Item {
	return &Item{
		appId:  appId,
		config: config,
	}
}

// 运行一次
func (this *Item) Run() (value interface{}, err error) {
	source := this.config.Source()
	if source == nil {
		errMsg := "item " + this.config.Name + " source '" + this.config.SourceCode + "' does not exist, please update this agent to latest version"
		PushEvent(NewItemEvent(runningAgent.Id, this.appId, this.config.Id, "", errors.New(errMsg), time.Now().Unix(), 0))
		return "", errors.New(errMsg)
	}
	return source.Execute(nil)
}

// 定时运行
func (this *Item) Schedule() {
	if this.lastTimer != nil {
		this.lastTimer.Stop()
	}
	this.lastTimer = timers.Every(this.config.IntervalDuration(), func(ticker *time.Ticker) {
		source := this.config.Source()
		if source == nil {
			errMsg := "item " + this.config.Name + " source '" + this.config.SourceCode + "' does not exist, please update this agent to latest version"
			logs.Println(errMsg)
			PushEvent(NewItemEvent(runningAgent.Id, this.appId, this.config.Id, "", errors.New(errMsg), time.Now().Unix(), 0))
			return
		}

		t := time.Now()
		value, err := source.Execute(nil)
		costMs := time.Since(t).Seconds() * 1000
		if err != nil {
			logs.Println("execute " + this.config.Name + " error:" + err.Error())
		} else {
			this.oldValue = value
		}

		// 执行动作
		for _, threshold := range this.config.Thresholds {
			if len(threshold.Actions) == 0 {
				continue
			}
			b, err1 := threshold.Test(value, this.oldValue)
			if err1 != nil {
				logs.Println(this.config.Name + " error:" + err1.Error())
				if err == nil {
					err = err1
				}
			}
			if b {
				logs.Println("run " + this.config.Name + " [" + threshold.Param + " " + threshold.Operator + " " + threshold.Value + "] actions")
				err1 := threshold.RunActions(map[string]string{})
				if err1 != nil {
					logs.Println(this.config.Name + " error:" + err1.Error())

					if err == nil {
						err = err1
					}
				}
			}
		}

		if value != nil {
			PushEvent(NewItemEvent(runningAgent.Id, this.appId, this.config.Id, value, err, t.Unix(), costMs))
		} else {
			PushEvent(NewItemEvent(runningAgent.Id, this.appId, this.config.Id, "", err, t.Unix(), costMs))
		}
	})
}

// 停止运行
func (this *Item) Stop() {
	source := this.config.Source()
	if source != nil {
		source.Stop()
	}

	if this.lastTimer != nil {
		this.lastTimer.Stop()
	}
}
