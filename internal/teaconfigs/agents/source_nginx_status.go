package agents

import (
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs/forms"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"time"
)

// Nginx Status数据源
type NginxStatusSource struct {
	Source `yaml:",inline"`

	URL string `yaml:"url" json:"url"`

	lastRequests int64
	lastTime     time.Time
}

// 获取新对象
func NewNginxStatusSource() *NginxStatusSource {
	return &NginxStatusSource{}
}

// 校验
func (this *NginxStatusSource) Validate() error {
	return nil
}

// 名称
func (this *NginxStatusSource) Name() string {
	return "Nginx Status"
}

// 代号
func (this *NginxStatusSource) Code() string {
	return "nginxStatus"
}

// 描述
func (this *NginxStatusSource) Description() string {
	return "利用ngx_http_stub_status_module模块读取Nginx相关信息"
}

// 执行
func (this *NginxStatusSource) Execute(params map[string]string) (value interface{}, err error) {
	before := time.Now()

	if len(this.URL) == 0 {
		return maps.Map{
			"cost": time.Since(before).Seconds(),
		}, errors.New("'url' should not be empty")
	}

	timeout := 10 * time.Second
	req, err := http.NewRequest(http.MethodGet, this.URL, nil)
	if err != nil {
		return maps.Map{
			"cost": time.Since(before).Seconds(),
		}, err
	}
	req.Header.Set("User-Agent", teaconst.TeaProductCode+"/"+teaconst.TeaVersion)

	client := teautils.SharedHttpClient(timeout)
	resp, err := client.Do(req)
	if err != nil {
		return maps.Map{
			"cost": time.Since(before).Seconds(),
		}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return maps.Map{
			"cost":   time.Since(before).Seconds(),
			"status": resp.StatusCode,
		}, errors.New("'status' should be 200")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return maps.Map{
			"cost":   time.Since(before).Seconds(),
			"status": resp.StatusCode,
		}, err
	}

	activeConnections := 0
	result := regexp.MustCompile("Active\\s+connections:\\s+(\\d+)").FindSubmatch(data)
	if len(result) > 1 {
		activeConnections = types.Int(result[1])
	}

	acceptedConnections := int64(0)
	handledConnections := int64(0)
	totalRequests := int64(0)
	result = regexp.MustCompile("server\\s+accepts\\s+handled\\s+requests\\s+(\\d+)\\s+(\\d+)\\s+(\\d+)").FindSubmatch(data)
	if len(result) > 1 {
		acceptedConnections = types.Int64(result[1])
		handledConnections = types.Int64(result[2])
		totalRequests = types.Int64(result[3])
	}

	reading := 0
	writing := 0
	waiting := 0
	result = regexp.MustCompile("Reading:\\s+(\\d+)\\s+Writing:\\s+(\\d+)\\s+Waiting:\\s+(\\d+)").FindSubmatch(data)
	if len(result) > 1 {
		reading = types.Int(result[1])
		writing = types.Int(result[2])
		waiting = types.Int(result[3])
	}

	// 平均每秒请求数
	requestsPerSecond := 0
	if totalRequests > this.lastRequests && this.lastRequests > 0 {
		requestsPerSecond = int(math.Ceil(float64(totalRequests-this.lastRequests) / time.Since(this.lastTime).Seconds()))
		if requestsPerSecond > 0 {
			requestsPerSecond-- // 减去监控系统的请求
		}
	}

	this.lastRequests = totalRequests
	this.lastTime = time.Now()

	return maps.Map{
		"status": resp.StatusCode,
		"cost":   time.Since(before).Seconds(),
		"result": maps.Map{
			"activeConnections":   activeConnections,
			"acceptedConnections": acceptedConnections,
			"handledConnections":  handledConnections,
			"totalRequests":       totalRequests,
			"readingConnections":  reading,
			"writingConnections":  writing,
			"waitingConnections":  waiting,
			"requestsPerSecond":   requestsPerSecond,
		},
	}, nil
}

// 选项表单
func (this *NginxStatusSource) Form() *forms.Form {
	form := forms.NewForm(this.Code())

	group := form.NewGroup()
	{
		field := forms.NewTextBox("Nginx状态检查URL", "")
		field.Rows = 2
		field.Comment = "<a href=\"http://nginx.org/en/docs/http/ngx_http_stub_status_module.html\" target=\"_blank\">相关配置文档&raquo;</a>"
		field.Code = "url"
		field.IsRequired = true
		field.MaxLength = 500
		field.Attr("style", "word-wrap:break-word;word-break:break-all;line-height:1.5")
		field.ValidateCode = `
if (value.length == 0) {
	throw new Error("状态检查URL")
}

if (!value.match(/^(http|https):\/\//i)) {
	throw new Error("URL地址必须以http或https开头");
}
`
		group.Add(field)
	}

	return form
}

// 变量
func (this *NginxStatusSource) Variables() []*SourceVariable {
	return []*SourceVariable{
		{
			Code:        "cost",
			Description: "请求耗时（秒）",
		},
		{
			Code:        "status",
			Description: "HTTP状态码",
		},
		{
			Code:        "result.activeConnections",
			Description: "当前活跃的连接",
		},
		{
			Code:        "result.acceptedConnections",
			Description: "已建立连接的连接数",
		},
		{
			Code:        "result.handledConnections",
			Description: "已处理的连接数",
		},
		{
			Code:        "result.totalRequests",
			Description: "所有请求数",
		},
		{
			Code:        "result.readingConnections",
			Description: "nginx正在读取的连接数",
		},
		{
			Code:        "result.writingConnections",
			Description: "nginx正在写入的连接数",
		},
		{
			Code:        "result.waitingConnections",
			Description: "正在等待请求的连接数",
		},
		{
			Code:        "result.requestsPerSecond",
			Description: "平均每秒请求数",
		},
	}
}

// 阈值
func (this *NginxStatusSource) Thresholds() []*Threshold {
	result := []*Threshold{}

	{
		t := NewThreshold()
		t.Param = "${status}"
		t.NoticeLevel = notices.NoticeLevelWarning
		t.Operator = ThresholdOperatorNot
		t.Value = "200"
		t.NoticeMessage = "Nginx没有正确的响应"
		result = append(result, t)
	}

	{
		t := NewThreshold()
		t.Param = "${cost}"
		t.NoticeLevel = notices.NoticeLevelWarning
		t.Operator = ThresholdOperatorGt
		t.Value = "5"
		t.NoticeMessage = "Nginx请求时间超过5秒"
		result = append(result, t)
	}

	{
		t := NewThreshold()
		t.Param = "${result.waitingConnections}"
		t.NoticeLevel = notices.NoticeLevelWarning
		t.Operator = ThresholdOperatorGt
		t.Value = "20000"
		t.NoticeMessage = "Nginx等待请求的连接数过多，超过20000"
		result = append(result, t)
	}

	return result
}

// 图表
func (this *NginxStatusSource) Charts() []*widgets.Chart {
	charts := []*widgets.Chart{}

	{
		// chart
		chart := widgets.NewChart()
		chart.Name = "等待连接数"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.SupportsTimeRange = true
		chart.Options = maps.Map{
			"code": `var chart = new charts.LineChart();

var ones = NewQuery().past(60, time.MINUTE).avg("result.waitingConnections");

var line = new charts.Line();
line.isFilled = true;

ones.$each(function (k, v) {
	line.addValue(v.value.result.waitingConnections);
	chart.addLabel(v.label);
});

chart.addLine(line);
chart.render();`,
		}

		charts = append(charts, chart)
	}

	{
		// chart
		chart := widgets.NewChart()
		chart.Name = "平均请求数<em>（每秒）</em>"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.SupportsTimeRange = true
		chart.Options = maps.Map{
			"code": `var chart = new charts.LineChart();

var ones = NewQuery().past(60, time.MINUTE).avg("result.requestsPerSecond");

var line = new charts.Line();
line.isFilled = true;

ones.$each(function (k, v) {
	line.addValue(v.value.result.requestsPerSecond);
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
func (this *NginxStatusSource) Presentation() *forms.Presentation {
	return &forms.Presentation{
		HTML: `
<tr>
	<td>Nginx状态检查URL</td>
	<td>{{source.url}}</td>
</tr>
`,
	}
}
