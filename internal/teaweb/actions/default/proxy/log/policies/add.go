package policies

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/tealogs"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/cmd"
)

type AddAction actions.Action

// 添加日志策略
func (this *AddAction) RunGet(params struct{}) {
	this.Data["selectedMenu"] = "add"

	this.Data["storages"] = tealogs.AllStorages()
	this.Data["formats"] = tealogs.AllStorageFormats()

	// 匹配条件运算符
	this.Data["condOperators"] = shared.AllRequestOperators()
	this.Data["condVariables"] = proxyutils.DefaultRequestVariables()

	// syslog
	this.Data["syslogPriorities"] = tealogs.SyslogStoragePriorities

	this.Show()
}

// 保存提交
func (this *AddAction) RunPost(params struct {
	Name            string
	StorageFormat   string
	StorageTemplate string
	StorageType     string

	// file
	FilePath       string
	FileAutoCreate bool

	// es
	EsEndpoint    string
	EsIndex       string
	EsMappingType string
	EsUsername    string
	EsPassword    string

	// mysql
	MysqlHost     string
	MysqlPort     int
	MysqlUsername string
	MysqlPassword string
	MysqlDatabase string
	MysqlTable    string
	MysqlLogField string

	// tcp
	TcpNetwork string
	TcpAddr    string

	// syslog
	SyslogProtocol   string
	SyslogServerAddr string
	SyslogServerPort int
	SyslogSocket     string
	SyslogTag        string
	SyslogPriority   int

	// command
	CommandCommand string
	CommandArgs    string
	CommandDir     string

	On bool

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入日志策略的名称").
		Field("storageType", params.StorageType).
		Require("请选择存储类型")

	var instance interface{} = nil
	switch params.StorageType {
	case tealogs.StorageTypeFile:
		params.Must.
			Field("filePath", params.FilePath).
			Require("请输入日志文件路径")

		storage := new(tealogs.FileStorage)
		storage.Format = params.StorageFormat
		storage.Template = params.StorageTemplate
		storage.Path = params.FilePath
		storage.AutoCreate = params.FileAutoCreate
		instance = storage
	case tealogs.StorageTypeES:
		params.Must.
			Field("esEndpoint", params.EsEndpoint).
			Require("请输入Endpoint").
			Field("esIndex", params.EsIndex).
			Require("请输入Index名称").
			Field("esMappingType", params.EsMappingType).
			Require("请输入Mapping名称")

		storage := new(tealogs.ESStorage)
		storage.Format = params.StorageFormat
		storage.Template = params.StorageTemplate
		storage.Endpoint = params.EsEndpoint
		storage.Index = params.EsIndex
		storage.MappingType = params.EsMappingType
		storage.Username = params.EsUsername
		storage.Password = params.EsPassword
		instance = storage
	case tealogs.StorageTypeMySQL:
		params.Must.
			Field("mysqlHost", params.MysqlHost).
			Require("请输入主机地址").
			Field("mysqlDatabase", params.MysqlDatabase).
			Require("请输入数据库名称").
			Field("mysqlTable", params.MysqlTable).
			Require("请输入数据表名称").
			Field("mysqlLogField", params.MysqlLogField).
			Require("请输入日志存储字段名")

		storage := new(tealogs.MySQLStorage)
		storage.AutoCreateTable = true
		storage.Format = params.StorageFormat
		storage.Template = params.StorageTemplate
		storage.Host = params.MysqlHost
		storage.Port = params.MysqlPort
		storage.Username = params.MysqlUsername
		storage.Password = params.MysqlPassword
		storage.Database = params.MysqlDatabase
		storage.Table = params.MysqlTable
		storage.LogField = params.MysqlLogField
		instance = storage
	case tealogs.StorageTypeTCP:
		params.Must.
			Field("tcpNetwork", params.TcpNetwork).
			Require("请选择网络协议").
			Field("tcpAddr", params.TcpAddr).
			Require("请输入网络地址")

		storage := new(tealogs.TCPStorage)
		storage.Format = params.StorageFormat
		storage.Template = params.StorageTemplate
		storage.Network = params.TcpNetwork
		storage.Addr = params.TcpAddr
		instance = storage
	case tealogs.StorageTypeSyslog:
		switch params.SyslogProtocol {
		case tealogs.SyslogStorageProtocolTCP, tealogs.SyslogStorageProtocolUDP:
			params.Must.
				Field("syslogServerAddr", params.SyslogServerAddr).
				Require("请输入网络地址")
		case tealogs.SyslogStorageProtocolSocket:
			params.Must.
				Field("syslogSocket", params.SyslogSocket).
				Require("请输入Socket路径")
		}

		storage := new(tealogs.SyslogStorage)
		storage.Format = params.StorageFormat
		storage.Template = params.StorageTemplate
		storage.Protocol = params.SyslogProtocol
		storage.ServerAddr = params.SyslogServerAddr
		storage.ServerPort = params.SyslogServerPort
		storage.Socket = params.SyslogSocket
		storage.Tag = params.SyslogTag
		storage.Priority = params.SyslogPriority
		instance = storage
	case tealogs.StorageTypeCommand:
		params.Must.
			Field("commandCommand", params.CommandCommand).
			Require("请输入可执行命令")

		storage := new(tealogs.CommandStorage)
		storage.Format = params.StorageFormat
		storage.Template = params.StorageTemplate
		storage.Command = params.CommandCommand
		storage.Args = cmd.ParseArgs(params.CommandArgs)
		storage.Dir = params.CommandDir
		instance = storage
	}

	if instance == nil {
		this.Fail("找不到选择的存储类型")
	}

	policy := teaconfigs.NewAccessLogStoragePolicy()
	policy.Type = params.StorageType
	policy.Name = params.Name
	policy.On = params.On

	options := map[string]interface{}{}
	err := teautils.ObjectToMapJSON(instance, &options)
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}
	policy.Options = options

	// 匹配条件
	conds, breakCond, err := proxyutils.ParseRequestConds(this.Request, "request")
	if err != nil {
		this.Fail("匹配条件\"" + breakCond.Param + " " + breakCond.Operator + " " + breakCond.Value + "\"校验失败：" + err.Error())
	}
	policy.Cond = conds

	err = policy.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	policyList := teaconfigs.SharedAccessLogStoragePolicyList()
	policyList.AddId(policy.Id)
	err = policyList.Save()
	if err != nil {
		_ = policy.Delete()
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
