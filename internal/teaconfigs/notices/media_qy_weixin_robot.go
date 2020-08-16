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

// 企业微信群机器人媒介
type NoticeQyWeixinRobotMedia struct {
	WebhookURL string     `yaml:"webhookURL" json:"webhookURL"`
	TextFormat TextFormat `yaml:"textFormat" json:"textFormat"`
}

// 获取新对象
func NewNoticeQyWeixinRobotMedia() *NoticeQyWeixinRobotMedia {
	return &NoticeQyWeixinRobotMedia{}
}

func (this *NoticeQyWeixinRobotMedia) Send(user string, subject string, body string) (resp []byte, err error) {
	if len(this.WebhookURL) == 0 {
		return nil, errors.New("webhook url should not be empty")
	}

	mobiles := []string{}
	if len(user) > 0 {
		for _, u := range strings.Split(user, ",") {
			u = strings.TrimSpace(u)
			if len(u) > 0 {
				mobiles = append(mobiles, u)
			}
		}
	}

	content := maps.Map{}
	if this.TextFormat == FormatMarkdown { // markdown
		content = maps.Map{
			"msgtype": "markdown",
			"markdown": maps.Map{
				"content":               subject + "\n" + body,
				"mentioned_mobile_list": mobiles,
			},
		}
	} else {
		content = maps.Map{
			"msgtype": "text",
			"text": maps.Map{
				"content":               subject + "\n" + body,
				"mentioned_mobile_list": mobiles,
			},
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
func (this *NoticeQyWeixinRobotMedia) RequireUser() bool {
	return false
}
