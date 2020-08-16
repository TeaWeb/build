package agents

import (
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs/forms"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
	"github.com/iwind/TeaGo/maps"
	"net"
	"time"
)

// Socket连通性
type SocketConnectivitySource struct {
	Source `yaml:",inline"`

	Address string `yaml:"address" json:"address"`
	Network string `yaml:"network" json:"network"`
}

// 获取新对象
func NewSocketConnectivitySource() *SocketConnectivitySource {
	return &SocketConnectivitySource{}
}

// 名称
func (this *SocketConnectivitySource) Name() string {
	return "端口连通性"
}

// 代号
func (this *SocketConnectivitySource) Code() string {
	return "socketConnectivity"
}

// 描述
func (this *SocketConnectivitySource) Description() string {
	return "获取端口连通性信息"
}

// 执行
func (this *SocketConnectivitySource) Execute(params map[string]string) (value interface{}, err error) {
	before := time.Now()

	if len(this.Address) == 0 {
		err = errors.New("'address' should not be empty")
		value = maps.Map{
			"cost": time.Since(before).Seconds(),
		}
		return
	}

	network := this.Network
	if len(network) == 0 {
		network = "tcp"
	}

	conn, err := net.Dial(network, this.Address)
	if err != nil {
		value = maps.Map{
			"cost": time.Since(before).Seconds(),
		}
		return value, err
	}

	_ = conn.Close()

	value = maps.Map{
		"cost":    time.Since(before).Seconds(),
		"success": 1,
	}

	return
}

// 表单信息
func (this *SocketConnectivitySource) Form() *forms.Form {
	form := forms.NewForm(this.Code())

	{
		group := form.NewGroup()

		{
			field := forms.NewTextField("目标地址", "")
			field.Code = "address"
			field.IsRequired = true
			field.Placeholder = ""
			field.Comment = "如果是tcp或者udp，地址需要加端口"
			field.MaxLength = 500
			field.ValidateCode = `
if (value.length == 0) {
	throw new Error("请输入目标地址")
}
`
			group.Add(field)
		}

		{
			field := forms.NewOptions("网络协议", "")
			field.Code = "network"
			field.IsRequired = true
			field.AddOption("TCP", "tcp")
			field.AddOption("UDP", "udp")
			field.AddOption("Unix", "unix")
			field.Attr("style", "width:10em")
			field.ValidateCode = `
if (value.length == 0) {
	throw new Error("请选择网络协议");
}
`
			group.Add(field)
		}
	}

	return form
}

// 变量
func (this *SocketConnectivitySource) Variables() []*SourceVariable {
	return []*SourceVariable{
		{
			Code:        "cost",
			Description: "连接耗时（秒）",
		},
	}
}

// 阈值
func (this *SocketConnectivitySource) Thresholds() []*Threshold {
	result := []*Threshold{}

	// 阈值
	{
		t := NewThreshold()
		t.Param = "${success}"
		t.Operator = ThresholdOperatorNot
		t.Value = "1"
		t.NoticeLevel = notices.NoticeLevelWarning
		t.NoticeMessage = "端口连接失败"
		result = append(result, t)
	}

	return result
}

// 图表
func (this *SocketConnectivitySource) Charts() []*widgets.Chart {
	charts := []*widgets.Chart{}

	{
		// chart
		chart := widgets.NewChart()
		chart.Name = "端口连通性（ms）"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.SupportsTimeRange = true
		chart.Options = maps.Map{
			"code": `var chart = new charts.LineChart();

var ones = NewQuery().past(60, time.MINUTE).avg("cost");

var line = new charts.Line();
line.isFilled = true;

ones.$each(function (k, v) {
	if (v.value == "") {
		return;
	}
	line.addValue(v.value.cost * 1000);
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
func (this *SocketConnectivitySource) Presentation() *forms.Presentation {
	p := forms.NewPresentation()
	p.HTML = `
<tr>
	<td>目标地址</td>
	<td>{{source.address}}</td>
</tr>
<tr>
	<td>网络协议</td>
	<td>{{source.network.toUpperCase()}}</td>
</tr>
`
	return p
}
