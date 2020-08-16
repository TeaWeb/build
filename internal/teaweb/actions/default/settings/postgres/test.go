package postgres

import (
	"github.com/TeaWeb/build/internal/teaconfigs/db"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/actions"
)

type TestAction actions.Action

// 测试数据库
func (this *TestAction) RunPost(params struct {
	DBType     string `alias:"dbType"`
	Addr       string
	Username   string
	Password   string
	DBName     string `alias:"dbName"`
	AutoCreate bool
}) {
	params.Addr = teautils.FormatAddress(params.Addr)

	if len(params.Addr) == 0 {
		this.Fail("请输入数据库主机地址")
	}

	if len(params.DBName) == 0 {
		this.Fail("请输入数据库名称")
	}

	config := db.NewPostgresConfig()
	config.Addr = params.Addr
	config.Username = params.Username
	config.Password = params.Password
	config.DBName = params.DBName
	dsn := config.ComposeDSN()
	driver := new(teadb.PostgresDriver)
	message, ok := driver.TestDSN(dsn, params.AutoCreate)
	if ok {
		this.Success()
	}

	this.Fail(message)
}
