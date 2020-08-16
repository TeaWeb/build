package agents

import (
	"github.com/TeaWeb/build/internal/teaconfigs/forms"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/time"
	"time"
)

// 日期相关
type DateSource struct {
	Source `yaml:",inline"`
}

// 获取新对象
func NewDateSource() *DateSource {
	return &DateSource{}
}

// 名称
func (this *DateSource) Name() string {
	return "日期时间"
}

// 代号
func (this *DateSource) Code() string {
	return "date"
}

// 描述
func (this *DateSource) Description() string {
	return "获取主机上日期时间信息"
}

// 执行
func (this *DateSource) Execute(params map[string]string) (value interface{}, err error) {
	t := time.Now()
	value = maps.Map{
		"timestamp": t.Unix(),
		"nano":      t.Nanosecond(),
		"date":      timeutil.Format("Y-m-d H:i:s", t),
		"offset":    timeutil.Format("O", t),
	}
	return
}

// 表单信息
func (this *DateSource) Form() *forms.Form {
	form := forms.NewForm(this.Code())
	return form
}

// 变量
func (this *DateSource) Variables() []*SourceVariable {
	return []*SourceVariable{
		{
			Code:        "timestamp",
			Description: "时间戳（秒）",
		},
		{
			Code:        "nano",
			Description: "Nano时间",
		},
		{
			Code:        "date",
			Description: "日期形式的时间，格式为Y-m-d H:i:s",
		},
		{
			Code:        "offset",
			Description: "时差，比如+0800",
		},
	}
}

// 阈值
func (this *DateSource) Thresholds() []*Threshold {
	result := []*Threshold{}

	// 阈值
	{
		t := NewThreshold()
		t.Param = "new Date().getTime() / 1000 - ${timestamp}"
		t.Value = "300"
		t.NoticeLevel = notices.NoticeLevelWarning
		t.Operator = ThresholdOperatorGte
		t.NoticeMessage = "主机时间出现很大偏差"
		result = append(result, t)
	}

	return result
}

// 图表
func (this *DateSource) Charts() []*widgets.Chart {
	charts := []*widgets.Chart{}

	// 时钟
	{
		chart := widgets.NewChart()
		chart.Id = "clock"
		chart.Name = "时钟"
		chart.Columns = 1
		chart.Type = "javascript"
		chart.Options = maps.Map{
			"code": `var chart = new charts.Clock();
var latest = NewQuery().latest(1);
if (latest.length > 0) {
	chart.timestamp = parseInt(new Date().getTime() / 1000) - (latest[0].createdAt - latest[0].value.timestamp);
}
chart.render();
`,
		}
		charts = append(charts, chart)
	}

	return charts
}

// 显示信息
func (this *DateSource) Presentation() *forms.Presentation {
	return nil
}
