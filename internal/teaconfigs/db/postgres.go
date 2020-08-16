package db

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"gopkg.in/yaml.v3"
	"github.com/iwind/TeaGo/Tea"
	"io/ioutil"
)

const (
	postgresFilename = "postgres.conf"
)

// PostgreSQL配置
type PostgresConfig struct {
	DSN string `yaml:"dsn" json:"dsn"`

	Addr     string `yaml:"addr" json:"addr"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	DBName   string `yaml:"dbName" json:"dbName"`

	PoolSize int `yaml:"poolSize" json:"poolSize"`

	// 日志访问配置
	AccessLog *PostgresAccessLogConfig `yaml:"accessLog" json:"accessLog"`
}

// 访问日志配置
type PostgresAccessLogConfig struct {
	CleanHour int `yaml:"cleanHour" json:"cleanHour"` // 清理时间，0-23
	KeepDays  int `yaml:"keepDays" json:"keepDays"`   // 保留挺熟
}

// 获取新对象
func NewPostgresConfig() *PostgresConfig {
	return &PostgresConfig{}
}

// 默认的PostgreSQL配置
func DefaultPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Addr:     "127.0.0.1:5432",
		Username: "postgres",
		DBName:   "teaweb",
	}
}

// 加载PostgreSQL配置
func LoadPostgresConfig() (*PostgresConfig, error) {
	data, err := ioutil.ReadFile(Tea.ConfigFile(postgresFilename))
	if err != nil {
		return nil, err
	}
	config := &PostgresConfig{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// 组合DSN
func (this *PostgresConfig) ComposeDSN() string {
	return "postgres://" + this.Username + ":" + this.Password + "@" + this.Addr + "/" + this.DBName + "?sslmode=disable"
}

// 保存
func (this *PostgresConfig) Save() error {
	shared.Locker.Lock()
	defer shared.Locker.WriteUnlockNotify()

	data, err := yaml.Marshal(this)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(Tea.ConfigFile(postgresFilename), data, 0777)
}
