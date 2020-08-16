package agents

import (
	"github.com/TeaWeb/build/internal/teaconfigs/forms"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
	"github.com/iwind/TeaGo/maps"
	"github.com/shirou/gopsutil/net"
)

// 网络连接数
type ConnectionsSource struct {
	Source `yaml:",inline"`
}

// 获取新对象
func NewConnectionsSource() *ConnectionsSource {
	return &ConnectionsSource{}
}

// 名称
func (this *ConnectionsSource) Name() string {
	return "网络连接数"
}

// 代号
func (this *ConnectionsSource) Code() string {
	return "connections"
}

// 描述
func (this *ConnectionsSource) Description() string {
	return "获取网络连接数"
}

// 执行
func (this *ConnectionsSource) Execute(params map[string]string) (value interface{}, err error) {
	stat, err := net.Connections("all")
	if err != nil {
		return maps.Map{
			"connections": 0,
		}, err
	}

	value = maps.Map{
		"connections": len(stat),
	}

	return
}

// 表单信息
func (this *ConnectionsSource) Form() *forms.Form {
	form := forms.NewForm(this.Code())
	return form
}

// 变量
func (this *ConnectionsSource) Variables() []*SourceVariable {
	return []*SourceVariable{
		{
			Code:        "connections",
			Description: "连接数",
		},
	}
}

// 阈值
func (this *ConnectionsSource) Thresholds() []*Threshold {
	result := []*Threshold{}

	// 阈值
	{
		t := NewThreshold()
		t.Param = "${connections}"
		t.Operator = ThresholdOperatorGte
		t.Value = "10000"
		t.NoticeLevel = notices.NoticeLevelWarning
		t.NoticeMessage = "当前连接数过多"
		result = append(result, t)
	}

	return result
}

// 图表
func (this *ConnectionsSource) Charts() []*widgets.Chart {
	charts := []*widgets.Chart{}

	{
		// chart
		chart := widgets.NewChart()
		chart.Name = "网络连接数"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.SupportsTimeRange = true
		chart.Options = maps.Map{
			"code": `var chart = new charts.LineChart();

var ones = NewQuery().past(60, time.MINUTE).avg("connections");

var line = new charts.Line();
line.isFilled = true;

ones.$each(function (k, v) {
	line.addValue(v.value.connections);
	chart.addLabel(v.label);
});

chart.addLine(line);
chart.render();

`,
		}

		charts = append(charts, chart)
	}

	return charts
}

// 显示信息
func (this *ConnectionsSource) Presentation() *forms.Presentation {
	return nil
}
