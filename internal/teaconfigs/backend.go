package teaconfigs

import (
	"crypto/tls"
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/rands"
	"github.com/iwind/TeaGo/timers"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

// 服务后端配置
type BackendConfig struct {
	nextBackend *BackendConfig // 等待切换的下一个Backend

	shared.HeaderList `yaml:",inline"`

	TeaVersion string `yaml:"teaVersion" json:"teaVersion"`

	On           bool   `yaml:"on" json:"on"`                     // 是否启用
	Id           string `yaml:"id" json:"id"`                     // ID
	Version      int    `yaml:"version" json:"version"`           // 版本号
	Code         string `yaml:"code" json:"code"`                 // 代号
	Address      string `yaml:"address" json:"address"`           // 地址
	Scheme       string `yaml:"scheme" json:"scheme"`             // 协议，http、https、tcp、tcp+tls、ftp
	Weight       uint   `yaml:"weight" json:"weight"`             // 权重
	IsBackup     bool   `yaml:"backup" json:"isBackup"`           // 是否为备份
	FailTimeout  string `yaml:"failTimeout" json:"failTimeout"`   // 连接失败超时
	ReadTimeout  string `yaml:"readTimeout" json:"readTimeout"`   // 读取超时时间
	IdleTimeout  string `yaml:"idleTimeout" json:"idleTimeout"`   // 空闲连接超时时间
	MaxFails     int32  `yaml:"maxFails" json:"maxFails"`         // 最多失败次数
	CurrentFails int32  `yaml:"currentFails" json:"currentFails"` // 当前已失败次数
	MaxConns     int32  `yaml:"maxConns" json:"maxConns"`         // 最大并发连接数
	CurrentConns int32  `yaml:"currentConns" json:"currentConns"` // 当前连接数
	IdleConns    int32  `yaml:"idleConns" json:"idleConns"`       // 最大空闲连接数

	IsDown   bool      `yaml:"down" json:"isDown"`                           // 是否下线
	DownTime time.Time `yaml:"downTime,omitempty" json:"downTime,omitempty"` // 下线时间

	RequestGroupIds []string               `yaml:"requestGroupIds" json:"requestGroupIds"` // 所属请求分组
	RequestURI      string                 `yaml:"requestURI" json:"requestURI"`           // 转发后的请求URI
	ResponseHeaders []*shared.HeaderConfig `yaml:"responseHeaders" json:"responseHeaders"` // 响应Header
	Host            string                 `yaml:"host" json:"host"`                       // 自定义主机名

	// 健康检查URL，目前支持：
	// - http|https 返回2xx-3xx认为成功
	CheckOn       bool   `yaml:"checkOn" json:"checkOn"` // 是否开启
	CheckURL      string `yaml:"checkURL" json:"checkURL"`
	CheckInterval int    `yaml:"checkInterval" json:"checkInterval"`
	CheckTimeout  string `yaml:"checkTimeout" json:"checkTimeout"` // 超时时间

	Cert *SSLCertConfig `yaml:"cert" json:"cert"` // 请求源服务器用的证书

	// ftp
	FTP *FTPBackendConfig `yaml:"ftp" json:"ftp"`

	failTimeoutDuration time.Duration
	readTimeoutDuration time.Duration
	idleTimeoutDuration time.Duration

	hasRequestURI bool
	requestPath   string
	requestArgs   string

	hasCheckURL          bool
	checkLooper          *timers.Looper
	checkTimeoutDuration time.Duration

	upCallbacks   []func(backend *BackendConfig)
	downCallbacks []func(backend *BackendConfig)

	hasRequestHeaders  bool
	hasResponseHeaders bool

	hasHost bool

	uniqueKey string

	hasAddrVariables bool // 地址中是否含有变量
}

// 获取新对象
func NewBackendConfig() *BackendConfig {
	return &BackendConfig{
		On:         true,
		Id:         rands.HexString(16),
		TeaVersion: teaconst.TeaVersion,
	}
}

// 校验
func (this *BackendConfig) Validate() error {
	this.Compatible()

	// 证书
	if this.Cert != nil {
		err := this.Cert.Validate()
		if err != nil {
			return err
		}
	}

	// unique key
	this.uniqueKey = this.Id + "@" + fmt.Sprintf("%d", this.Version)

	// failTimeout
	if len(this.FailTimeout) > 0 {
		this.failTimeoutDuration, _ = time.ParseDuration(this.FailTimeout)
	}

	// readTimeout
	if len(this.ReadTimeout) > 0 {
		this.readTimeoutDuration, _ = time.ParseDuration(this.ReadTimeout)
	}

	// idleTimeout
	if len(this.IdleTimeout) > 0 {
		this.idleTimeoutDuration, _ = time.ParseDuration(this.IdleTimeout)
	}

	// 是否有端口
	_, _, err := net.SplitHostPort(this.Address)
	if err != nil {
		if this.Scheme == "https" {
			this.Address += ":443"
		} else {
			this.Address += ":80"
		}
	}

	// Headers
	err = this.ValidateHeaders()
	if err != nil {
		return err
	}

	// request uri
	if len(this.RequestURI) == 0 || this.RequestURI == "${requestURI}" {
		this.hasRequestURI = false
	} else {
		this.hasRequestURI = true

		if strings.Contains(this.RequestURI, "?") {
			pieces := strings.SplitN(this.RequestURI, "?", -1)
			this.requestPath = pieces[0]
			this.requestArgs = pieces[1]
		} else {
			this.requestPath = this.RequestURI
		}
	}

	// check
	this.hasCheckURL = this.CheckOn && len(this.CheckURL) > 0
	if len(this.CheckTimeout) > 0 {
		this.checkTimeoutDuration, _ = time.ParseDuration(this.CheckTimeout)
	}

	// headers
	this.hasRequestHeaders = len(this.RequestHeaders) > 0
	this.hasResponseHeaders = len(this.ResponseHeaders) > 0

	// host
	this.hasHost = len(this.Host) > 0

	// variables
	this.hasAddrVariables = shared.RegexpNamedVariable.MatchString(this.Address)

	return nil
}

// 处理兼容性
func (this *BackendConfig) Compatible() {
	if len(this.TeaVersion) == 0 {
		this.CheckOn = len(this.CheckURL) > 0
	}
}

// 连接超时时间
func (this *BackendConfig) FailTimeoutDuration() time.Duration {
	return this.failTimeoutDuration
}

// 读取超时时间
func (this *BackendConfig) ReadTimeoutDuration() time.Duration {
	return this.readTimeoutDuration
}

// 保持空闲连接时间
func (this *BackendConfig) IdleTimeoutDuration() time.Duration {
	return this.idleTimeoutDuration
}

// 候选对象代号
func (this *BackendConfig) CandidateCodes() []string {
	codes := []string{this.Id}
	if len(this.Code) > 0 {
		codes = append(codes, this.Code)
	}
	return codes
}

// 候选对象权重
func (this *BackendConfig) CandidateWeight() uint {
	return this.Weight
}

// 增加错误次数
func (this *BackendConfig) IncreaseFails() int32 {
	atomic.AddInt32(&this.CurrentFails, 1)
	return this.CurrentFails
}

// 增加连接数，并返回增加之后的数字
func (this *BackendConfig) IncreaseConn() int32 {
	return atomic.AddInt32(&this.CurrentConns, 1)
}

// 减少连接数，并返回减少之后的数字
func (this *BackendConfig) DecreaseConn() int32 {
	if this.nextBackend != nil {
		return this.nextBackend.DecreaseConn()
	}
	if this.CurrentConns == 0 {
		return 0
	}
	return atomic.AddInt32(&this.CurrentConns, -1)
}

// 添加请求分组
func (this *BackendConfig) AddRequestGroupId(requestGroupId string) {
	this.RequestGroupIds = append(this.RequestGroupIds, requestGroupId)
}

// 删除某个请求分组
func (this *BackendConfig) RemoveRequestGroupId(requestGroupId string) {
	result := []string{}
	for _, groupId := range this.RequestGroupIds {
		if groupId == requestGroupId {
			continue
		}
		result = append(result, groupId)
	}
	this.RequestGroupIds = result
}

// 判断是否有某个情趣分组ID
func (this *BackendConfig) HasRequestGroupId(requestGroupId string) bool {
	if requestGroupId == "default" && len(this.RequestGroupIds) == 0 {
		return true
	}
	return lists.ContainsString(this.RequestGroupIds, requestGroupId)
}

// 判断是否设置RequestURI
func (this *BackendConfig) HasRequestURI() bool {
	return this.hasRequestURI
}

// 获取转发后的Path
func (this *BackendConfig) RequestPath() string {
	return this.requestPath
}

// 获取转发后的附加参数
func (this *BackendConfig) RequestArgs() string {
	return this.requestArgs
}

// 健康检查
func (this *BackendConfig) CheckHealth() bool {
	if !this.CheckOn {
		return true
	}

	timeout := 5 * time.Second
	if this.checkTimeoutDuration > 0 {
		timeout = this.checkTimeoutDuration
	} else if this.failTimeoutDuration > 0 {
		timeout = this.failTimeoutDuration
	}

	// http, https
	if this.IsHTTP() {
		if len(this.CheckURL) == 0 {
			return true
		}
		req, err := http.NewRequest(http.MethodGet, this.CheckURL, nil)
		if err != nil {
			logs.Error(err)
			return false
		}
		req.Header.Set("User-Agent", teaconst.TeaProductCode+"/"+teaconst.TeaVersion)
		client := teautils.SharedHttpClient(timeout)
		resp, err := client.Do(req)
		if err != nil {
			return false
		}
		defer func() {
			err = resp.Body.Close()
			if err != nil {
				logs.Error(err)
			}
		}()
		return resp.StatusCode >= 200 && resp.StatusCode < 400
	}

	// tcp, tcp+tls
	if this.IsTCP() {
		if this.Scheme == "tcp" {
			conn, err := net.DialTimeout("tcp", this.Address, timeout)
			if err != nil {
				return false
			}
			err = conn.Close()
			if err != nil {
				logs.Error(err)
			}
			return true
		} else if this.Scheme == "tcp+tls" {
			conn, err := tls.DialWithDialer(&net.Dialer{
				Timeout: timeout,
			}, "tcp", this.Address, &tls.Config{
				InsecureSkipVerify: true,
			})
			if err != nil {
				return false
			}
			err = conn.Close()
			if err != nil {
				logs.Error(err)
			}
		}
	}

	return true
}

// 重启检查
func (this *BackendConfig) RestartChecking() {
	if this.checkLooper != nil {
		this.checkLooper.Stop()
		this.checkLooper = nil
	}

	if !this.CheckOn {
		return
	}

	if this.IsHTTP() && len(this.CheckURL) == 0 {
		return
	}

	interval := this.CheckInterval
	if interval <= 0 {
		interval = 30
	}

	this.checkLooper = timers.Loop(time.Duration(interval)*time.Second, func(looper *timers.Looper) {
		if this.CheckHealth() {
			if this.IsDown {
				this.CurrentFails = 0
				this.IsDown = false

				this.OnUp()
			}
		} else {
			this.CurrentFails++
			if this.MaxFails > 0 && this.CurrentFails >= this.MaxFails && !this.IsDown {
				this.IsDown = true
				this.DownTime = time.Now()

				this.OnDown()
			}
		}
	})
}

// 停止Checking
func (this *BackendConfig) StopChecking() {
	if this.checkLooper != nil {
		this.checkLooper.Stop()
		this.checkLooper = nil
	}
}

// 判断是否有URL Check
func (this *BackendConfig) HasCheckURL() bool {
	return this.hasCheckURL
}

// 装载事件
func (this *BackendConfig) OnAttach() {
	this.downCallbacks = []func(backend *BackendConfig){}
	this.RestartChecking()
}

// 卸载事件
func (this *BackendConfig) OnDetach() {
	this.StopChecking()
}

// 下线事件
func (this *BackendConfig) OnDown() {
	for _, callback := range this.downCallbacks {
		callback(this)
	}
}

// 上线事件
func (this *BackendConfig) OnUp() {
	for _, callback := range this.upCallbacks {
		callback(this)
	}
}

// 增加下线回调
func (this *BackendConfig) DownCallback(callback func(backend *BackendConfig)) {
	this.downCallbacks = append(this.downCallbacks, callback)
}

// 增加上线回调
func (this *BackendConfig) UpCallback(callback func(backend *BackendConfig)) {
	this.upCallbacks = append(this.upCallbacks, callback)
}

// 添加响应Header
func (this *BackendConfig) AddResponseHeader(header *shared.HeaderConfig) {
	this.ResponseHeaders = append(this.ResponseHeaders, header)
}

// 判断是否有响应Header
func (this *BackendConfig) HasResponseHeaders() bool {
	return this.hasResponseHeaders
}

// 判断是否有自定义主机名
func (this *BackendConfig) HasHost() bool {
	return this.hasHost
}

// 克隆状态
func (this *BackendConfig) CloneState(oldBackend *BackendConfig) {
	if oldBackend == nil {
		return
	}
	oldBackend.nextBackend = this
	this.IsDown = oldBackend.IsDown
	this.DownTime = oldBackend.DownTime
	this.CurrentFails = oldBackend.CurrentFails
	atomic.StoreInt32(&this.CurrentConns, oldBackend.CurrentConns)
}

// 获取唯一ID
func (this *BackendConfig) UniqueKey() string {
	return this.uniqueKey
}

// 更新
func (this *BackendConfig) Touch() {
	this.Version++
}

//是否为HTTP
func (this *BackendConfig) IsHTTP() bool {
	return len(this.Scheme) == 0 || this.Scheme == "http" || this.Scheme == "https"
}

// 是否为TCP
func (this *BackendConfig) IsTCP() bool {
	return this.Scheme == "tcp" || this.Scheme == "tcp+tls"
}

// 是否为FTP
func (this *BackendConfig) IsFTP() bool {
	return this.Scheme == "ftp"
}

// 地址中是否含有变量
func (this *BackendConfig) HasAddrVariables() bool {
	return this.hasAddrVariables
}
