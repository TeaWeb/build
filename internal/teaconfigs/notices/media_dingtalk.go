package notices

import (
	"bytes"
	"errors"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/string"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// 钉钉群机器人媒介
type NoticeDingTalkMedia struct {
	WebhookURL string `yaml:"webhookURL" json:"webhookURL"`
}

// 获取新对象
func NewNoticeDingTalkMedia() *NoticeDingTalkMedia {
	return &NoticeDingTalkMedia{}
}

func (this *NoticeDingTalkMedia) Send(user string, subject string, body string) (resp []byte, err error) {
	if len(this.WebhookURL) == 0 {
		return nil, errors.New("webhook url should not be empty")
	}

	content := maps.Map{
		"msgtype": "text",
		"text": maps.Map{
			"content": "标题：" + subject + "\n内容：" + body,
		},
	}
	if len(user) > 0 {
		mobiles := []string{}
		for _, u := range strings.Split(user, ",") {
			u = strings.TrimSpace(u)
			if len(u) > 0 {
				mobiles = append(mobiles, u)
			}
		}

		content["at"] = maps.Map{
			"atMobiles": mobiles,
		}
	}

	reader := bytes.NewBufferString(stringutil.JSONEncode(content))
	req, err := http.NewRequest(http.MethodPost, this.WebhookURL, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := teautils.SharedHttpClient(5 * time.Second)
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	resp, err = ioutil.ReadAll(response.Body)
	return
}

// 是否需要用户标识
func (this *NoticeDingTalkMedia) RequireUser() bool {
	return false
}
