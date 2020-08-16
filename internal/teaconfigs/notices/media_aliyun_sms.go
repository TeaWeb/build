package notices

import (
	"encoding/json"
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/string"
	"strings"
)

// 阿里云短信
type NoticeAliyunSmsMedia struct {
	Sign            string             `yaml:"sign" json:"sign"`                       // 签名名称
	TemplateCode    string             `yaml:"templateCode" json:"templateCode"`       // 模板CODE
	Variables       []*shared.Variable `yaml:"variables" json:"variables"`             // 变量
	AccessKeyId     string             `yaml:"accessKeyId" json:"accessKeyId"`         // AccessKeyId
	AccessKeySecret string             `yaml:"accessKeySecret" json:"accessKeySecret"` // AccessKeySecret
}

// 获取新对象
func NewNoticeAliyunSmsMedia() *NoticeAliyunSmsMedia {
	return &NoticeAliyunSmsMedia{}
}

func (this *NoticeAliyunSmsMedia) Send(user string, subject string, body string) (resp []byte, err error) {
	// {"Message":"OK","RequestId":"xxxx","BizId":"xxxx","Code":"OK"}

	client, err := sdk.NewClientWithAccessKey("cn-hangzhou", this.AccessKeyId, this.AccessKeySecret)
	if err != nil {
		panic(err)
	}

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-hangzhou"
	request.QueryParams["PhoneNumbers"] = user
	request.QueryParams["SignName"] = this.Sign
	request.QueryParams["TemplateCode"] = this.TemplateCode

	varMap := maps.Map{}
	for _, v := range this.Variables {
		value := v.Value
		value = strings.Replace(value, "${NoticeUser}", user, -1)
		value = strings.Replace(value, "${NoticeSubject}", subject, -1)
		value = strings.Replace(value, "${NoticeBody}", body, -1)

		// 阿里云的限制参数长度不能超过20：
		maxLen := 20
		if len([]rune(value)) > maxLen {
			value = string([]rune(value)[:maxLen-4]) + "..."
		}

		varMap[v.Name] = value
	}
	request.QueryParams["TemplateParam"] = stringutil.JSONEncode(varMap)

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return nil, err
	}

	data := response.GetHttpContentBytes()
	m := maps.Map{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return data, err
	}

	if m.GetString("Code") == "OK" {
		return data, nil
	}
	return data, errors.New("fail to send sms：" + string(data))
}

// 是否需要用户标识
func (this *NoticeAliyunSmsMedia) RequireUser() bool {
	return true
}
