package install

import (
	"github.com/TeaWeb/build/internal/teaconfigs/db"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/iwind/TeaGo/actions"
)

type SaveAction actions.Action

func (this *SaveAction) RunPost(params struct {
	DBType   string `alias:"dbType"`
	Addr     string
	Username string
	Password string
	DBName   string `alias:"dbName"`

	// mongodb专用
	AuthEnabled             bool
	AuthMechanism           string
	AuthMechanismProperties string
}) {
	switch params.DBType {
	case db.DBTypeMongo:
		config := db.NewMongoConfig()
		config.Addr = params.Addr
		config.Username = params.Username
		config.Password = params.Password
		config.DBName = params.DBName
		config.AuthEnabled = params.AuthEnabled
		config.AuthMechanism = params.AuthMechanism
		config.LoadAuthMechanismProperties(params.AuthMechanismProperties)
		config.URI = config.ComposeURI()
		err := config.Save()
		if err != nil {
			this.Fail("保存失败：" + err.Error())
		}

		sharedDB := db.SharedDBConfig()
		sharedDB.Type = db.DBTypeMongo
		sharedDB.IsInitialized = true
		err = sharedDB.Save()
		if err != nil {
			this.Fail("保存失败：" + err.Error())
		}

		// 切换数据库驱动
		teadb.ChangeDB()

		this.Success()
	case db.DBTypeMySQL:
		config := db.NewMySQLConfig()
		config.Addr = params.Addr
		config.Username = params.Username
		config.Password = params.Password
		config.DBName = params.DBName
		config.DSN = config.ComposeDSN()
		err := config.Save()
		if err != nil {
			this.Fail("保存失败：" + err.Error())
		}

		sharedDB := db.SharedDBConfig()
		sharedDB.Type = db.DBTypeMySQL
		sharedDB.IsInitialized = true
		err = sharedDB.Save()
		if err != nil {
			this.Fail("保存失败：" + err.Error())
		}

		// 切换数据库驱动
		teadb.ChangeDB()

		this.Success()
	case db.DBTypePostgres:
		config := db.NewPostgresConfig()
		config.Addr = params.Addr
		config.Username = params.Username
		config.Password = params.Password
		config.DBName = params.DBName
		config.DSN = config.ComposeDSN()
		err := config.Save()
		if err != nil {
			this.Fail("保存失败：" + err.Error())
		}

		sharedDB := db.SharedDBConfig()
		sharedDB.Type = db.DBTypePostgres
		sharedDB.IsInitialized = true
		err = sharedDB.Save()
		if err != nil {
			this.Fail("保存失败：" + err.Error())
		}

		// 切换数据库驱动
		teadb.ChangeDB()

		this.Success()
	}
	this.Success()
}
