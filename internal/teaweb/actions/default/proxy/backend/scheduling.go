package backend

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/scheduling"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

type SchedulingAction actions.Action

// 调度算法
func (this *SchedulingAction) Run(params struct {
	ServerId   string
	LocationId string
	Websocket  bool
	From       string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	this.Data["server"] = server
	if len(params.LocationId) > 0 {
		this.Data["selectedTab"] = "location"
	} else {
		this.Data["selectedTab"] = "backend"
	}
	this.Data["locationId"] = params.LocationId
	this.Data["websocket"] = types.Int(params.Websocket)
	this.Data["from"] = params.From

	backendList, err := server.FindBackendList(params.LocationId, params.Websocket)
	if err != nil {
		this.Fail(err.Error())
	}
	if backendList.SchedulingConfig() == nil {
		backendList.SetSchedulingConfig(&teaconfigs.SchedulingConfig{
			Code:    "random",
			Options: maps.Map{},
		})
	}
	this.Data["scheduling"] = backendList.SchedulingConfig()

	// 调度类型
	schedulingTypes := []maps.Map{}
	for _, m := range scheduling.AllSchedulingTypes() {
		networks, ok := m["networks"]
		if !ok {
			continue
		}
		if !types.IsSlice(networks) {
			continue
		}
		if (server.IsHTTP() && lists.Contains(networks, "http")) || (server.IsTCP() && lists.Contains(networks, "tcp")) {
			schedulingTypes = append(schedulingTypes, m)
		}
	}
	this.Data["schedulingTypes"] = schedulingTypes

	this.Show()
}

// 保存提交
func (this *SchedulingAction) RunPost(params struct {
	ServerId    string
	LocationId  string
	Websocket   bool
	Type        string
	HashKey     string
	StickyType  string
	StickyParam string
	Must        *actions.Must
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	backendList, err := server.FindBackendList(params.LocationId, params.Websocket)
	if err != nil {
		this.Fail(err.Error())
	}

	options := maps.Map{}
	if params.Type == "hash" {
		params.Must.
			Field("hashKey", params.HashKey).
			Require("请输入Key")

		options["key"] = params.HashKey
	} else if params.Type == "sticky" {
		params.Must.
			Field("stickyType", params.StickyType).
			Require("请选择参数类型").
			Field("stickyParam", params.StickyParam).
			Require("请输入参数名").
			Match("^[a-zA-Z0-9]+$", "参数名只能是英文字母和数字的组合").
			MaxCharacters(50, "参数名长度不能超过50位")

		options["type"] = params.StickyType
		options["param"] = params.StickyParam
	}

	if scheduling.FindSchedulingType(params.Type) == nil {
		this.Fail("不支持此种算法")
	}

	backendList.SetSchedulingConfig(&teaconfigs.SchedulingConfig{
		Code:    params.Type,
		Options: options,
	})

	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	if len(backendList.AllBackends()) > 0 {
		proxyutils.NotifyChange()
	}

	this.Success()
}
