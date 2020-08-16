package fastcgi

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"regexp"
)

type UpdateAction actions.Action

// 修改
func (this *UpdateAction) Run(params struct {
	From       string
	ServerId   string
	LocationId string
	FastcgiId  string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	fastcgiList, err := server.FindFastcgiList(params.LocationId)
	if err != nil {
		this.Fail(err.Error())
	}
	fastcgi := fastcgiList.FindFastcgi(params.FastcgiId)
	if fastcgi == nil {
		this.Fail("找不到要修改的Fastcgi")
	}

	m := maps.Map{
		"on":              fastcgi.On,
		"id":              fastcgi.Id,
		"pass":            fastcgi.Pass,
		"poolSize":        fastcgi.PoolSize,
		"params":          fastcgi.Params,
		"pathInfoPattern": fastcgi.PathInfoPattern,
	}
	if fastcgi.ReadTimeout != "0s" {
		m["readTimeoutSeconds"] = int(fastcgi.ReadTimeoutDuration().Seconds())
	} else {
		m["readTimeoutSeconds"] = 0
	}
	this.Data["fastcgi"] = m

	this.Data["from"] = params.From
	this.Data["server"] = maps.Map{
		"id": params.ServerId,
	}
	this.Data["locationId"] = params.LocationId

	this.Show()
}

// 修改
func (this *UpdateAction) RunPost(params struct {
	ServerId        string
	LocationId      string
	On              bool
	Pass            string
	ReadTimeout     int
	ParamNames      []string
	ParamValues     []string
	PoolSize        int
	PathInfoPattern string

	FastcgiId string

	Must *actions.Must
}) {
	params.Must.
		Field("pass", params.Pass).
		Require("请输入Fastcgi地址").
		Field("poolSize", params.PoolSize).
		Gte(0, "连接池尺寸不能小于0")

	// PATH_INFO
	if len(params.PathInfoPattern) > 0 {
		_, err := regexp.Compile(params.PathInfoPattern)
		if err != nil {
			this.FailField("pathInfoPattern", "PATH_INFO匹配规则错误："+err.Error())
		}
	}

	paramsMap := map[string]string{}
	for index, paramName := range params.ParamNames {
		if index < len(params.ParamValues) {
			paramsMap[paramName] = params.ParamValues[index]
		}
	}

	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	fastcgiList, err := server.FindFastcgiList(params.LocationId)
	if err != nil {
		this.Fail(err.Error())
	}

	fastcgi := fastcgiList.FindFastcgi(params.FastcgiId)
	if fastcgi == nil {
		this.Fail("找不到要修改的Fastcgi")
	}

	fastcgi.On = params.On
	fastcgi.Pass = teautils.FormatAddress(params.Pass)
	fastcgi.ReadTimeout = fmt.Sprintf("%ds", params.ReadTimeout)
	fastcgi.Params = paramsMap
	fastcgi.PoolSize = params.PoolSize
	fastcgi.PathInfoPattern = params.PathInfoPattern
	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()

	this.Success()
}
