package agents

import (
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
)

// App定义
type AppConfig struct {
	Id    string        `yaml:"id" json:"id"`       // ID
	On    bool          `yaml:"on" json:"on"`       // 是否启用
	Tasks []*TaskConfig `yaml:"tasks" json:"tasks"` // 任务设置
	Items []*Item       `yaml:"item" json:"items"`  // 监控项
	Name  string        `yaml:"name" json:"name"`   // 名称

	IsSharedWithGroup bool     `yaml:"issharedwithgroup" json:"isSharedWithGroup"` // 是否与当前组共享，使用issharedwithgroup是为了兼容v0.1.6之前的版本
	SharedAgentIds    []string `yaml:"sharedAgentIds" json:"sharedAgentIds"`       // 共享的Agents TODO 暂不实现

	NoticeSetting map[notices.NoticeLevel][]*notices.NoticeReceiver `yaml:"noticeSetting" json:"noticeSetting"`
}

// 获取新对象
func NewAppConfig() *AppConfig {
	return &AppConfig{
		Id:            rands.HexString(16),
		On:            true,
		NoticeSetting: map[notices.NoticeLevel][]*notices.NoticeReceiver{},
	}
}

// 校验
func (this *AppConfig) Validate() error {
	// 任务
	for _, t := range this.Tasks {
		err := t.Validate()
		if err != nil {
			return err
		}
	}

	// 监控项
	for _, item := range this.Items {
		err := item.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// Schedule Tasks
func (this *AppConfig) FindSchedulingTasks() []*TaskConfig {
	result := []*TaskConfig{}
	for _, t := range this.Tasks {
		if len(t.Schedule) > 0 {
			result = append(result, t)
		}
	}
	return result
}

// Boot Tasks
func (this *AppConfig) FindBootingTasks() []*TaskConfig {
	result := []*TaskConfig{}
	for _, t := range this.Tasks {
		if t.IsBooting {
			result = append(result, t)
		}
	}
	return result
}

// Manual Tasks
func (this *AppConfig) FindManualTasks() []*TaskConfig {
	result := []*TaskConfig{}
	for _, t := range this.Tasks {
		if t.IsManual {
			result = append(result, t)
		}
	}
	return result
}

// 添加任务
func (this *AppConfig) AddTask(task *TaskConfig) {
	this.Tasks = append(this.Tasks, task)
}

// 删除任务
func (this *AppConfig) RemoveTask(taskId string) {
	result := []*TaskConfig{}
	for _, t := range this.Tasks {
		if t.Id == taskId {
			continue
		}
		result = append(result, t)
	}
	this.Tasks = result
}

// 查找任务
func (this *AppConfig) FindTask(taskId string) *TaskConfig {
	for _, t := range this.Tasks {
		if t.Id == taskId {
			return t
		}
	}
	return nil
}

// 添加监控项
func (this *AppConfig) AddItem(item *Item) {
	this.Items = append(this.Items, item)
}

// 删除监控项
func (this *AppConfig) RemoveItem(itemId string) {
	result := []*Item{}
	for _, item := range this.Items {
		if item.Id == itemId {
			continue
		}
		result = append(result, item)
	}
	this.Items = result
}

// 查找监控项
func (this *AppConfig) FindItem(itemId string) *Item {
	for _, item := range this.Items {
		if item.Id == itemId {
			item.Validate()
			return item
		}
	}
	return nil
}

// 添加通知接收者
func (this *AppConfig) AddNoticeReceiver(level notices.NoticeLevel, receiver *notices.NoticeReceiver) {
	if this.NoticeSetting == nil {
		this.NoticeSetting = map[notices.NoticeLevel][]*notices.NoticeReceiver{}
	}
	receivers, found := this.NoticeSetting[level]
	if !found {
		receivers = []*notices.NoticeReceiver{}
	}
	receivers = append(receivers, receiver)
	this.NoticeSetting[level] = receivers
}

// 删除通知接收者
func (this *AppConfig) RemoveNoticeReceiver(level notices.NoticeLevel, receiverId string) {
	if this.NoticeSetting == nil {
		return
	}
	receivers, found := this.NoticeSetting[level]
	if !found {
		return
	}

	result := []*notices.NoticeReceiver{}
	for _, r := range receivers {
		if r.Id == receiverId {
			continue
		}
		result = append(result, r)
	}
	this.NoticeSetting[level] = result
}

// 获取通知接收者数量
func (this *AppConfig) CountNoticeReceivers() int {
	count := 0
	for _, receivers := range this.NoticeSetting {
		count += len(receivers)
	}
	return count
}

// 删除媒介
func (this *AppConfig) RemoveMedia(mediaId string) (found bool) {
	for level, receivers := range this.NoticeSetting {
		result := []*notices.NoticeReceiver{}
		for _, receiver := range receivers {
			if receiver.MediaId == mediaId {
				found = true
				continue
			}
			result = append(result, receiver)
		}
		this.NoticeSetting[level] = result
	}
	return
}

// 查找一个或多个级别对应的接收者，并合并相同的接收者
func (this *AppConfig) FindAllNoticeReceivers(level ...notices.NoticeLevel) []*notices.NoticeReceiver {
	if len(level) == 0 {
		return []*notices.NoticeReceiver{}
	}

	m := maps.Map{} // mediaId_user => bool
	result := []*notices.NoticeReceiver{}
	for _, l := range level {
		receivers, ok := this.NoticeSetting[l]
		if !ok {
			continue
		}
		for _, receiver := range receivers {
			if !receiver.On {
				continue
			}
			key := receiver.Key()
			if m.Has(key) {
				continue
			}
			m[key] = true
			result = append(result, receiver)
		}
	}
	return result
}

// 匹配关键词
func (this *AppConfig) MatchKeyword(keyword string) (matched bool, name string, tags []string) {
	if teautils.MatchKeyword(this.Name, keyword) {
		matched = true
		name = this.Name
		return
	}

	return
}
