package update

import (
	"encoding/json"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"
)

type IndexAction actions.Action

// 检查版本更新
func (this *IndexAction) Run(params struct{}) {
	this.Data["currentVersion"] = teaconst.TeaVersion

	this.Show()
}

// 开始检查
func (this *IndexAction) RunPost() {
	url := "http://teaos.cn/services/version?version=" + teaconst.TeaVersion

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		this.Fail("发生错误：" + err.Error())
	}
	req.Header.Set("User-Agent", runtime.GOOS+"/"+runtime.GOARCH)

	client := teautils.SharedHttpClient(5 * time.Second)
	resp, err := client.Do(req)
	if err != nil {
		this.Fail("无法连接TeaWeb服务：" + err.Error())
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		this.Fail("发生错误：" + err.Error())
	}

	m := maps.Map{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		this.Fail("发生错误：" + err.Error())
	}

	dataMap := maps.NewMap(m.Get("data"))

	hasNew := dataMap.GetBool("hasNew")
	this.Data["hasNew"] = hasNew
	this.Data["latest"] = dataMap.GetString("latest")

	this.Success()
}
