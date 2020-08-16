package notices

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/maps"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// 企业微信媒介
type NoticeQyWeixinMedia struct {
	CorporateId string     `yaml:"corporateId" json:"corporateId"`
	AgentId     string     `yaml:"agentId" json:"agentId"`
	AppSecret   string     `yaml:"appSecret" json:"appSecret"`
	TextFormat  TextFormat `yaml:"textFormat" json:"textFormat"`
}

// 获取新对象
func NewNoticeQyWeixinMedia() *NoticeQyWeixinMedia {
	return &NoticeQyWeixinMedia{}
}

func (this *NoticeQyWeixinMedia) Send(user string, subject string, body string) (respData []byte, err error) {
	// 获取Token
	u := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + url.QueryEscape(this.CorporateId) + "&corpsecret=" + url.QueryEscape(this.AppSecret)
	req1, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	client := teautils.SharedHttpClient(5 * time.Second)
	resp1, err := client.Do(req1)
	if err != nil {
		return nil, err
	}
	defer resp1.Body.Close()

	if resp1.StatusCode != http.StatusOK {
		return nil, errors.New("status code not 200")
	}

	data, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		return nil, err
	}

	m := maps.Map{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return data, err
	}
	errCode := m.GetInt("errcode")
	if errCode > 0 {
		return data, errors.New("error code:" + fmt.Sprintf("%d", errCode))
	}

	accessToken := m.GetString("access_token")

	if len(user) == 0 {
		user = "@all"
	}

	// 发送消息
	u = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=" + url.QueryEscape(accessToken)

	msg := maps.Map{
		"touser":  user,
		"toparty": "",
		"totag":   "",
		"toall":   0,
		"agentid": this.AgentId,
		"safe":    0,
	}
	if this.TextFormat == FormatMarkdown {
		msg["msgtype"] = "markdown"
		msg["markdown"] = maps.Map{
			"content": subject + "\n" + body,
		}
	} else {
		msg["msgtype"] = "text"
		msg["text"] = maps.Map{
			"content": subject + "\n" + body,
		}
	}
	data, err = json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, u, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("status not 200")
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	m = maps.Map{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	errCode = m.GetInt("errcode")
	if errCode != 0 {
		return data, errors.New("errcode " + fmt.Sprintf("%d", errCode))
	}

	invalidUser := m.GetString("invaliduser")
	if len(invalidUser) > 0 {
		return data, errors.New("invalid users:" + invalidUser)
	}
	return data, nil
}

// 是否需要用户标识
func (this *NoticeQyWeixinMedia) RequireUser() bool {
	return false
}
