package agents

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs/forms"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teautils"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// WebHook
type WebHookSource struct {
	Source `yaml:",inline"`

	URL      string             `yaml:"url" json:"url"`
	Timeout  string             `yaml:"timeout" json:"timeout"`
	Method   string             `yaml:"method" json:"method"` // 请求方法
	Headers  []*shared.Variable `yaml:"headers" json:"headers"`
	Params   []*shared.Variable `yaml:"params" json:"params"`
	TextBody string             `yaml:"textBody" json:"textBody"`

	timeoutDuration time.Duration
}

// 获取新对象
func NewWebHookSource() *WebHookSource {
	return &WebHookSource{}
}

// 校验
func (this *WebHookSource) Validate() error {
	this.timeoutDuration, _ = time.ParseDuration(this.Timeout)
	if len(this.Method) == 0 {
		this.Method = http.MethodPost
	} else {
		this.Method = strings.ToUpper(this.Method)
	}

	if len(this.URL) == 0 {
		return errors.New("url should not be empty")
	}

	return nil
}

// 名称
func (this *WebHookSource) Name() string {
	return "WebHook"
}

// 代号
func (this *WebHookSource) Code() string {
	return "webhook"
}

// 描述
func (this *WebHookSource) Description() string {
	return "通过HTTP或者HTTPS接口获取数据"
}

// 执行
func (this *WebHookSource) Execute(params map[string]string) (value interface{}, err error) {
	if this.timeoutDuration.Seconds() <= 0 {
		this.timeoutDuration = 10 * time.Second
	}

	client := teautils.SharedHttpClient(this.timeoutDuration)

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

	req, err := http.NewRequest(this.Method, urlString, body)
	if err != nil {
		return nil, err
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

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("response status code should be 200, but is " + fmt.Sprintf("%d", resp.StatusCode))
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return DecodeSource(respBytes, this.DataFormat)
}

// 选项表单
func (this *WebHookSource) Form() *forms.Form {
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
	"timeout": values.timeout
};
`
			group.Add(field)
		}
	}

	form.ValidateCode = `
var url = values.webhook.url;
if (url.length == 0) {
	return FieldError("url", "请输入URL")
}

if (!url.match(/^(http|https):\/\//i)) {
	return FieldError("url", "URL地址必须以http或https开头");
}

var method = values.webhook.method;
if (method.length == 0) {
	return FieldError("method", "请选择请求方法");
}

var timeout = values.webhook.timeout;
if (!timeout.match(/^\d+(s|ms)$/)) {
	return FieldError("timeout", "超时时间只能是一个整数");
}

return {
	"url": values.webhook.url,
	"method": values.webhook.method,
	"params": values.webhook.params,
	"headers": values.webhook.headers,
	"textBody": values.webhook.textBody,
	"timeout": values.webhook.timeout
}
`

	return form
}

func (this *WebHookSource) Presentation() *forms.Presentation {
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
			<pre class="webhook-block-body">{{source.textBody}}</pre>
		</div>
	</td>
</tr>
<tr>
	<td class="color-border">请求超时<em>（Timeout）</em></td>
	<td>{{source.timeout}}</td>
</tr>`,
		CSS: `.webhook-block-body {
    border: 1px #eee solid;
    padding: 0.4em;
    background: rgba(0, 0, 0, 0.01);
    font-size: 0.9em;
    max-height: 10em;
    overflow-y: auto;
    margin: 0;
}

.webhook-block-body::-webkit-scrollbar {
    width: 4px;
}
`,
	}
}

func (this *WebHookSource) lookupHeader(name string) (value string, ok bool) {
	for _, h := range this.Headers {
		if h.Name == name {
			return h.Value, true
		}
	}
	return "", false
}
