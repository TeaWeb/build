package api

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
	"github.com/iwind/TeaGo/utils/time"
	"sync"
	"time"
)

// API测试报告
type APITestPlanReport struct {
	Filename     string           `yaml:"filename" json:"filename"`         // 文件名
	StartedAt    int64            `yaml:"startedAt" json:"startedAt"`       // 开始时间
	FinishedAt   int64            `yaml:"finishedAt" json:"finishedAt"`     // 结束时间
	TotalAPIs    int              `yaml:"totalApis" json:"totalApis"`       // 总体API数量
	TotalScripts int              `yaml:"totalScripts" json:"totalScripts"` // 总体脚本总数
	Results      []*APITestResult `yaml:"results" json:"results"`           // 统计

	locker sync.Mutex // 操作锁
}

// 获取新对象
func NewAPITestPlanReport() *APITestPlanReport {
	return &APITestPlanReport{}
}

// 从配置文件中加载测试报告
func NewAPITestPlanReportFromFile(filename string) *APITestPlanReport {
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

	report := NewAPITestPlanReport()
	err = reader.ReadYAML(report)
	if err != nil {
		return nil
	}
	return report
}

// 初始化文件信息
func (this *APITestPlanReport) InitFile() {
	this.Filename = "report." + rands.HexString(16)+ ".conf"
}

// 计算总结果数
func (this *APITestPlanReport) CountResults() int {
	this.locker.Lock()
	defer this.locker.Unlock()

	return len(this.Results)
}

// 计算失败的结果数
func (this *APITestPlanReport) CountFailedResults() int {
	this.locker.Lock()
	defer this.locker.Unlock()

	i := 0
	for _, result := range this.Results {
		if !result.IsPassed {
			i ++
		}
	}
	return i
}

// 计算脚本数
func (this *APITestPlanReport) CountScripts() int {
	this.locker.Lock()
	defer this.locker.Unlock()

	count := 0
	for _, result := range this.Results {
		count += len(result.Scripts)
	}
	return count
}

// 添加API执行结果
func (this *APITestPlanReport) AddAPIResult(apiResult *APITestResult) {
	this.locker.Lock()
	defer this.locker.Unlock()

	this.Results = append(this.Results, apiResult)
}

// 保存
func (this *APITestPlanReport) Save() error {
	this.locker.Lock()
	defer this.locker.Unlock()

	if len(this.Filename) == 0 {
		this.InitFile()
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

// 取得报告的综合信息
func (this *APITestPlanReport) Summary() maps.Map {
	countApis := this.CountResults()
	countScripts := this.CountScripts()

	reportMap := maps.Map{
		"isStarted":            this.StartedAt > 0,
		"isFinished":           this.FinishedAt > 0,
		"totalApis":            this.TotalAPIs,
		"totalScripts":         this.TotalScripts,
		"countFinishedApis":    countApis,
		"countFinishedScripts": countScripts,
	}

	if this.StartedAt > 0 {
		reportMap["startedTime"] = timeutil.Format("Y-m-d H:i:s", time.Unix(this.StartedAt, 0))
	} else {
		reportMap["startedTime"] = ""
	}

	if this.FinishedAt > 0 {
		reportMap["finishedTime"] = timeutil.Format("Y-m-d H:i:s", time.Unix(this.FinishedAt, 0))
	} else {
		reportMap["finishedTime"] = ""
	}

	if this.TotalAPIs > 0 && this.TotalScripts > 0 {
		reportMap["percentApis"] = countApis * 100 / this.TotalAPIs
		reportMap["percentScripts"] = countScripts * 100 / this.TotalScripts
	} else {
		reportMap["percentApis"] = 0
		reportMap["percentScripts"] = 0
	}

	return reportMap
}
