package db

import (
	"gopkg.in/yaml.v3"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"io/ioutil"
)

// 数据库配置
type DBType = string

const (
	DBConfigFile = "db.conf"

	DBTypeMongo    = "mongo"
	DBTypeMySQL    = "mysql"
	DBTypePostgres = "postgres"
)

// 变量
var sharedDBConfig *DBConfig = nil

// 数据库配置
type DBConfig struct {
	Type          DBType `yaml:"type" json:"type"`                   // 类型：mongo, mysql, postgres ...
	IsInitialized bool   `yaml:"isInitialized" json:"isInitialized"` // 是否初始化
}

// 取得共享的配置
func SharedDBConfig() *DBConfig {
	if sharedDBConfig != nil {
		return sharedDBConfig
	}
	config := &DBConfig{
		Type:          DBTypeMongo,
		IsInitialized: true,
	}
	data, err := ioutil.ReadFile(Tea.ConfigFile(DBConfigFile))
	if err != nil {
		return config
	}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		logs.Error(err)
		return config
	}
	sharedDBConfig = config
	return config
}

// 保存
func (this *DBConfig) Save() error {
	data, err := yaml.Marshal(this)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(Tea.ConfigFile(DBConfigFile), data, 0777)
}

// 数据库名称
func (this *DBConfig) TypeName() string {
	switch this.Type {
	case DBTypeMongo:
		return "MongoDB"
	case DBTypeMySQL:
		return "MySQL"
	case DBTypePostgres:
		return "PostgreSQL"
	}
	return ""
}
