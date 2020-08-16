package db

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/iwind/TeaGo/Tea"
	_ "github.com/iwind/TeaGo/bootstrap"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/types"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net"
	"net/url"
	"regexp"
	"strings"
)

const (
	mongoFilename = "mongo.conf"
)

// MongoDB配置
type MongoConfig struct {
	TeaVersion string `yaml:"teaVersion" json:"teaVersion"`

	URI string `yaml:"uri" json:"uri"`

	Scheme                  string             `yaml:"scheme" json:"scheme"`
	Username                string             `yaml:"username" json:"username"`
	Password                string             `yaml:"password" json:"password"`
	Addr                    string             `yaml:"addr" json:"addr"`
	AuthEnabled             bool               `yaml:"authEnabled" json:"authEnabled"`
	AuthMechanism           string             `yaml:"authMechanism" json:"authMechanism"`
	AuthMechanismProperties []*shared.Variable `yaml:"authMechanismProperties" json:"authMechanismProperties"`
	RequestURI              string             `yaml:"requestURI" json:"requestURI"` // @TODO 未来版本需要实现
	DBName                  string             `yaml:"dbName" json:"dbName"`

	PoolSize int `yaml:"poolSize" json:"poolSize"` // 连接池大小
	Timeout  int `yaml:"timeout" json:"timeout"`   // 超时时间（秒）

	// 日志访问配置
	AccessLog *MongoAccessLogConfig `yaml:"accessLog" json:"accessLog"`
}

// 访问日志配置
type MongoAccessLogConfig struct {
	CleanHour int `yaml:"cleanHour" json:"cleanHour"` // 清理时间，0-23
	KeepDays  int `yaml:"keepDays" json:"keepDays"`   // 保留挺熟
}

// 获取新对象
func NewMongoConfig() *MongoConfig {
	return &MongoConfig{}
}

// 加载MongoDB配置
func LoadMongoConfig() (*MongoConfig, error) {
	data, err := ioutil.ReadFile(Tea.ConfigFile(mongoFilename))
	if err != nil {
		return nil, err
	}
	config := &MongoConfig{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	// 老版本的处理
	if len(config.TeaVersion) == 0 {
		err = config.ParseFromURI(config.URI)
		if err != nil {
			logs.Error(err)
		}
		config.AuthEnabled = len(config.Username) > 0
		config.DBName = "teaweb"
	}

	return config, nil
}

// 默认的MongoDB配置
func DefaultMongoConfig() *MongoConfig {
	return &MongoConfig{
		Addr: "127.0.0.1:27017",
	}
}

// 从URI中分析配置
func (this *MongoConfig) ParseFromURI(uri string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}

	this.Scheme = u.Scheme
	this.Addr = u.Host

	if u.User != nil {
		this.AuthEnabled = true
		this.Username = u.User.Username()
		this.Password, _ = u.User.Password()
	}

	this.AuthMechanism = u.Query().Get("authMechanism")
	properties := u.Query().Get("authMechanismProperties")
	if len(properties) > 0 {
		for _, property := range strings.Split(properties, ",") {
			if strings.Contains(property, ":") {
				pieces := strings.Split(property, ":")
				this.AuthMechanismProperties = append(this.AuthMechanismProperties, shared.NewVariable(pieces[0], pieces[1]))
			}
		}
	}

	return nil
}

// 保存
func (this *MongoConfig) Save() error {
	shared.Locker.Lock()
	defer shared.Locker.WriteUnlockNotify()

	this.URI = this.ComposeURI()
	this.TeaVersion = teaconst.TeaVersion
	data, err := yaml.Marshal(this)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(Tea.ConfigFile(mongoFilename), data, 0666)
}

// 组合后的URI
func (this *MongoConfig) ComposeURI() string {
	return this.ComposeURIMask(false)
}

// 组合后的URI
func (this *MongoConfig) ComposeURIMask(mask bool) string {
	uri := ""
	if len(this.Scheme) > 0 {
		uri += this.Scheme + "://"
	} else {
		uri += "mongodb://"
	}

	if this.AuthEnabled && len(this.Username) > 0 {
		uri += this.Username
		if len(this.Password) > 0 {
			if mask {
				uri += ":" + strings.Repeat("*", len(this.Password))
			} else {
				uri += ":" + this.Password
			}
		}
		uri += "@"
	}

	uri += this.Addr

	if this.AuthEnabled && len(this.AuthMechanism) > 0 {
		uri += "/?authMechanism=" + this.AuthMechanism

		if len(this.AuthMechanismProperties) > 0 {
			properties := []string{}
			for _, v := range this.AuthMechanismProperties {
				properties = append(properties, v.Name+":"+v.Value)
			}
			uri += "&authMechanismProperties=" + strings.Join(properties, ",")
		}
	}

	return uri
}

// 分析认证参数
func (this *MongoConfig) LoadAuthMechanismProperties(properties string) {
	if len(properties) == 0 {
		this.AuthMechanismProperties = []*shared.Variable{}
		return
	}
	for _, property := range regexp.MustCompile("\\s*,\\s*").Split(properties, -1) {
		if strings.Contains(property, ":") {
			pieces := strings.Split(property, ":")
			this.AuthMechanismProperties = append(this.AuthMechanismProperties, shared.NewVariable(pieces[0], pieces[1]))
		}
	}
}

// 将认证参数转化为字符串
func (this *MongoConfig) AuthMechanismPropertiesString() string {
	s := []string{}
	for _, v := range this.AuthMechanismProperties {
		s = append(s, v.Name+":"+v.Value)
	}
	return strings.Join(s, ",")
}

// 取得Map形式的认证属性
func (this *MongoConfig) AuthMechanismPropertiesMap() map[string]string {
	m := map[string]string{}
	for _, v := range this.AuthMechanismProperties {
		m[v.Name] = v.Value
	}
	return m
}

// 设置地址
func (this *MongoConfig) SetAddr(host string, port uint) {
	if port > 0 {
		this.Addr = host + ":" + fmt.Sprintf("%d", port)
	} else {
		this.Addr = host
	}
}

// 获取Host
func (this *MongoConfig) Host() string {
	host, _, err := net.SplitHostPort(this.Addr)
	if err != nil {
		return this.Addr
	}
	return host
}

// 获取Port
func (this *MongoConfig) Port() int {
	_, port, err := net.SplitHostPort(this.Addr)
	if err != nil {
		return 0
	}
	return types.Int(port)
}
