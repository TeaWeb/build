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

// DNS检查
type DNSSource struct {
	Source `yaml:",inline"`

	Domain string `yaml:"domain" json:"domain"`
	Type   string `yaml:"type" json:"type"` // A, MX等
}

// 获取新对象
func NewDNSSource() *DNSSource {
	return &DNSSource{}
}

// 名称
func (this *DNSSource) Name() string {
	return "DNS解析"
}

// 代号
func (this *DNSSource) Code() string {
	return "dns"
}

// 描述
func (this *DNSSource) Description() string {
	return "使用DNS解析域名信息"
}

// 执行
func (this *DNSSource) Execute(params map[string]string) (value interface{}, err error) {
	if len(this.Domain) == 0 {
		return nil, errors.New("'domain' should not be empty")
	}
	switch this.Type {
	case "A", "AAAA":
		before := time.Now()
		ipList, err := net.LookupIP(this.Domain)
		if err != nil {
			return nil, err
		}

		result := []string{}
		for _, ip := range ipList {
			result = append(result, ip.String())
		}

		value = maps.Map{
			"result": result,
			"cost":   time.Since(before).Seconds(),
		}
	case "CHANGE":
		before := time.Now()
		result, err := net.LookupCNAME(this.Domain)
		if err != nil {
			return nil, err
		}
		value = maps.Map{
			"result": result,
			"cost":   time.Since(before).Seconds(),
		}
	case "MX":
		before := time.Now()
		mxList, err := net.LookupMX(this.Domain)
		if err != nil {
			return nil, err
		}

		result := []maps.Map{}
		for _, mx := range mxList {
			result = append(result, maps.Map{
				"host": mx.Host,
				"pref": mx.Pref,
			})
		}

		value = maps.Map{
			"result": result,
			"cost":   time.Since(before).Seconds(),
		}
	case "NS":
		before := time.Now()
		nxList, err := net.LookupNS(this.Domain)
		if err != nil {
			return nil, err
		}
		result := []maps.Map{}
		for _, ns := range nxList {
			result = append(result, maps.Map{
				"host": ns.Host,
			})
		}
		value = maps.Map{
			"result": result,
			"cost":   time.Since(before).Seconds(),
		}
	case "TXT":
		before := time.Now()
		result, err := net.LookupTXT(this.Domain)
		if err != nil {
			return nil, err
		}
		value = maps.Map{
			"result": result,
			"cost":   time.Since(before).Seconds(),
		}
	default:
		before := time.Now()

		ipList, err := net.LookupIP(this.Domain)
		if err != nil {
			return nil, err
		}

		result := []string{}
		for _, ip := range ipList {
			result = append(result, ip.String())
		}

		value = maps.Map{
			"result": result,
			"cost":   time.Since(before).Seconds(),
		}
	}
	return
}

// 表单信息
func (this *DNSSource) Form() *forms.Form {
	form := forms.NewForm(this.Code())
	{
		group := form.NewGroup()

		{
			field := forms.NewTextField("域名", "Domain")
			field.IsRequired = true
			field.Code = "domain"
			field.MaxLength = 500
			group.Add(field)
		}

		{
			field := forms.NewOptions("记录类型", "Type")
			field.IsRequired = true
			field.Attr("style", "width:10em")
			field.Code = "type"
			field.AddOption("A", "A")
			field.AddOption("AAAA", "AAAA")
			field.AddOption("CHANGE", "CHANGE")
			field.AddOption("MX", "MX")
			field.AddOption("TXT", "TXT")
			field.AddOption("NS", "NS")
			group.Add(field)
		}
	}

	return form
}

// 变量
func (this *DNSSource) Variables() []*SourceVariable {
	return []*SourceVariable{
		{
			Code:        "cost",
			Description: "解析耗时（秒）",
		},
		{
			Code:        "result",
			Description: "解析结果",
		},
	}
}

// 阈值
func (this *DNSSource) Thresholds() []*Threshold {
	result := []*Threshold{}

	// 阈值
	{
		t := NewThreshold()
		t.Param = "${result}"
		t.Value = ""
		t.NoticeLevel = notices.NoticeLevelWarning
		t.Operator = ThresholdOperatorEq
		t.NoticeMessage = "无法解析当前域名"
		result = append(result, t)
	}

	{
		t := NewThreshold()
		t.Param = "${cost}"
		t.Value = "5"
		t.NoticeLevel = notices.NoticeLevelWarning
		t.Operator = ThresholdOperatorGte
		t.NoticeMessage = "解析耗时过长"
		result = append(result, t)
	}

	return result
}

// 图表
func (this *DNSSource) Charts() []*widgets.Chart {
	charts := []*widgets.Chart{}

	{
		// chart
		chart := widgets.NewChart()
		chart.Name = "DNS解析耗时（ms）"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.SupportsTimeRange = true
		chart.Options = maps.Map{
			"code": `var chart = new charts.LineChart();

var ones = NewQuery().past(60, time.MINUTE).avg("cost");

var line = new charts.Line();
line.isFilled = true;

var maxValue = 0;

ones.$each(function (k, v) {
	var ms = v.value.cost * 1000;
	line.addValue(ms);
	if (maxValue < ms) {
		maxValue = ms;	
	}

	chart.addLabel(v.label);
});

if (maxValue < 10) {
	maxValue = 10;	
}
if (maxValue > 0) {
	chart.max = maxValue;	
}

chart.addLine(line);
chart.render();
`,
		}

		charts = append(charts, chart)
	}

	return charts
}

// 显示信息
func (this *DNSSource) Presentation() *forms.Presentation {
	p := forms.NewPresentation()
	p.HTML = `
<tr>
	<td>域名<em>（Domain）</em></td>
	<td>{{source.domain}}</td>
</tr>
<tr>
	<td>记录类型<em>（Type）</em></td>
	<td>{{source.type}}</td>
</tr>
`
	return p
}
