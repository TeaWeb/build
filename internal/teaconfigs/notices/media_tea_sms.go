package notices

import (
	"encoding/json"
	"errors"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/maps"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TeaOS云短信
type NoticeTeaSmsMedia struct {
	Sign         string `yaml:"sign" json:"sign"`                 // 签名名称
	AccessId     string `yaml:"accessId" json:"accessId"`         // AccessId
	AccessSecret string `yaml:"accessSecret" json:"accessSecret"` // AccessSecret
}

// 获取新对象
func NewNoticeTeaSmsMedia() *NoticeTeaSmsMedia {
	return &NoticeTeaSmsMedia{}
}

func (this *NoticeTeaSmsMedia) Send(user string, subject string, body string) (respBytes []byte, err error) {
	apiURL := "http://cloud.teaos.cn/api/v1/sms/send"
	params := url.Values{}
	params.Set("mobile", user)
	params.Set("subject", subject)
	params.Set("body", body)

	// 调试用
	/**if Tea.IsTesting() {
		params.Set("debug", "1")
	}**/

	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "TeaWeb/"+teaconst.TeaVersion)
	req.Header.Set("Tea-Access-Id", this.AccessId)
	req.Header.Set("Tea-Access-Secret", this.AccessSecret)

	client := teautils.SharedHttpClient(5 * time.Second)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	m := maps.Map{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return data, err
	}

	if m.GetInt("code") == 200 {
		return data, nil
	}

	return nil, errors.New("发送失败：" + m.GetString("message"))
}

// 是否需要用户标识
func (this *NoticeTeaSmsMedia) RequireUser() bool {
	return true
}
