package notices

import (
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/utils/string"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Webhook媒介
type NoticeWebhookMedia struct {
	URL         string             `yaml:"url" json:"url"` // URL中可以使用${NoticeSubject}, ${NoticeBody}两个变量
	Method      string             `yaml:"method" json:"method"`
	ContentType string             `yaml:"contentType" json:"contentType"` // 内容类型：params|body
	Headers     []*shared.Variable `yaml:"headers" json:"headers"`
	Params      []*shared.Variable `yaml:"params" json:"params"`
	Body        string             `yaml:"body" json:"body"`
}

// 获取新对象
func NewNoticeWebhookMedia() *NoticeWebhookMedia {
	return &NoticeWebhookMedia{}
}

// 发送
func (this *NoticeWebhookMedia) Send(user string, subject string, body string) (resp []byte, err error) {
	if len(this.URL) == 0 {
		return nil, errors.New("'url' should be specified")
	}

	timeout := 10 * time.Second

	if len(this.Method) == 0 {
		this.Method = http.MethodGet
	}

	this.URL = strings.Replace(this.URL, "${NoticeUser}", url.QueryEscape(user), -1)
	this.URL = strings.Replace(this.URL, "${NoticeSubject}", url.QueryEscape(subject), -1)
	this.URL = strings.Replace(this.URL, "${NoticeBody}", url.QueryEscape(body), -1)

	var req *http.Request
	if this.Method == http.MethodGet {
		req, err = http.NewRequest(this.Method, this.URL, nil)
	} else {
		params := url.Values{
			"NoticeUser":    []string{user},
			"NoticeSubject": []string{subject},
			"NoticeBody":    []string{body},
		}

		postBody := ""
		if this.ContentType == "params" {
			for _, param := range this.Params {
				param.Value = strings.Replace(param.Value, "${NoticeUser}", user, -1)
				param.Value = strings.Replace(param.Value, "${NoticeSubject}", subject, -1)
				param.Value = strings.Replace(param.Value, "${NoticeBody}", body, -1)
				params.Add(param.Name, param.Value)
			}
			postBody = params.Encode()
		} else if this.ContentType == "body" {
			userJSON := stringutil.JSONEncode(user)
			subjectJSON := stringutil.JSONEncode(subject)
			bodyJSON := stringutil.JSONEncode(body)
			if len(userJSON) > 0 {
				userJSON = userJSON[1 : len(userJSON)-1]
			}
			if len(subjectJSON) > 0 {
				subjectJSON = subjectJSON[1 : len(subjectJSON)-1]
			}
			if len(bodyJSON) > 0 {
				bodyJSON = bodyJSON[1 : len(bodyJSON)-1]
			}
			postBody = strings.Replace(this.Body, "${NoticeUser}", userJSON, -1)
			postBody = strings.Replace(postBody, "${NoticeSubject}", subjectJSON, -1)
			postBody = strings.Replace(postBody, "${NoticeBody}", bodyJSON, -1)
		} else {
			postBody = params.Encode()
		}

		req, err = http.NewRequest(this.Method, this.URL, strings.NewReader(postBody))
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", teaconst.TeaProductCode+"/"+teaconst.TeaVersion)

	if len(this.Headers) > 0 {
		for _, h := range this.Headers {
			req.Header.Set(h.Name, h.Value)
		}
	}

	client := teautils.SharedHttpClient(timeout)
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	data, err := ioutil.ReadAll(response.Body)
	return data, err
}

// 是否需要用户标识
func (this *NoticeWebhookMedia) RequireUser() bool {
	return false
}

// 添加Header
func (this *NoticeWebhookMedia) AddHeader(name string, value string) {
	this.Headers = append(this.Headers, &shared.Variable{
		Name:  name,
		Value: value,
	})
}

// 添加参数
func (this *NoticeWebhookMedia) AddParam(name string, value string) {
	this.Params = append(this.Params, &shared.Variable{
		Name:  name,
		Value: value,
	})
}
