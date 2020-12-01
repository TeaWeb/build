package api

import (
	"errors"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/rands"
	"time"
)

// API测试计划
type APITestPlan struct {
	On       bool     `yaml:"on" json:"on"`             // 是否开启
	Filename string   `yaml:"filename" json:"filename"` // 配置文件名
	Hour     int      `yaml:"hour" json:"hour"`         // 小时
	Minute   int      `yaml:"minute" json:"minute"`     // 分钟
	Second   int      `yaml:"second" json:"second"`     // 秒
	Weekdays []int    `yaml:"weekdays" json:"weekdays"` // 周
	Reports  []string `yaml:"reports" json:"reports"`   // 报告文件名
	APIs     []string `yaml:"apis" json:"apis"`         // 参与计划的API TODO 需要实现
}

// 获取新对象
func NewAPITestPlan() *APITestPlan {
	return &APITestPlan{
		On: true,
	}
}

// 从文件中加载测试计划
func NewAPITestPlanFromFile(filename string) *APITestPlan {
	if len(filename) == 0 {
		return nil
	}
	file := files.NewFile(Tea.ConfigFile(filename))
	reader, err := file.Reader()
	if err != nil {
		logs.Error(err)
		return nil
	}
	defer reader.Close()

	plan := NewAPITestPlan()
	err = reader.ReadYAML(plan)
	if err != nil {
		logs.Error(err)
		return nil
	}

	plan.Filename = filename
	return plan
}

// 保存当前测试计划
func (this *APITestPlan) Save() error {
	if len(this.Filename) == 0 {
		this.Filename = "plan." + rands.HexString(16)+ ".conf"
	}
	file := files.NewFile(Tea.ConfigFile(this.Filename))
	writer, err := file.Writer()
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = writer.WriteYAML(this)
	return nil
}

// 删除当前测试计划
func (this *APITestPlan) Delete() error {
	if len(this.Filename) == 0 {
		return errors.New("filename should not be empty")
	}
	file := files.NewFile(Tea.ConfigFile(this.Filename))
	err := file.Delete()
	if err == nil {
		// 删除报告
		for _, report := range this.Reports {
			err := files.NewFile(Tea.ConfigFile(report)).Delete()
			if err != nil {
				logs.Error(err)
			}
		}
	}
	return err
}

// 添加测试报告
func (this *APITestPlan) AddReport(reportFilename string) {
	this.Reports = append(this.Reports, reportFilename)
}

// 读取最后一次报告
func (this *APITestPlan) LastReport() *APITestPlanReport {
	if len(this.Reports) == 0 {
		return nil
	}
	reportFile := this.Reports[len(this.Reports)-1]
	return NewAPITestPlanReportFromFile(reportFile)
}

// 取得周内日期
func (this *APITestPlan) WeekdayNames() []string {
	result := []string{}
	for _, weekday := range this.Weekdays {
		switch weekday {
		case 1:
			result = append(result, "周一")
		case 2:
			result = append(result, "周二")
		case 3:
			result = append(result, "周三")
		case 4:
			result = append(result, "周四")
		case 5:
			result = append(result, "周五")
		case 6:
			result = append(result, "周六")
		case 7:
			result = append(result, "周日")
		}
	}
	return result
}

// 检查时间是否匹配
func (this *APITestPlan) MatchTime(currentTime time.Time) bool {
	weekday := int(currentTime.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return currentTime.Hour() == this.Hour &&
		currentTime.Minute() == this.Minute &&
		currentTime.Second() == this.Second &&
		lists.Contains(this.Weekdays, weekday)
}
