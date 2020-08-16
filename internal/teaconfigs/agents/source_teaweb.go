package agents

import (
	"encoding/json"
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs/forms"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/maps"
	"io/ioutil"
	"net/http"
	"time"
)

// TeaWeb相关数据源
type TeaWebSource struct {
	Source `yaml:",inline"`

	API     string `yaml:"api" json:"api"`
	Timeout int    `yaml:"timeout" json:"timeout"`
}

// 获取新对象
func NewTeaWebSource() *TeaWebSource {
	return &TeaWebSource{}
}

// 名称
func (this *TeaWebSource) Name() string {
	return "TeaWeb"
}

// 代号
func (this *TeaWebSource) Code() string {
	return "teaweb"
}

// 描述
func (this *TeaWebSource) Description() string {
	return "通过TeaWeb API监控其他TeaWeb"
}

// 执行
func (this *TeaWebSource) Execute(params map[string]string) (value interface{}, err error) {
	if len(this.API) == 0 {
		return nil, errors.New("API address should not be empty")
	}

	before := time.Now()
	req, err := http.NewRequest(http.MethodGet, this.API, nil)
	if err != nil {
		value = maps.Map{
			"cost":   time.Since(before).Seconds(),
			"status": 0,
			"result": "",
			"length": 0,
		}
		return value, err
	}
	req.Header.Set("User-Agent", teaconst.TeaProductCode+"/"+teaconst.TeaVersion)

	timeout := this.Timeout
	if timeout <= 0 {
		timeout = 30
	}

	client := teautils.SharedHttpClient(time.Duration(timeout) * time.Second)
	resp, err := client.Do(req)
	if err != nil {
		return maps.Map{
			"status": 0,
			"cost":   time.Since(before).Seconds(),
			"result": maps.Map{},
		}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return maps.Map{
			"status": resp.StatusCode,
			"cost":   time.Since(before).Seconds(),
			"result": maps.Map{},
		}, nil
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return maps.Map{
			"status": resp.StatusCode,
			"cost":   time.Since(before).Seconds(),
			"result": maps.Map{},
		}, err
	}

	m := maps.Map{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return maps.Map{
			"status": resp.StatusCode,
			"cost":   time.Since(before).Seconds(),
			"result": maps.Map{},
		}, err
	}

	return maps.Map{
		"status": resp.StatusCode,
		"cost":   time.Since(before).Seconds(),
		"result": m,
	}, nil
}

// 表单信息
func (this *TeaWebSource) Form() *forms.Form {
	form := forms.NewForm(this.Code())
	{
		group := form.NewGroup()
		{
			field := forms.NewTextBox("API地址", "")
			field.Rows = 2
			field.Comment = "格式为：http://TeaWeb访问地址/api/monitor?TeaKey=登录用户密钥 <a href=\"http://teaos.cn/doc/advanced/APIMonitor.md\" target=\"_blank\">说明文档&raquo;</a>"
			field.Code = "api"
			field.IsRequired = true
			field.MaxLength = 200
			field.Attr("style", "word-wrap:break-word;word-break:break-all;line-height:1.5")
			field.ValidateCode = `
if (value.length == 0) {
	throw new Error("请输入API地址")
}

if (!value.match(/^(http|https):\/\//i)) {
	throw new Error("URL地址必须以http或https开头");
}
`
			group.Add(field)
		}

		{
			group := form.NewGroup()

			{
				field := forms.NewTextField("请求超时", "Timeout")
				field.Code = "timeout"
				field.Value = 10
				field.MaxLength = 10
				field.RightLabel = "秒"
				field.Attr("style", "width:5em")
				field.ValidateCode = `
var intValue = parseInt(value);
if (isNaN(intValue)) {
	throw new Error("超时时间需要是一个整数");
}

return intValue;
`

				group.Add(field)
			}
		}
	}
	return form
}

// 变量
func (this *TeaWebSource) Variables() []*SourceVariable {
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
			Code:        "result.arch",
			Description: "CPU架构",
		},
		{
			Code:        "result.heap",
			Description: "Heap值（字节）",
		},
		{
			Code:        "result.memory",
			Description: "总内存值（字节）",
		},
		{
			Code:        "result.mongo",
			Description: "MongoDB连接是否正常",
		},
		{
			Code:        "result.os",
			Description: "操作系统代号",
		},
		{
			Code:        "result.routines",
			Description: "go routine数量",
		},
		{
			Code:        "result.version",
			Description: "TeaWeb版本",
		},
	}
}

// 阈值
func (this *TeaWebSource) Thresholds() []*Threshold {
	result := []*Threshold{}

	{
		t := NewThreshold()
		t.Param = "${status}"
		t.NoticeLevel = notices.NoticeLevelWarning
		t.Operator = ThresholdOperatorNot
		t.Value = "200"
		t.NoticeMessage = "TeaWeb没有正确的响应"
		result = append(result, t)
	}

	{
		t := NewThreshold()
		t.Param = "${cost}"
		t.NoticeLevel = notices.NoticeLevelWarning
		t.Operator = ThresholdOperatorGt
		t.Value = "5"
		t.NoticeMessage = "TeaWeb请求时间超过5秒"
		result = append(result, t)
	}

	{
		t := NewThreshold()
		t.Param = "${result.memory}"
		t.NoticeLevel = notices.NoticeLevelWarning
		t.Operator = ThresholdOperatorGt
		t.Value = "1073741824"
		t.NoticeMessage = "TeaWeb使用内存超过1G"
		result = append(result, t)
	}

	return result
}

// 图表
func (this *TeaWebSource) Charts() []*widgets.Chart {
	charts := []*widgets.Chart{}
	return charts
}

// 显示信息
func (this *TeaWebSource) Presentation() *forms.Presentation {
	p := forms.NewPresentation()
	p.HTML = `
<tr>
	<td>API地址</td>
	<td>{{source.api}}</td>
</tr>
<tr>
	<td>请求超时</td>
	<td>{{source.timeout}}s</td>
</tr>
`
	return p
}
