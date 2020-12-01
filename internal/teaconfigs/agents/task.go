package agents

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/rands"
	"runtime"
	"strings"
	"time"
)

// 任务
// 日志存储在 task.${id}
type TaskConfig struct {
	Id        string             `yaml:"id" json:"id"`               // ID
	On        bool               `yaml:"on" json:"on"`               // 是否启用
	Name      string             `yaml:"name" json:"name"`           // 名称
	Cwd       string             `yaml:"cwd" json:"cwd"`             // 当前工作目录（Current Working Directory）
	Env       []*shared.Variable `yaml:"env" json:"env"`             // 环境变量设置
	Script    string             `yaml:"script" json:"script"`       // 脚本
	IsBooting bool               `yaml:"isBooting" json:"isBooting"` // 在Boot时启动
	Schedule  []*ScheduleConfig  `yaml:"schedule" json:"schedule"`   // 定时
	IsManual  bool               `yaml:"isManual" json:"isManual"`   // 是否手工调用
	Version   uint               `yaml:"version" json:"version"`     // 版本
}

// 获取新对象
func NewTaskConfig() *TaskConfig {
	return &TaskConfig{
		On: true,
		Id: rands.HexString(16),
	}
}

// 校验
func (this *TaskConfig) Validate() error {
	for _, s := range this.Schedule {
		err := s.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

// 添加环境变量
func (this *TaskConfig) AddEnv(name string, value string) {
	this.Env = append(this.Env, &shared.Variable{
		Name:  name,
		Value: value,
	})
}

// 添加定时
func (this *TaskConfig) AddSchedule(schedule *ScheduleConfig) {
	this.Schedule = append(this.Schedule, schedule)
}

// 取得下次运行时间
func (this *TaskConfig) Next(now time.Time) (next time.Time, ok bool) {
	if len(this.Schedule) == 0 {
		return
	}
	last := time.Date(0, 0, 0, 0, 0, 0, 0, time.Local)
	for _, s := range this.Schedule {
		next, ok := s.Next(now)
		if ok && (last.Year() <= 0 || next.Before(last)) {
			last = next
		}
	}

	if last.Unix() >= now.Unix() {
		next = last
		ok = true
		return
	}
	return
}

// 格式化脚本
func (this *TaskConfig) FormattedScript() string {
	script := this.Script
	script = strings.Replace(script, "\r", "", -1)
	return script
}

// 保存到本地
func (this *TaskConfig) Generate() (path string, err error) {
	if runtime.GOOS == "windows" {
		path = Tea.ConfigFile("agents/task." + this.Id + ".bat")
	} else {
		path = Tea.ConfigFile("agents/task." + this.Id + ".script")
	}
	shFile := files.NewFile(path)
	if !shFile.Exists() {
		err = shFile.WriteString(this.FormattedScript())
		if err != nil {
			return
		}
		err = shFile.Chmod(0777)
		if err != nil {
			return
		}
	}
	return
}

// 重新生成
func (this *TaskConfig) GenerateAgain() (path string, err error) {
	if runtime.GOOS == "windows" {
		path = Tea.ConfigFile("agents/task." + this.Id + ".bat")
	} else {
		path = Tea.ConfigFile("agents/task." + this.Id + ".script")
	}
	shFile := files.NewFile(path)
	err = shFile.WriteString(this.FormattedScript())
	if err != nil {
		return
	}
	err = shFile.Chmod(0777)
	if err != nil {
		return
	}
	return
}

// 匹配关键词
func (this *TaskConfig) MatchKeyword(keyword string) (matched bool, name string, tags []string) {
	if teautils.MatchKeyword(this.Name, keyword) {
		matched = true
		name = this.Name
	}
	return
}
