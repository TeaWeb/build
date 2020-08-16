package headers

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type AddIgnoreAction actions.Action

// 添加屏蔽的Header
func (this *AddIgnoreAction) Run(params struct {
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
func (this *AddIgnoreAction) RunPost(params struct {
	ServerId   string
	LocationId string
	RewriteId  string
	FastcgiId  string
	BackendId  string
	Name       string
	Must       *actions.Must
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	params.Must.
		Field("name", params.Name).
		Require("请输入Name")

	headerList, err := server.FindHeaderList(params.LocationId, params.BackendId, params.RewriteId, params.FastcgiId)
	if err != nil {
		this.Fail(err.Error())
	}

	if headerList.ContainsIgnoreResponseHeader(params.Name) {
		this.Fail("已经存在，不需要重复添加")
	}

	headerList.AddIgnoreResponseHeader(params.Name)

	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()

	this.Success()
}
