package groups

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type AddAction actions.Action

// 添加分组
func (this *AddAction) Run(params struct {
	ServerId   string
	LocationId string
	Websocket  int
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	if len(params.LocationId) > 0 {
		this.Data["selectedTab"] = "location"
	} else {
		this.Data["selectedTab"] = "backend"
	}
	this.Data["server"] = server

	this.Data["isTCP"] = server.IsTCP()
	this.Data["isHTTP"] = server.IsHTTP()

	this.Data["locationId"] = params.LocationId
	this.Data["websocket"] = params.Websocket

	this.Data["operators"] = shared.AllRequestOperators()

	// 请求变量
	this.Data["variables"] = proxyutils.DefaultRequestVariables()

	this.Show()
}

// 提交保存
func (this *AddAction) RunPost(params struct {
	ServerId string

	Name string

	IPRangeTypeList     []string `alias:"ipRangeTypeList"`
	IPRangeFromList     []string `alias:"ipRangeFromList"`
	IPRangeToList       []string `alias:"ipRangeToList"`
	IPRangeCIDRIPList   []string `alias:"ipRangeCIDRIPList"`
	IPRangeCIDRBitsList []string `alias:"ipRangeCIDRBitsList"`
	IPRangeVarList      []string `alias:"ipRangeVarList"`

	RequestHeaderNames  []string
	RequestHeaderValues []string

	ResponseHeaderNames  []string
	ResponseHeaderValues []string

	Must *actions.Must
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到要修改的Server")
	}

	params.Must.
		Field("name", params.Name).
		Require("请输入分组名")

	group := teaconfigs.NewRequestGroup()
	group.Name = params.Name

	// 匹配条件
	conds, breakCond, err := proxyutils.ParseRequestConds(this.Request, "request")
	if err != nil {
		this.Fail("匹配条件\"" + breakCond.Param + " " + breakCond.Operator + " " + breakCond.Value + "\"校验失败：" + err.Error())
	}
	group.Cond = conds

	// IP范围
	if len(params.IPRangeTypeList) > 0 {
		for index, ipRangeType := range params.IPRangeTypeList {
			if index < len(params.IPRangeFromList) && index < len(params.IPRangeToList) && index < len(params.IPRangeCIDRIPList) && index < len(params.IPRangeCIDRBitsList) && index < len(params.IPRangeVarList) {
				if ipRangeType == "range" {
					config := shared.NewIPRangeConfig()
					config.Type = shared.IPRangeTypeRange
					config.IPFrom = params.IPRangeFromList[index]
					config.IPTo = params.IPRangeToList[index]
					config.Param = params.IPRangeVarList[index]
					err := config.Validate()
					if err != nil {
						this.Fail("校验失败：" + err.Error())
					}
					group.AddIPRange(config)
				} else if ipRangeType == "cidr" {
					config := shared.NewIPRangeConfig()
					config.Type = shared.IPRangeTypeCIDR
					config.CIDR = params.IPRangeCIDRIPList[index] + "/" + params.IPRangeCIDRBitsList[index]
					config.Param = params.IPRangeVarList[index]
					err := config.Validate()
					if err != nil {
						this.Fail("校验失败：" + err.Error())
					}
					group.AddIPRange(config)
				}
			}
		}
	}

	// 请求Header
	if len(params.RequestHeaderNames) > 0 {
		for index, headerName := range params.RequestHeaderNames {
			if index < len(params.RequestHeaderValues) {
				header := shared.NewHeaderConfig()
				header.Name = headerName
				header.Value = params.RequestHeaderValues[index]
				group.AddRequestHeader(header)
			}
		}
	}

	// 响应Header
	if len(params.ResponseHeaderNames) > 0 {
		for index, headerName := range params.ResponseHeaderNames {
			if index < len(params.ResponseHeaderValues) {
				header := shared.NewHeaderConfig()
				header.Name = headerName
				header.Value = params.ResponseHeaderValues[index]
				group.AddResponseHeader(header)
			}
		}
	}

	// 保存
	server.AddRequestGroup(group)
	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 通知改变
	proxyutils.NotifyChange()

	this.Success()
}
