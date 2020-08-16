package headers

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateRequestHeaderAction actions.Action

// 修改
func (this *UpdateRequestHeaderAction) Run(params struct {
	From       string
	ServerId   string
	LocationId string
	RewriteId  string
	FastcgiId  string
	BackendId  string
	HeaderId   string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	this.Data["from"] = params.From
	this.Data["server"] = maps.Map{
		"id": params.ServerId,
	}
	this.Data["locationId"] = params.LocationId
	this.Data["rewriteId"] = params.RewriteId
	this.Data["fastcgiId"] = params.FastcgiId
	this.Data["backendId"] = params.BackendId

	headerList, err := server.FindHeaderList(params.LocationId, params.BackendId, params.RewriteId, params.FastcgiId)
	if err != nil {
		this.Fail(err.Error())
	}

	header := headerList.FindRequestHeader(params.HeaderId)
	if header == nil {
		this.Fail("找不到要修改的Header")
	}

	this.Data["header"] = header

	this.Show()
}

// 提交修改
func (this *UpdateRequestHeaderAction) RunPost(params struct {
	ServerId   string
	LocationId string
	RewriteId  string
	FastcgiId  string
	BackendId  string
	HeaderId   string

	On    bool
	Name  string
	Value string

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

	header := headerList.FindRequestHeader(params.HeaderId)
	if header == nil {
		this.Fail("找不到要修改的Header")
	}

	header.On = params.On
	header.Name = params.Name
	header.Value = params.Value

	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()

	this.Success()
}
