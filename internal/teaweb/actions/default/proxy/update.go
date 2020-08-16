package proxy

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/types"
	"strconv"
)

type UpdateAction actions.Action

// 修改代理服务信息
func (this *UpdateAction) Run(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	this.Data["server"] = server
	this.Data["selectedTab"] = "basic"
	this.Data["isTCP"] = server.IsTCP()
	this.Data["isForwardHTTP"] = server.ForwardHTTP

	if server.Gzip == nil {
		server.Gzip = &teaconfigs.GzipConfig{
			Level:     0,
			MinLength: "",
			MimeTypes: nil,
		}
	}

	this.Data["usualCharsets"] = teautils.UsualCharsets
	this.Data["charsets"] = teautils.AllCharsets

	this.Data["accessLogs"] = proxyutils.FormatAccessLog(server.AccessLog)

	// 通知设置
	server.SetupNoticeItems()
	this.Data["noticeItems"] = server.NoticeItems

	// 目录补全
	security := configs.SharedAdminConfig().Security
	if security != nil {
		this.Data["dirAutoComplete"] = security.DirAutoComplete
	} else {
		this.Data["dirAutoComplete"] = false
	}

	this.Show()
}

// 保存提交
func (this *UpdateAction) RunPost(params struct {
	// 通用
	ServerId    string
	Description string
	Name        []string
	Listen      []string

	// HTTP
	HttpOn      bool
	Root        string
	Charset     string
	Index       []string
	MaxBodySize float64
	MaxBodyUnit string

	EnableStat bool

	// gzip
	GzipLevel          uint8
	GzipMinLength      float64
	GzipMinUnit        string
	GzipMimeTypeValues []string

	// cache
	CacheStatic bool

	// pages
	PageStatusList    []string
	PageURLList       []string
	PageNewStatusList []string

	ShutdownPageOn     bool
	ShutdownPageURL    string
	ShutdownPageStatus int

	RedirectToHttps bool

	// TCP
	TcpOn              bool
	TcpReadBufferSize  int
	TcpWriteBufferSize int

	// ForwardHTTP
	ForwardHTTPEnableMITM bool

	Must *actions.Must
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	params.Must.
		Field("description", params.Description).
		Require("代理服务名称不能为空")

	server.Description = params.Description
	server.Name = params.Name
	server.Listen = teautils.FormatAddressList(params.Listen)

	if server.TCP != nil { // TCP
		server.TCP.TCPOn = params.TcpOn
		if params.TcpReadBufferSize >= 0 {
			server.TCP.ReadBufferSize = params.TcpReadBufferSize
		}
		if params.TcpWriteBufferSize >= 0 {
			server.TCP.WriteBufferSize = params.TcpWriteBufferSize
		}
	} else if server.ForwardHTTP != nil { // ForwardHTTP
		server.Http = params.HttpOn

		server.ForwardHTTP.EnableMITM = params.ForwardHTTPEnableMITM

		// 访问日志
		server.AccessLog = proxyutils.ParseAccessLogForm(this.Request)

		server.DisableStat = !params.EnableStat
	} else { // HTTP
		server.Http = params.HttpOn
		server.Root = params.Root
		server.Charset = params.Charset
		server.Index = params.Index
		server.MaxBodySize = strconv.FormatFloat(params.MaxBodySize, 'f', -1, 64) + params.MaxBodyUnit

		// 访问日志
		server.AccessLog = proxyutils.ParseAccessLogForm(this.Request)

		server.DisableStat = !params.EnableStat
		server.CacheStatic = params.CacheStatic

		// gzip
		if params.GzipLevel > 0 && params.GzipLevel <= 9 {
			minLength := strconv.FormatFloat(params.GzipMinLength, 'f', -1, 64) + params.GzipMinUnit
			gzip := &teaconfigs.GzipConfig{
				Level:     int8(params.GzipLevel),
				MinLength: minLength,
				MimeTypes: params.GzipMimeTypeValues,
			}
			server.Gzip = gzip
		} else {
			server.Gzip = nil
		}

		server.Pages = []*teaconfigs.PageConfig{}
		for index, status := range params.PageStatusList {
			page := teaconfigs.NewPageConfig()
			page.Status = []string{status}
			if index < len(params.PageURLList) {
				page.URL = params.PageURLList[index]
			}
			if index < len(params.PageNewStatusList) {
				page.NewStatus = types.Int(params.PageNewStatusList[index])
				if page.NewStatus < 0 {
					page.NewStatus = 0
				}
			}
			server.AddPage(page)
		}

		// shutdown page
		if server.Shutdown != nil {
			server.Shutdown.On = params.ShutdownPageOn
			server.Shutdown.URL = params.ShutdownPageURL
			server.Shutdown.Status = params.ShutdownPageStatus
		} else if params.ShutdownPageOn {
			server.Shutdown = teaconfigs.NewShutdownConfig()
			server.Shutdown.On = params.ShutdownPageOn
			server.Shutdown.URL = params.ShutdownPageURL
			server.Shutdown.Status = params.ShutdownPageStatus
		}
		if server.Shutdown != nil && server.Shutdown.On && len(server.Shutdown.URL) == 0 {
			this.FailField("shutdownPageURL", "请输入临时关闭页面文件路径")
		}

		server.RedirectToHttps = params.RedirectToHttps
	}

	// 通知设置
	server.SetupNoticeItemsFromRequest(this.Request)

	err := server.Validate()
	if err != nil {
		this.Fail("校验失败：" + err.Error())
	}

	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 重启
	proxyutils.NotifyChange()

	this.Success()
}
