package agents

import (
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/rands"
	"time"
)

// 数据指标项
type Item struct {
	On               bool                   `yaml:"on" json:"on"`
	Id               string                 `yaml:"id" json:"id"`
	Name             string                 `yaml:"name" json:"name"`
	SourceCode       string                 `yaml:"sourceCode" json:"sourceCode"`             // 数据源代号
	SourceOptions    map[string]interface{} `yaml:"sourceOptions" json:"sourceOptions"`       // 数据源选项
	Interval         string                 `yaml:"interval" json:"interval"`                 // 刷新间隔
	Thresholds       []*Threshold           `yaml:"thresholds" json:"thresholds"`             // 阈值设置
	Charts           []*widgets.Chart       `yaml:"charts" json:"charts"`                     // 图表
	RecoverSuccesses int                    `yaml:"recoverSuccesses" json:"recoverSuccesses"` // 恢复的成功次数

	source           SourceInterface
	intervalDuration time.Duration
}

// 获取新对象
func NewItem() *Item {
	return &Item{
		On:         true,
		Id:         rands.HexString(16),
		Thresholds: []*Threshold{},
	}
}

// 校验
func (this *Item) Validate() error {
	this.source = FindDataSourceInstance(this.SourceCode, this.SourceOptions)

	if len(this.Charts) == 0 {
		this.Charts = []*widgets.Chart{}
	}

	if this.source != nil {
		err := this.source.Validate()
		if err != nil {
			return err
		}
	}

	for _, t := range this.Thresholds {
		err := t.Validate()
		if err != nil {
			return err
		}
	}

	for _, c := range this.Charts {
		err := c.Validate()
		if err != nil {
			return err
		}
	}

	this.intervalDuration, _ = time.ParseDuration(this.Interval)

	return nil
}

// 获取刷新间隔
func (this *Item) IntervalDuration() time.Duration {
	if this.intervalDuration.Seconds() > 0 {
		return this.intervalDuration
	}
	return 30 * time.Second
}

// 添加阈值
func (this *Item) AddThreshold(t ...*Threshold) {
	this.Thresholds = append(this.Thresholds, t...)
}

// 数据源对象
func (this *Item) Source() SourceInterface {
	return this.source
}

// 检查某个值对应的通知级别
func (this *Item) TestValue(value interface{}, oldValue interface{}) (threshold *Threshold, row interface{}, level notices.NoticeLevel, message string, err error) {
	if len(this.Thresholds) == 0 {
		return nil, nil, notices.NoticeLevelNone, "", nil
	}
	for _, t := range this.Thresholds {
		b, row, testErr := t.TestRow(value, oldValue)
		if testErr != nil {
			return nil, nil, notices.NoticeLevelNone, "", errors.New("[threshold] " + testErr.Error())
		}
		if b {
			if len(t.NoticeMessage) > 0 {
				return t, row, t.NoticeLevel, t.NoticeMessage, nil
			} else {
				return t, row, t.NoticeLevel, "", nil
			}
		}
	}
	return nil, nil, notices.NoticeLevelNone, "", nil
}

// 添加图表
func (this *Item) AddChart(chart *widgets.Chart) {
	this.Charts = append(this.Charts, chart)
}

// 添加一组图表中的某几个
func (this *Item) AddFilterCharts(charts []*widgets.Chart, chartId ...string) {
	for _, c := range charts {
		if lists.ContainsString(chartId, c.Id) {
			this.AddChart(c)
		}
	}
}

// 查找图表
func (this *Item) FindChart(chartId string) *widgets.Chart {
	for _, c := range this.Charts {
		if c.Id == chartId {
			return c
		}
	}
	return nil
}

// 删除图表
func (this *Item) RemoveChart(chartId string) {
	result := []*widgets.Chart{}
	for _, c := range this.Charts {
		if c.Id == chartId {
			continue
		}
		result = append(result, c)
	}
	this.Charts = result
}

// 匹配关键词
func (this *Item) MatchKeyword(keyword string) (matched bool, name string, tags []string) {
	if teautils.MatchKeyword(this.Name, keyword) || teautils.MatchKeyword(this.SourceCode, keyword) {
		matched = true
		name = this.Name
		if len(this.SourceCode) > 0 {
			tags = []string{"数据源：" + this.SourceCode}
		}
	}
	return
}
