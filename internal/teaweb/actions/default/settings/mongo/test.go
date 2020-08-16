package mongo

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs/db"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/iwind/TeaGo/actions"
)

type TestAction actions.Action

func (this *TestAction) Run(params struct {
	Host                    string
	Port                    uint
	DBName                  string `alias:"dbName"`
	Username                string
	Password                string
	AuthEnabled             bool
	AuthMechanism           string
	AuthMechanismProperties string
}) {
	config := db.NewMongoConfig()
	config.Addr = params.Host
	if params.Port > 0 {
		config.Addr += ":" + fmt.Sprintf("%d", params.Port)
	}
	config.DBName = params.DBName
	config.AuthEnabled = params.AuthEnabled
	config.Username = params.Username
	config.Password = params.Password
	config.AuthMechanism = params.AuthMechanism
	config.LoadAuthMechanismProperties(params.AuthMechanismProperties)

	oldConfig, err := db.LoadMongoConfig()
	if err != nil {
		this.Fail(err.Error())
	}
	if oldConfig != nil && len(config.Password) == 0 {
		config.Password = oldConfig.Password
	}

	driver := new(teadb.MongoDriver)
	message, ok := driver.TestConfig(config)
	if !ok {
		this.Fail(message)
	}

	this.Success()
}
