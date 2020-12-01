package teaconfigs

import (
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
	"net/http"
	"path/filepath"
	"regexp"
	"time"
)

// Fastcgi配置
type FastcgiConfig struct {
	shared.HeaderList `yaml:",inline"`

	On bool   `yaml:"on" json:"on"`
	Id string `yaml:"id" json:"id"`

	// fastcgi地址配置
	// 支持unix:/tmp/php-fpm.sock ...
	Pass string `yaml:"pass" json:"pass"`

	Index           string            `yaml:"index" json:"index"`                     //@TODO
	Params          map[string]string `yaml:"params" json:"params"`                   //@TODO
	ReadTimeout     string            `yaml:"readTimeout" json:"readTimeout"`         // @TODO 读取超时时间
	SendTimeout     string            `yaml:"sendTimeout" json:"sendTimeout"`         // @TODO 发送超时时间
	ConnTimeout     string            `yaml:"connTimeout" json:"connTimeout"`         // @TODO 连接超时时间
	PoolSize        int               `yaml:"poolSize" json:"poolSize"`               // 连接池尺寸
	PathInfoPattern string            `yaml:"pathInfoPattern" json:"pathInfoPattern"` // PATH_INFO匹配正则

	network string // 协议：tcp, unix
	address string // 地址

	paramsMap      maps.Map
	readTimeout    time.Duration
	pathInfoRegexp *regexp.Regexp
}

// 获取新对象
func NewFastcgiConfig() *FastcgiConfig {
	return &FastcgiConfig{
		On: true,
		Id: rands.HexString(16),
	}
}

// 校验配置
func (this *FastcgiConfig) Validate() error {
	this.paramsMap = maps.NewMap(this.Params)
	if !this.paramsMap.Has("SCRIPT_FILENAME") {
		this.paramsMap["SCRIPT_FILENAME"] = ""
	}
	if !this.paramsMap.Has("SERVER_SOFTWARE") {
		this.paramsMap["SERVER_SOFTWARE"] = teaconst.TeaProductCode + "/" + teaconst.TeaVersion
	}
	if !this.paramsMap.Has("REDIRECT_STATUS") {
		this.paramsMap["REDIRECT_STATUS"] = "200"
	}
	if !this.paramsMap.Has("GATEWAY_INTERFACE") {
		this.paramsMap["GATEWAY_INTERFACE"] = "CGI/1.1"
	}

	// 校验地址
	if regexp.MustCompile("^\\d+$").MatchString(this.Pass) {
		this.network = "tcp"
		this.address = "127.0.0.1:" + this.Pass
	} else if regexp.MustCompile("^(.*):(\\d+)$").MatchString(this.Pass) {
		matches := regexp.MustCompile("^(.*):(\\d+)$").FindStringSubmatch(this.Pass)
		ip := matches[1]
		port := matches[2]
		if len(ip) == 0 {
			ip = "127.0.0.1"
		}
		this.network = "tcp"
		this.address = ip + ":" + port
	} else if regexp.MustCompile("^\\d+\\.\\d+.\\d+.\\d+$").MatchString(this.Pass) {
		this.network = "tcp"
		this.address = this.Pass + ":9000"
	} else if regexp.MustCompile("^unix:(.+)$").MatchString(this.Pass) {
		matches := regexp.MustCompile("^unix:(.+)$").FindStringSubmatch(this.Pass)
		path := matches[1]
		this.network = "unix"
		this.address = path
	} else if regexp.MustCompile("^[./].+$").MatchString(this.Pass) {
		this.network = "unix"
		this.address = this.Pass
	} else {
		return errors.New("invalid 'pass' format")
	}

	// 超时时间
	if len(this.ReadTimeout) > 0 {
		duration, err := time.ParseDuration(this.ReadTimeout)
		if err != nil {
			return err
		}
		this.readTimeout = duration
	} else {
		this.readTimeout = 3 * time.Second
	}

	// 校验Header
	err := this.ValidateHeaders()
	if err != nil {
		return err
	}

	// PATH_INFO
	if len(this.PathInfoPattern) > 0 {
		reg, err := regexp.Compile(this.PathInfoPattern)
		if err != nil {
			return err
		}
		this.pathInfoRegexp = reg
	}

	return nil
}

// 过滤参数
func (this *FastcgiConfig) FilterParams(req *http.Request) maps.Map {
	params := maps.NewMap(this.paramsMap)

	// 自动添加参数
	script := params.GetString("SCRIPT_FILENAME")
	if len(script) > 0 {
		if !params.Has("SCRIPT_NAME") {
			params["SCRIPT_NAME"] = filepath.Base(script)
		}
		if !params.Has("DOCUMENT_ROOT") {
			params["DOCUMENT_ROOT"] = filepath.Dir(script)
		}
		if !params.Has("PWD") {
			params["PWD"] = filepath.Dir(script)
		}
	}

	return params
}

// 超时时间
func (this *FastcgiConfig) ReadTimeoutDuration() time.Duration {
	if this.readTimeout <= 0 {
		this.readTimeout = 30 * time.Second
	}
	return this.readTimeout
}

//网络协议
func (this *FastcgiConfig) Network() string {
	return this.network
}

// 网络地址
func (this *FastcgiConfig) Address() string {
	return this.address
}

// 读取参数
func (this *FastcgiConfig) Param(paramName string) string {
	v, _ := this.Params[paramName]
	return v
}

// PATH_INFO正则
func (this *FastcgiConfig) PathInfoRegexp() *regexp.Regexp {
	return this.pathInfoRegexp
}
