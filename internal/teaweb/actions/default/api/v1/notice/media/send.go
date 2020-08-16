package media

import (
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/pquerna/ffjson/ffjson"
	"io/ioutil"
)

type SendAction actions.Action

// 通过某个媒介发送信息
func (this *SendAction) RunPost(params struct {
	MediaId string
}) {
	setting := notices.SharedNoticeSetting()
	media := setting.FindMedia(params.MediaId)
	if media == nil {
		apiutils.Fail(this, "can not find media")
		return
	}

	rawMedia, err := media.Raw()
	if err != nil {
		apiutils.Fail(this, err.Error())
		return
	}

	requestBodyBytes, err := ioutil.ReadAll(this.Request.Body)
	if err != nil {
		apiutils.Fail(this, err.Error())
		return
	}

	if len(requestBodyBytes) == 0 {
		apiutils.Fail(this, "request body should be a valid json")
		return
	}

	paramMap := maps.Map{}
	err = ffjson.Unmarshal(requestBodyBytes, &paramMap)
	if err != nil {
		apiutils.Fail(this, err.Error())
		return
	}

	if paramMap.Len() == 0 {
		apiutils.Fail(this, "request body should be a valid json")
		return
	}

	user := paramMap.GetString("user")
	subject := paramMap.GetString("subject")
	body := paramMap.GetString("body")

	if rawMedia.RequireUser() {
		if len(user) == 0 {
			apiutils.Fail(this, "please input 'user' parameter")
			return
		}
	}

	resp, err := rawMedia.Send(user, subject, body)
	if err != nil {
		apiutils.Fail(this, err.Error())
		return
	}

	apiutils.Success(this, maps.Map{
		"result": string(resp),
	})
}
