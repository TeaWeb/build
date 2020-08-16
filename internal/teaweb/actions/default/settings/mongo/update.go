package mongo

import (
	"github.com/TeaWeb/build/internal/teaconfigs/db"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateAction actions.Action

// 修改连接
func (this *UpdateAction) Run(params struct {
	Action string
}) {
	this.Data["action"] = params.Action
	this.Data["typeIsChanged"] = db.SharedDBConfig().Type != db.DBTypeMongo

	config, err := db.LoadMongoConfig()
	if err != nil {
		config = db.DefaultMongoConfig()
	}
	this.Data["config"] = maps.Map{
		"scheme":                  config.Scheme,
		"username":                config.Username,
		"password":                "",
		"host":                    config.Host(),
		"port":                    config.Port(),
		"dbName":                  config.DBName,
		"authEnabled":             config.AuthEnabled,
		"authMechanism":           config.AuthMechanism,
		"authMechanismProperties": config.AuthMechanismPropertiesString(),
		"requestURI":              config.RequestURI,
		"poolSize":                config.PoolSize,
		"timeout":                 config.Timeout,
	}

	this.Show()
}

// 提交保存
func (this *UpdateAction) RunPost(params struct {
	Host                    string
	Port                    uint
	DBName                  string `alias:"dbName"`
	Username                string
	Password                string
	AuthEnabled             bool
	AuthMechanism           string
	AuthMechanismProperties string
	PoolSize                int
	Timeout                 int

	Must *actions.Must
}) {
	// 是否已改变
	sharedConfig := db.SharedDBConfig()
	isChanged := sharedConfig.Type != db.DBTypeMongo

	params.Must.
		Field("host", params.Host).
		Require("请输入主机地址").
		Field("port", params.Port).
		Require("请输入端口").
		Gt(0, "请输入正确的端口").
		Field("poolSize", params.PoolSize).
		Gte(0, "连接池不能小于0").
		Field("timeout", params.Timeout).
		Gte(0, "超时时间不能小于0")

	config, err := db.LoadMongoConfig()
	if err != nil {
		this.Fail(err.Error())
	}

	config.SetAddr(params.Host, params.Port)
	config.DBName = params.DBName
	config.AuthEnabled = params.AuthEnabled
	config.AuthMechanism = params.AuthMechanism
	config.LoadAuthMechanismProperties(params.AuthMechanismProperties)
	config.Username = params.Username
	if len(params.Username) > 0 {
		if len(params.Password) > 0 {
			config.Password = params.Password
		}
	} else {
		config.Password = ""
	}
	config.PoolSize = params.PoolSize
	config.Timeout = params.Timeout
	err = config.Save()

	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	if isChanged {
		shouldRestart = true

		sharedConfig.Type = db.DBTypeMongo
		err = sharedConfig.Save()
		if err != nil {
			this.Fail("保存失败：" + err.Error())
		}
	}

	// 重新连接
	teadb.ChangeDB()

	this.Next("/settings/mongo", nil).Success("保存成功")
}
