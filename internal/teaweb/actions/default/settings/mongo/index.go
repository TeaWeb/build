package mongo

import (
	"github.com/TeaWeb/build/internal/teaconfigs/db"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/maps"
	"runtime"
	"strings"
)

type IndexAction actions.Action

// MongoDB连接信息
func (this *IndexAction) Run(params struct{}) {
	this.Data["shouldRestart"] = shouldRestart

	config, err := db.LoadMongoConfig()
	if err != nil {
		this.Data["error"] = err.Error()
		config = db.DefaultMongoConfig()
	}

	this.Data["config"] = maps.Map{
		"scheme":                  config.Scheme,
		"authEnabled":             config.AuthEnabled,
		"username":                config.Username,
		"password":                strings.Repeat("*", len(config.Password)),
		"host":                    config.Host(),
		"port":                    config.Port(),
		"authMechanism":           config.AuthMechanism,
		"authMechanismProperties": config.AuthMechanismPropertiesString(),
		"requestURI":              config.RequestURI,
	}
	this.Data["uri"] = config.ComposeURIMask(true)

	// 连接状态
	err = teadb.SharedDB().Test()
	if err != nil {
		this.Data["error"] = err.Error()
	} else {
		this.Data["error"] = ""
	}

	// 检测是否已安装
	mongodbPath := Tea.Root + "/mongodb/bin/mongod"
	if files.NewFile(mongodbPath).Exists() {
		this.Data["isInstalled"] = true
	} else {
		this.Data["isInstalled"] = false
	}

	// 当前系统
	this.Data["os"] = runtime.GOOS

	this.Show()
}
