package agents

import (
	"bytes"
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs/forms"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/maps"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// URL连通性
type URLConnectivitySource struct {
	Source `yaml:",inline"`

	Timeout  int                `yaml:"timeout" json:"timeout"` // 连接超时
	URL      string             `yaml:"url" json:"url"`
	Method   string             `yaml:"method" json:"method"`
	Headers  []*shared.Variable `yaml:"headers" json:"headers"`
	Params   []*shared.Variable `yaml:"params" json:"params"`
	TextBody string             `yaml:"textBody" json:"textBody"`
}

// 获取新对象
func NewURLConnectivitySource() *URLConnectivitySource {
	return &URLConnectivitySource{}
}

// 名称
func (this *URLConnectivitySource) Name() string {
	return "URL连通性"
}

// 代号
func (this *URLConnectivitySource) Code() string {
	return "urlConnectivity"
}

// 描述
func (this *URLConnectivitySource) Description() string {
	return "获取URL连通性信息"
}

// 执行
func (this *URLConnectivitySource) Execute(params map[string]string) (value interface{}, err error) {
	if len(this.URL) == 0 {
		err = errors.New("'url' should not be empty")
		return maps.Map{
			"status": 0,
		}, err
	}

	method := this.Method
	if len(method) == 0 {
		method = http.MethodGet
	}

	urlString := this.URL
	var body io.Reader = nil
	if this.Method == "PUT" {
		body = bytes.NewReader([]byte(this.TextBody))
	} else {
		query := url.Values{}
		for name, value := range params {
			query[name] = []string{value}
		}
		for _, param := range this.Params {
			_, ok := query[param.Name]
			if ok {
				query[param.Name] = append(query[param.Name], param.Value)
			} else {
				query[param.Name] = []string{param.Value}
			}
		}
		rawQuery := query.Encode()

		if len(query) > 0 {
			if this.Method == "GET" {
				if strings.Index(this.URL, "?") > 0 {
					urlString += "&" + rawQuery
				} else {
					urlString += "?" + rawQuery
				}
			} else {
				body = bytes.NewReader([]byte(rawQuery))
			}
		} else if this.Method == "POST" {
			body = bytes.NewReader([]byte(this.TextBody))
		}
	}

	before := time.Now()
	req, err := http.NewRequest(method, this.URL, body)
	if err != nil {
		value = maps.Map{
			"cost":   time.Since(before).Seconds(),
			"status": 0,
			"result": "",
			"length": 0,
		}
		return value, err
	}

	for _, h := range this.Headers {
		req.Header.Add(h.Name, h.Value)
	}

	_, ok := this.lookupHeader("User-Agent")
	if !ok {
		req.Header.Set("User-Agent", teaconst.TeaProductCode+"/"+teaconst.TeaVersion)
	}

	if this.Method == "POST" {
		_, ok := this.lookupHeader("Content-Type")
		if !ok {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	timeout := this.Timeout
	if timeout <= 0 {
		timeout = 30
	}

	client := teautils.SharedHttpClient(time.Duration(timeout) * time.Second)
	resp, err := client.Do(req)
	if err != nil {
		return maps.Map{
			"status": 0,
		}, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return maps.Map{
			"status": 0,
		}, err
	}

	if len(data) > 1024 {
		data = data[:1024]
	}

	value = maps.Map{
		"cost":   time.Since(before).Seconds(),
		"status": resp.StatusCode,
		"result": string(data),
		"length": len(data),
	}

	return
}

// 表单信息
func (this *URLConnectivitySource) Form() *forms.Form {
	form := forms.NewForm(this.Code())
	{
		group := form.NewGroup()
		{
			field := forms.NewHTTPBox("", "")
			field.Code = this.Code()
			field.InitCode = `
return {
	"url": values.url,
	"method": values.method,
	"params": values.params,
	"headers": values.headers,
	"textBody": values.textBody,
	"timeout": values.timeout + "s"
};
`
			group.Add(field)
		}
	}

	form.ValidateCode = `
var url = values.urlConnectivity.url;
if (url.length == 0) {
	return FieldError("url", "请输入URL")
}

if (!url.match(/^(http|https):\/\//i)) {
	return FieldError("url", "URL地址必须以http或https开头");
}

var method = values.urlConnectivity.method;
if (method.length == 0) {
	return FieldError("method", "请选择请求方法");
}

var timeout = values.urlConnectivity.timeout;
if (!timeout.match(/^\d+(s|ms)$/)) {
	return FieldError("timeout", "超时时间只能是一个整数");
}
timeout = parseInt(timeout.replace(/(s|ms)/, ""));

return {
	"url": values.urlConnectivity.url,
	"method": values.urlConnectivity.method,
	"params": values.urlConnectivity.params,
	"headers": values.urlConnectivity.headers,
	"textBody": values.urlConnectivity.textBody,
	"timeout": timeout
}
`

	return form
}

// 变量
func (this *URLConnectivitySource) Variables() []*SourceVariable {
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
			Code:        "result",
			Description: "响应内容文本，最多只记录前1024个字节",
		},
		{
			Code:        "length",
			Description: "响应的内容长度",
		},
	}
}

// 阈值
func (this *URLConnectivitySource) Thresholds() []*Threshold {
	result := []*Threshold{}

	// 阈值
	{
		t := NewThreshold()
		t.Param = "${status}"
		t.Operator = ThresholdOperatorGte
		t.Value = "400"
		t.NoticeLevel = notices.NoticeLevelWarning
		t.NoticeMessage = "请求状态码错误：${status}"
		result = append(result, t)
	}

	// 阈值
	{
		t := NewThreshold()
		t.Param = "${status}"
		t.Operator = ThresholdOperatorEq
		t.Value = "0"
		t.NoticeLevel = notices.NoticeLevelWarning
		t.NoticeMessage = "URL请求失败"
		result = append(result, t)
	}

	return result
}

// 图表
func (this *URLConnectivitySource) Charts() []*widgets.Chart {
	charts := []*widgets.Chart{}

	{
		// chart
		chart := widgets.NewChart()
		chart.Name = "URL连通性（ms）"
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
chart.render();`,
		}

		charts = append(charts, chart)
	}

	return charts
}

// 显示信息
func (this *URLConnectivitySource) Presentation() *forms.Presentation {
	return &forms.Presentation{
		HTML: `
			<tr>
				<td class="color-border">URL</td>
				<td>{{source.url}}</td>
			</tr>
			<tr>
				<td class="color-border">请求方法</td>
				<td>{{source.method}}</td>
			</tr>
			<tr>
				<td class="color-border">自定义Header</td>
				<td>
					<span v-if="source.headers == null || source.headers.length == 0" class="disabled">还没有自定义Header。</span>
					<div v-if="source.headers != null && source.headers.length > 0">
						<span class="ui label tiny" v-for="header in source.headers">{{header.name}}: {{header.value}}</span>
					</div>
				</td>
			</tr>
			<tr v-if="source.method == 'POST' || source.method == 'PUT'">
				<td class="color-border">自定义请求内容</td>
				<td>
					<span v-if="(source.params == null || source.params.length == 0) && (source.textBody == null || source.textBody.length == 0)" class="disabled">还没有自定义请求内容。</span>
					<div v-if="source.params != null && source.params.length > 0">
						<span class="ui label tiny" v-for="param in source.params">{{param.name}}: {{param.value}}</span>
					</div>
					<div v-if="source.textBody != null && source.textBody.length > 0">
						<pre class="urlConnectivity-block-body">{{source.textBody}}</pre>
					</div>
				</td>
			</tr>
			<tr>
				<td class="color-border">请求超时<em>（Timeout）</em></td>
				<td>{{source.timeout}}s</td>
			</tr>`,
		CSS: `.urlConnectivity-block-body {
			border: 1px #eee solid;
			padding: 0.4em;
			background: rgba(0, 0, 0, 0.01);
			font-size: 0.9em;
			max-height: 10em;
			overflow-y: auto;
			margin: 0;
		}
		
		.urlConnectivity-block-body::-webkit-scrollbar {
			width: 4px;
		}
		`,
	}
}

func (this *URLConnectivitySource) lookupHeader(name string) (value string, ok bool) {
	for _, h := range this.Headers {
		if h.Name == name {
			return h.Value, true
		}
	}
	return "", false
}
