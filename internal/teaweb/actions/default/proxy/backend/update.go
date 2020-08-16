package backend

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/certs/certutils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"regexp"
	"strings"
)

type UpdateAction actions.Action

// 修改后端服务器
func (this *UpdateAction) Run(params struct {
	ServerId   string
	LocationId string
	Websocket  bool
	Backend    string
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

	this.Data["isTCP"] = server.IsTCP()
	this.Data["isHTTP"] = server.IsHTTP()

	this.Data["locationId"] = params.LocationId
	this.Data["websocket"] = types.Int(params.Websocket)
	this.Data["from"] = params.From

	backendList, err := server.FindBackendList(params.LocationId, params.Websocket)
	if err != nil {
		this.Fail(err.Error())
	}
	backend := backendList.FindBackend(params.Backend)
	if backend == nil {
		this.Fail("找不到要修改的后端服务器")
	}

	err = backend.Validate()
	if err != nil {
		logs.Error(err)
	}

	if len(backend.RequestGroupIds) == 0 {
		backend.AddRequestGroupId("default")
	}

	if len(backend.RequestURI) == 0 {
		backend.RequestURI = "${requestURI}"
	}

	if backend.FTP == nil {
		backend.FTP = &teaconfigs.FTPBackendConfig{}
	}

	this.Data["backend"] = maps.Map{
		"id":              backend.Id,
		"address":         backend.Address,
		"scheme":          backend.Scheme,
		"code":            backend.Code,
		"weight":          backend.Weight,
		"failTimeout":     int(backend.FailTimeoutDuration().Seconds()),
		"readTimeout":     int(backend.ReadTimeoutDuration().Seconds()),
		"idleTimeout":     strings.TrimSuffix(backend.IdleTimeout, "s"),
		"on":              backend.On,
		"maxConns":        backend.MaxConns,
		"idleConns":       backend.IdleConns,
		"maxFails":        backend.MaxFails,
		"isDown":          backend.IsDown,
		"isBackup":        backend.IsBackup,
		"requestGroupIds": backend.RequestGroupIds,
		"requestURI":      backend.RequestURI,
		"checkOn":         backend.CheckOn,
		"checkURL":        backend.CheckURL,
		"checkInterval":   backend.CheckInterval,
		"checkTimeout":    strings.TrimSuffix(backend.CheckTimeout, "s"),
		"requestHeaders":  backend.RequestHeaders,
		"responseHeaders": backend.ResponseHeaders,
		"host":            backend.Host,
		"cert":            backend.Cert,
		"ftp":             backend.FTP,
	}

	// 公共可以使用的证书
	this.Data["sharedCerts"] = certutils.ListPairCertsMap()

	this.Show()
}

// 提交
func (this *UpdateAction) RunPost(params struct {
	ServerId   string
	LocationId string
	Websocket  bool
	BackendId  string
	Address    string
	Scheme     string

	UseCert        bool
	CertId         string
	CertServerName string

	Weight          uint
	On              bool
	Code            string
	FailTimeout     uint
	ReadTimeout     uint
	IdleTimeout     string
	MaxFails        int32
	MaxConns        int32
	IdleConns       int32
	IsBackup        bool
	RequestGroupIds []string
	RequestURI      string

	CheckOn       bool
	CheckURL      string
	CheckInterval int
	CheckTimeout  string

	FtpUsername string
	FtpPassword string
	FtpDir      string

	RequestHeaderNames  []string
	RequestHeaderValues []string

	ResponseHeaderNames  []string
	ResponseHeaderValues []string

	Host string

	Must *actions.Must
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	params.Must.
		Field("address", params.Address).
		Require("请输入后端服务器地址")

	// 健康检查
	if params.CheckOn {
		if server.IsHTTP() {
			if len(params.CheckURL) == 0 {
				this.FailField("checkURL", "健康检查URL不能为空")
			}

			if !regexp.MustCompile("(?i)(http://|https://)").MatchString(params.CheckURL) {
				this.FailField("checkURL", "健康检查URL必须以http://或https://开头")
			}
		}
	}

	// 证书
	if params.UseCert {
		if len(params.CertId) == 0 {
			this.Fail("在请求后端服务器时使用SSL证书时，需要选择一个证书")
		}

		cert := teaconfigs.SharedSSLCertList().FindCert(params.CertId)
		if cert == nil {
			this.Fail("选择的SSL证书不存在")
		}
	}

	backendList, err := server.FindBackendList(params.LocationId, params.Websocket)
	if err != nil {
		this.Fail(err.Error())
	}

	backend := backendList.FindBackend(params.BackendId)
	if backend == nil {
		this.Fail("找不到要修改的后端服务器")
	}
	backend.Touch()
	backend.TeaVersion = teaconst.TeaVersion

	backend.Address = teautils.FormatAddress(params.Address)
	backend.Scheme = params.Scheme

	// 证书
	if params.UseCert {
		backend.Cert = teaconfigs.NewSSLCertConfig("", "")
		backend.Cert.IsShared = true
		backend.Cert.Id = params.CertId
		backend.Cert.ServerName = params.CertServerName
	} else {
		backend.Cert = nil
	}

	backend.Weight = params.Weight
	backend.On = params.On
	backend.IsDown = false
	backend.Code = params.Code
	backend.FailTimeout = fmt.Sprintf("%d", params.FailTimeout) + "s"
	backend.ReadTimeout = fmt.Sprintf("%d", params.ReadTimeout) + "s"
	backend.MaxFails = params.MaxFails
	backend.MaxConns = params.MaxConns
	backend.IsBackup = params.IsBackup
	backend.RequestGroupIds = params.RequestGroupIds
	backend.RequestURI = params.RequestURI

	backend.CheckOn = params.CheckOn
	backend.CheckURL = params.CheckURL
	backend.CheckInterval = params.CheckInterval
	backend.CheckTimeout = params.CheckTimeout + "s"

	backend.IdleConns = params.IdleConns
	backend.IdleTimeout = params.IdleTimeout + "s"

	// 请求Header
	backend.RequestHeaders = []*shared.HeaderConfig{}
	if len(params.RequestHeaderNames) > 0 {
		for index, headerName := range params.RequestHeaderNames {
			if index < len(params.RequestHeaderValues) {
				header := shared.NewHeaderConfig()
				header.Name = headerName
				header.Value = params.RequestHeaderValues[index]
				backend.AddRequestHeader(header)
			}
		}
	}

	// 响应Header
	backend.ResponseHeaders = []*shared.HeaderConfig{}
	if len(params.ResponseHeaderNames) > 0 {
		for index, headerName := range params.ResponseHeaderNames {
			if index < len(params.ResponseHeaderValues) {
				header := shared.NewHeaderConfig()
				header.Name = headerName
				header.Value = params.ResponseHeaderValues[index]
				backend.AddResponseHeader(header)
			}
		}
	}

	backend.Host = params.Host

	// ftp
	if params.Scheme == "ftp" {
		backend.FTP = &teaconfigs.FTPBackendConfig{
			Username: params.FtpUsername,
			Password: params.FtpPassword,
			Dir:      params.FtpDir,
		}
	}

	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()

	this.Success()
}
