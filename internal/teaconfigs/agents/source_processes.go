package agents

import (
	"github.com/TeaWeb/build/internal/teaconfigs/forms"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
	"github.com/iwind/TeaGo/maps"
	"github.com/shirou/gopsutil/process"
)

// 进程数
type ProcessesSource struct {
	Source `yaml:",inline"`
}

// 获取新对象
func NewProcessesSource() *ProcessesSource {
	return &ProcessesSource{}
}

// 名称
func (this *ProcessesSource) Name() string {
	return "进程数"
}

// 代号
func (this *ProcessesSource) Code() string {
	return "processes"
}

// 描述
func (this *ProcessesSource) Description() string {
	return "获取当前主机运行的进程数"
}

// 执行
func (this *ProcessesSource) Execute(params map[string]string) (value interface{}, err error) {
	stat, err := process.Pids()
	if err != nil {
		return maps.Map{
			"processes": 0,
		}, err
	}

	value = maps.Map{
		"processes": len(stat),
	}

	return
}

// 表单信息
func (this *ProcessesSource) Form() *forms.Form {
	form := forms.NewForm(this.Code())
	return form
}

// 变量
func (this *ProcessesSource) Variables() []*SourceVariable {
	return []*SourceVariable{
		{
			Code:        "processes",
			Description: "进程数",
		},
	}
}

// 阈值
func (this *ProcessesSource) Thresholds() []*Threshold {
	result := []*Threshold{}

	// 阈值
	{
		t := NewThreshold()
		t.Param = "${processes}"
		t.Operator = ThresholdOperatorGte
		t.Value = "1000"
		t.NoticeLevel = notices.NoticeLevelWarning
		t.NoticeMessage = "当前系统进程数过多"
		result = append(result, t)
	}

	return result
}

// 图表
func (this *ProcessesSource) Charts() []*widgets.Chart {
	charts := []*widgets.Chart{}

	{
		// chart
		chart := widgets.NewChart()
		chart.Name = "系统进程数"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.SupportsTimeRange = true
		chart.Options = maps.Map{
			"code": `var chart = new charts.LineChart();

var ones = NewQuery().past(60, time.MINUTE).avg("processes");

var line = new charts.Line();
line.isFilled = true;

ones.$each(function (k, v) {
	line.addValue(v.value.processes);
	chart.addLabel(v.label);
});

chart.addLine(line);
chart.render();`,
		}

		charts = append(charts, chart)
	}

	return charts
}

// 显示信息
func (this *ProcessesSource) Presentation() *forms.Presentation {
	return nil
}
