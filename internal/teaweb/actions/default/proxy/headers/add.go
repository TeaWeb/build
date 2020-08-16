package headers

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type AddAction actions.Action

// 添加Header
func (this *AddAction) Run(params struct {
	From       string
	ServerId   string
	LocationId string
	RewriteId  string
	FastcgiId  string
	BackendId  string
}) {
	this.Data["from"] = params.From
	this.Data["server"] = maps.Map{
		"id": params.ServerId,
	}
	this.Data["locationId"] = params.LocationId
	this.Data["rewriteId"] = params.RewriteId
	this.Data["fastcgiId"] = params.FastcgiId
	this.Data["backendId"] = params.BackendId

	this.Show()
}

// 提交保存
func (this *AddAction) RunPost(params struct {
	ServerId   string
	LocationId string
	RewriteId  string
	FastcgiId  string
	BackendId  string

	On         bool
	Name       string
	Value      string
	AllStatus  bool
	StatusList []int

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入名称")

	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	headerList, err := server.FindHeaderList(params.LocationId, params.BackendId, params.RewriteId, params.FastcgiId)
	if err != nil {
		this.Fail(err.Error())
	}

	header := shared.NewHeaderConfig()
	header.On = params.On
	header.Name = params.Name
	header.Value = params.Value
	header.Always = params.AllStatus
	header.Status = params.StatusList
	headerList.AddResponseHeader(header)

	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()

	this.Success()
}
