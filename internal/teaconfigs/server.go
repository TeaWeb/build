package teaconfigs

import (
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs/api"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teaevents"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/build/internal/teawaf"
	"github.com/iwind/TeaGo/rands"
	"gopkg.in/yaml.v3"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/string"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// 服务配置
type ServerConfig struct {
	shared.HeaderList `yaml:",inline"`
	FastcgiList       `yaml:",inline"`
	RewriteList       `yaml:",inline"`
	BackendList       `yaml:",inline"`

	On bool `yaml:"on" json:"on"`

	Id              string   `yaml:"id" json:"id"`                           // ID
	TeaVersion      string   `yaml:"teaVersion" json:"teaVersion"`           // Tea版本
	Description     string   `yaml:"description" json:"description"`         // 描述
	Name            []string `yaml:"name" json:"name"`                       // 域名
	Http            bool     `yaml:"http" json:"http"`                       // 是否支持HTTP
	RedirectToHttps bool     `yaml:"redirectToHttps" json:"redirectToHttps"` // 是否自动跳转到Https
	IsDefault       bool     `yaml:"isDefault" json:"isDefault"`             // 是否默认的服务，找不到匹配域名时有限使用此配置

	// 监听地址
	Listen []string `yaml:"listen" json:"listen"`

	Root        string            `yaml:"root" json:"root"`               // 资源根目录
	Index       []string          `yaml:"index" json:"index"`             // 默认文件
	Charset     string            `yaml:"charset" json:"charset"`         // 字符集
	Locations   []*LocationConfig `yaml:"locations" json:"locations"`     // 地址配置
	MaxBodySize string            `yaml:"maxBodySize" json:"maxBodySize"` // 请求body最大尺寸

	GzipLevel1     uint8       `yaml:"gzipLevel" json:"gzipLevel"`         // deprecated in v0.1.8: Gzip压缩级别
	GzipMinLength1 string      `yaml:"gzipMinLength" json:"gzipMinLength"` // deprecated in v0.1.8: 需要压缩的最小内容尺寸
	Gzip           *GzipConfig `yaml:"gzip" json:"gzip"`

	// 访问日志
	AccessLog []*AccessLogConfig `yaml:"accessLog" json:"accessLog"` // 访问日志配置

	DisableAccessLog1 bool  `yaml:"disableAccessLog" json:"disableAccessLog"` // deprecated: 是否禁用访问日志
	AccessLogFields1  []int `yaml:"accessLogFields" json:"accessLogFields"`   // deprecated: 访问日志保留的字段，如果为nil，则表示没有设置

	// 统计
	DisableStat bool `yaml:"disableStat" json:"disableStat"` // 是否禁用统计

	// SSL
	SSL *SSLConfig `yaml:"ssl" json:"ssl"`

	// TCP，如果有此配置的说明为TCP代理
	TCP *TCPConfig `yaml:"tcp" json:"tcp"`

	// ForwardHTTP，如果有此配置的说明为正向HTTP代理
	ForwardHTTP *ForwardHTTPConfig `yaml:"forwardHTTP" json:"forwardHTTP"`

	// 参考：http://nginx.org/en/docs/http/ngx_http_access_module.html
	Allow []string `yaml:"allow" json:"allow"` //TODO
	Deny  []string `yaml:"deny" json:"deny"`   //TODO

	Filename string `yaml:"filename" json:"filename"` // 配置文件名

	Proxy string `yaml:"proxy" json:"proxy"` //  代理配置 TODO

	CachePolicy string `yaml:"cachePolicy" json:"cachePolicy"` // 缓存策略
	CacheOn     bool   `yaml:"cacheOn" json:"cacheOn"`         // 缓存是否打开
	cachePolicy *shared.CachePolicy

	WAFOn bool        `yaml:"wafOn" json:"wafOn"` // 是否启用
	WafId string      `yaml:"wafId" json:"wafId"` // WAF ID
	waf   *teawaf.WAF // waf object

	// API相关
	API *api.APIConfig `yaml:"api" json:"api"` // API配置

	// 看板
	RealtimeBoard *Board   `yaml:"realtimeBoard" json:"realtimeBoard"` // 即时看板
	StatBoard     *Board   `yaml:"statBoard" json:"statBoard"`         // 统计看板
	StatItems     []string `yaml:"statItems" json:"statItems"`         // 统计指标

	// 是否开启静态文件加速
	CacheStatic bool `yaml:"cacheStatic" json:"cacheStatic"`

	// 请求分组
	RequestGroups          []*RequestGroup `yaml:"requestGroups" json:"requestGroups"` // 请求条件分组
	defaultRequestGroup    *RequestGroup
	hasRequestGroupFilters bool

	Pages           []*PageConfig   `yaml:"pages" json:"pages"`                   // 特殊页，更高级的需求应该通过Location来设置
	Shutdown        *ShutdownConfig `yaml:"shutdown" json:"shutdown"`             // 关闭页
	ShutdownPageOn1 bool            `yaml:"shutdownPageOn" json:"shutdownPageOn"` // deprecated: v0.1.8, 是否开启临时关闭页面
	ShutdownPage1   string          `yaml:"shutdownPage" json:"shutdownPage"`     // deprecated: v0.1.8, 临时关闭页面

	Version int `yaml:"version" json:"version"` // 版本

	// 隧道相关
	Tunnel *TunnelConfig `yaml:"tunnel" json:"tunnel"`

	// 通知设置
	NoticeSetting map[notices.NoticeLevel][]*notices.NoticeReceiver `yaml:"noticeSetting" json:"noticeSetting"`

	// 通知条目设置
	NoticeItems struct {
		BackendDown *notices.Item `yaml:"backendDown" json:"backendDown"`
		BackendUp   *notices.Item `yaml:"backendUp" json:"backendUp"`
	} `yaml:"noticeItems" json:"noticeItems"`

	maxBodySize int64
}

// 从目录中加载配置
func LoadServerConfigsFromDir(dirPath string) []*ServerConfig {
	servers := []*ServerConfig{}

	dir := files.NewFile(dirPath)
	subFiles := dir.Glob("*.proxy.conf")
	for _, configFile := range subFiles {
		reader, err := configFile.Reader()
		if err != nil {
			logs.Error(err)
			continue
		}

		// sample
		if configFile.Name() == "server.sample.www.proxy.conf" {
			_ = reader.Close()
			continue
		}

		config := &ServerConfig{}
		err = reader.ReadYAML(config)
		if err != nil {
			_ = reader.Close()
			continue
		}
		config.Filename = configFile.Name()

		// API
		if config.API == nil {
			config.API = api.NewAPIConfig()
		}

		servers = append(servers, config)
		_ = reader.Close()
	}

	return servers
}

// 取得一个新的服务配置
func NewServerConfig() *ServerConfig {
	server := &ServerConfig{
		On:      true,
		Id:      rands.HexString(16),
		API:     api.NewAPIConfig(),
		CacheOn: true,
		WAFOn:   true,
	}
	server.SetupNoticeItems()
	return server
}

// 从配置文件中读取配置信息
func NewServerConfigFromFile(filename string) (*ServerConfig, error) {
	if len(filename) == 0 {
		return nil, errors.New("filename should not be empty")
	}
	reader, err := files.NewReader(Tea.ConfigFile(filename))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = reader.Close()
	}()

	config := &ServerConfig{}
	err = reader.ReadYAML(config)
	if err != nil {
		return nil, err
	}

	config.compatible()
	config.Filename = filename

	// 初始化
	if len(config.Locations) == 0 {
		config.Locations = []*LocationConfig{}
	}
	if config.API == nil {
		config.API = api.NewAPIConfig()
	}
	if config.Headers == nil {
		config.Headers = []*shared.HeaderConfig{}
	}
	if config.IgnoreHeaders == nil {
		config.IgnoreHeaders = []string{}
	}

	return config, nil
}

// 通过ID读取配置信息
func NewServerConfigFromId(serverId string) *ServerConfig {
	if len(serverId) == 0 {
		return nil
	}

	filename := "server." + serverId + ".proxy.conf"
	file := files.NewFile(Tea.ConfigFile(filename))
	if !file.Exists() {
		// 遍历查找
		for _, server := range LoadServerConfigsFromDir(Tea.ConfigDir()) {
			if server.Id == serverId {
				server.compatible()
				return server
			}
		}

		return nil
	}
	data, err := file.ReadAll()
	if err != nil {
		logs.Error(err)
		return nil
	}

	server := &ServerConfig{}
	err = yaml.Unmarshal(data, server)
	if err != nil {
		logs.Error(err)
		return nil
	}

	server.compatible()

	return server
}

// 校验配置
func (this *ServerConfig) Validate() error {
	// 兼容设置
	this.compatible()

	// 最大Body尺寸
	maxBodySize, _ := stringutil.ParseFileSize(this.MaxBodySize)
	this.maxBodySize = int64(maxBodySize)

	// gzip
	if this.Gzip != nil {
		err := this.Gzip.Validate()
		if err != nil {
			return err
		}
	}

	// ssl
	if this.SSL != nil {
		err := this.SSL.Validate()
		if err != nil {
			return err
		}
	}

	// backends
	err := this.ValidateBackends()
	if err != nil {
		return err
	}

	// locations
	for _, location := range this.Locations {
		// 复制request group
		location.requestGroups = []*RequestGroup{}
		if len(location.Backends) > 0 {
			for _, group := range this.RequestGroups {
				location.AddRequestGroup(group.Copy())
			}
		}

		err := location.Validate()
		if err != nil {
			return err
		}
	}

	// fastcgi
	err = this.ValidateFastcgi()
	if err != nil {
		return err
	}

	// rewrite rules
	err = this.ValidateRewriteRules()
	if err != nil {
		return err
	}

	// headers
	err = this.ValidateHeaders()
	if err != nil {
		return err
	}

	// 校验缓存配置
	if len(this.CachePolicy) > 0 && this.CacheOn {
		policy := shared.NewCachePolicyFromFile(this.CachePolicy)
		if policy != nil {
			err := policy.Validate()
			if err != nil {
				return err
			}
			this.cachePolicy = policy
		}
	}

	// waf
	if len(this.WafId) > 0 && this.WAFOn {
		waf := SharedWAFList().FindWAF(this.WafId)
		if waf != nil {
			err := waf.Init()
			if err != nil {
				return err
			}
			this.waf = waf
		}
	}

	// api
	if this.API == nil {
		this.API = api.NewAPIConfig()
	}

	err = this.API.Validate()
	if err != nil {
		return err
	}

	// request groups
	for _, group := range this.RequestGroups {
		group.Backends = []*BackendConfig{}
		group.Scheduling = this.Scheduling

		if group.IsDefault {
			this.defaultRequestGroup = group
		}

		for _, backend := range this.Backends {
			if len(backend.RequestGroupIds) == 0 && group.IsDefault {
				group.AddBackend(backend)
			} else if backend.HasRequestGroupId(group.Id) {
				group.AddBackend(backend)
			}
		}

		err := group.Validate()
		if err != nil {
			return err
		}
		if group.HasFilters() {
			this.hasRequestGroupFilters = true
		}
	}

	// pages
	for _, page := range this.Pages {
		err := page.Validate()
		if err != nil {
			return err
		}
	}

	// shutdown
	if this.Shutdown != nil {
		err = this.Shutdown.Validate()
		if err != nil {
			return err
		}
	}

	// tunnel
	if this.Tunnel != nil {
		err = this.Tunnel.Validate()
		if err != nil {
			return err
		}
	}

	// access log
	if this.AccessLog != nil {
		for _, a := range this.AccessLog {
			err = a.Validate()
			if err != nil {
				return err
			}
		}
	}

	// tcp
	if this.TCP != nil {
		err = this.TCP.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// 版本相关兼容性
func (this *ServerConfig) compatible() {
	// 版本相关
	if len(this.TeaVersion) == 0 {
		// cache 默认值
		this.CacheOn = true

		// waf 默认值
		this.WAFOn = true
	}

	// v0.1.3
	if stringutil.VersionCompare(this.TeaVersion, "0.1.3") < 0 {
		// waf 默认值
		this.WAFOn = true
	}

	// v0.1.5
	if len(this.TeaVersion) == 0 || stringutil.VersionCompare(this.TeaVersion, "0.1.5") <= 0 {
		if len(this.AccessLog) == 0 {
			this.AccessLog = []*AccessLogConfig{
				{
					Id:      rands.HexString(16),
					On:      !this.DisableAccessLog1,
					Fields:  this.AccessLogFields1,
					Status1: true,
					Status2: true,
					Status3: true,
					Status4: true,
					Status5: true,
				},
			}
		}
	}

	// v0.1.8
	if stringutil.VersionCompare(this.TeaVersion, "0.1.8") < 0 {
		// gzip
		if this.GzipLevel1 > 0 {
			this.Gzip = &GzipConfig{
				Level:     int8(this.GzipLevel1),
				MinLength: this.GzipMinLength1,
			}
			this.GzipLevel1 = 0
		}

		// pages
		for _, page := range this.Pages {
			page.NewStatus = 200
		}

		// shutdown
		if this.ShutdownPageOn1 {
			shutdown := NewShutdownConfig()
			shutdown.On = true
			shutdown.URL = this.ShutdownPage1
			shutdown.Status = 200
			this.Shutdown = shutdown
		}
	}

	for _, location := range this.Locations {
		location.Compatible(this.TeaVersion)
	}
}

// 最大Body尺寸
func (this *ServerConfig) MaxBodyBytes() int64 {
	return this.maxBodySize
}

// 添加域名
func (this *ServerConfig) AddName(name ...string) {
	this.Name = append(this.Name, name...)
}

// 添加监听地址
func (this *ServerConfig) AddListen(address string) {
	this.Listen = append(this.Listen, address)
}

// 分解所有监听地址
func (this *ServerConfig) ParseListenAddresses() []string {
	result := []string{}
	var reg = regexp.MustCompile(`\[\s*(\d+)\s*[,:-]\s*(\d+)\s*]$`)
	for _, addr := range this.Listen {
		match := reg.FindStringSubmatch(addr)
		if len(match) == 0 {
			result = append(result, addr)
		} else {
			min := types.Int(match[1])
			max := types.Int(match[2])
			if min > max {
				min, max = max, min
			}
			for i := min; i <= max; i++ {
				newAddr := reg.ReplaceAllString(addr, ":"+strconv.Itoa(i))
				result = append(result, newAddr)
			}
		}
	}
	return result
}

// 获取某个位置上的配置
func (this *ServerConfig) LocationAtIndex(index int) *LocationConfig {
	if index < 0 {
		return nil
	}
	if index >= len(this.Locations) {
		return nil
	}
	location := this.Locations[index]
	err := location.Validate()
	if err != nil {
		logs.Error(err)
	}
	return location
}

// 保存
func (this *ServerConfig) Save() error {
	shared.Locker.Lock()
	defer shared.Locker.WriteUnlockNotify()

	if len(this.Filename) == 0 {
		return errors.New("'filename' should be specified")
	}

	this.TeaVersion = teaconst.TeaVersion
	this.Version++

	writer, err := files.NewWriter(Tea.ConfigFile(this.Filename))
	if err != nil {
		return err
	}
	_, err = writer.WriteYAML(this)
	_ = writer.Close()
	return err
}

// 删除
func (this *ServerConfig) Delete() error {
	if len(this.Filename) == 0 {
		return errors.New("'filename' should be specified")
	}

	// 删除key
	if this.SSL != nil {
		err := this.SSL.DeleteFiles()
		if err != nil {
			logs.Error(err)
		}
	}

	return files.NewFile(Tea.ConfigFile(this.Filename)).Delete()
}

// 判断是否和域名匹配
func (this *ServerConfig) MatchName(name string) (matchedName string, matched bool) {
	isMatched := teautils.MatchDomains(this.Name, name)
	if isMatched {
		return name, true
	}
	return
}

// 取得第一个非泛解析的域名
func (this *ServerConfig) FirstName() string {
	for _, name := range this.Name {
		if strings.Contains(name, "*") {
			continue
		}
		return name
	}
	return ""
}

// 添加路径规则
func (this *ServerConfig) AddLocation(location *LocationConfig) {
	this.Locations = append(this.Locations, location)
}

// 缓存策略
func (this *ServerConfig) CachePolicyObject() *shared.CachePolicy {
	return this.cachePolicy
}

// WAF策略
func (this *ServerConfig) WAF() *teawaf.WAF {
	return this.waf
}

// 根据Id查找Location
func (this *ServerConfig) FindLocation(locationId string) *LocationConfig {
	for _, location := range this.Locations {
		if location.Id == locationId {
			err := location.Validate()
			if err != nil {
				logs.Error(err)
			}
			return location
		}
	}
	return nil
}

// 删除Location
func (this *ServerConfig) RemoveLocation(locationId string) {
	result := []*LocationConfig{}
	for _, location := range this.Locations {
		if location.Id == locationId {
			continue
		}
		result = append(result, location)
	}
	this.Locations = result
}

// 查找HeaderList
func (this *ServerConfig) FindHeaderList(locationId string, backendId string, rewriteId string, fastcgiId string) (headerList shared.HeaderListInterface, err error) {
	if len(rewriteId) > 0 { // Rewrite
		if len(locationId) > 0 { // Server > Location > Rewrite
			location := this.FindLocation(locationId)
			if location == nil {
				err = errors.New("找不到要修改的location")
				return
			}

			rewrite := location.FindRewriteRule(rewriteId)
			if rewrite == nil {
				err = errors.New("找不到要修改的rewrite")
				return
			}
			headerList = rewrite
		} else { // Server > Rewrite
			rewrite := this.FindRewriteRule(rewriteId)
			if rewrite == nil {
				err = errors.New("找不到要修改的rewrite")
				return
			}
			headerList = rewrite
		}
	} else if len(fastcgiId) > 0 { // Fastcgi
		if len(locationId) > 0 { // Server > Location > Fastcgi
			location := this.FindLocation(locationId)
			if location == nil {
				err = errors.New("找不到要修改的location")
				return
			}

			fastcgi := location.FindFastcgi(fastcgiId)
			if fastcgi == nil {
				err = errors.New("找不到要修改的Fastcgi")
				return
			}
			headerList = fastcgi
		} else { // Server > Fastcgi
			fastcgi := this.FindFastcgi(fastcgiId)
			if fastcgi == nil {
				err = errors.New("找不到要修改的Fastcgi")
				return
			}
			headerList = fastcgi
		}
	} else if len(backendId) > 0 { // Backend
		if len(locationId) > 0 { // Server > Location > Backend
			location := this.FindLocation(locationId)
			if location == nil {
				err = errors.New("找不到要修改的location")
				return
			}

			backend := location.FindBackend(backendId)
			if backend == nil {
				err = errors.New("找不到要修改的Backend")
				return
			}
			headerList = backend
		} else { // Server > Backend
			backend := this.FindBackend(backendId)
			if backend == nil {
				err = errors.New("找不到要修改的Backend")
				return
			}
			headerList = backend
		}
	} else if len(locationId) > 0 { // Location
		location := this.FindLocation(locationId)
		if location == nil {
			err = errors.New("找不到要修改的location")
			return
		}
		headerList = location
	} else { // Server
		headerList = this
	}

	return
}

// 查找FastcgiList
func (this *ServerConfig) FindFastcgiList(locationId string) (fastcgiList FastcgiListInterface, err error) {
	if len(locationId) > 0 {
		location := this.FindLocation(locationId)
		if location == nil {
			err = errors.New("找不到要修改的location")
			return
		}
		fastcgiList = location
		return
	}
	fastcgiList = this
	return
}

// 查找重写规则
func (this *ServerConfig) FindRewriteList(locationId string) (rewriteList RewriteListInterface, err error) {
	if len(locationId) > 0 {
		location := this.FindLocation(locationId)
		if location == nil {
			err = errors.New("找不到要修改的location")
			return
		}
		rewriteList = location
		return
	}
	rewriteList = this
	return
}

// 查找后端服务器列表
func (this *ServerConfig) FindBackendList(locationId string, websocket bool) (backendList BackendListInterface, err error) {
	if len(locationId) > 0 {
		location := this.FindLocation(locationId)
		if location == nil {
			err = errors.New("找不到要修改的location")
			return
		}
		if websocket {
			if location.Websocket == nil {
				err = errors.New("websocket未设置")
				return
			}
			return location.Websocket, nil
		} else {
			return location, nil
		}
	}
	return this, nil
}

// 移动位置
func (this *ServerConfig) MoveLocation(fromIndex int, toIndex int) {
	if fromIndex < 0 || fromIndex >= len(this.Locations) {
		return
	}
	if toIndex < 0 || toIndex >= len(this.Locations) {
		return
	}
	if fromIndex == toIndex {
		return
	}

	location := this.Locations[fromIndex]
	newList := []*LocationConfig{}
	for i := 0; i < len(this.Locations); i++ {
		if i == fromIndex {
			continue
		}
		if fromIndex > toIndex && i == toIndex {
			newList = append(newList, location)
		}
		newList = append(newList, this.Locations[i])
		if fromIndex < toIndex && i == toIndex {
			newList = append(newList, location)
		}
	}

	this.Locations = newList
}

// 是否在引用某个代理
func (this *ServerConfig) RefersProxy(proxyId string) (description string, referred bool) {
	if this.Proxy == proxyId {
		return "server", true
	}
	for _, l := range this.Locations {
		if l.RefersProxy(proxyId) {
			return l.Pattern, true
		}
	}
	for _, r := range this.Rewrite {
		if r.RefersProxy(proxyId) {
			return r.Pattern, true
		}
	}
	return "", false
}

// 添加请求分组
func (this *ServerConfig) AddRequestGroup(group *RequestGroup) {
	this.RequestGroups = append(this.RequestGroups, group)
}

// 删除请求分组
func (this *ServerConfig) RemoveRequestGroup(groupId string) {
	result := []*RequestGroup{}
	for _, g := range this.RequestGroups {
		if g.Id == groupId {
			continue
		}
		result = append(result, g)
	}
	this.RequestGroups = result
}

// 查找请求分组
func (this *ServerConfig) FindRequestGroup(groupId string) *RequestGroup {
	for _, g := range this.RequestGroups {
		if g.Id == groupId {
			return g
		}
	}
	return nil
}

// 使用请求匹配分组
func (this *ServerConfig) MatchRequestGroup(formatter func(source string) string) *RequestGroup {
	if !this.hasRequestGroupFilters {
		return nil
	}
	for _, group := range this.RequestGroups {
		if group.HasFilters() && group.Match(formatter) {
			return group
		}
	}
	return nil
}

// 取得下一个可用的后端服务
func (this *ServerConfig) NextBackend(call *shared.RequestCall) *BackendConfig {
	if this.hasRequestGroupFilters {
		group := this.MatchRequestGroup(call.Formatter)
		if group != nil {
			// request
			if group.HasRequestHeaders() {
				for _, h := range group.RequestHeaders {
					if h.HasVariables() {
						call.Request.Header.Set(h.Name, call.Formatter(h.Value))
					} else {
						call.Request.Header.Set(h.Name, h.Value)
					}
				}
			}

			// response
			if group.HasResponseHeaders() {
				call.AddResponseCall(func(resp http.ResponseWriter) {
					// TODO 应用ignore headers
					for _, h := range group.ResponseHeaders {
						resp.Header().Set(h.Name, call.Formatter(h.Value))
					}
				})
			}

			return group.BackendList.NextBackend(call)
		}
	}

	// 默认分组
	if this.defaultRequestGroup != nil {
		// request
		if this.defaultRequestGroup.HasRequestHeaders() {
			for _, h := range this.defaultRequestGroup.RequestHeaders {
				if h.HasVariables() {
					call.Request.Header.Set(h.Name, call.Formatter(h.Value))
				} else {
					call.Request.Header.Set(h.Name, h.Value)
				}
			}
		}

		// response
		if this.defaultRequestGroup.HasResponseHeaders() {
			call.AddResponseCall(func(resp http.ResponseWriter) {
				for _, h := range this.defaultRequestGroup.ResponseHeaders {
					// TODO 应用ignore headers
					resp.Header().Set(h.Name, call.Formatter(h.Value))
				}
			})
		}

		return this.defaultRequestGroup.NextBackend(call)
	}

	return this.BackendList.NextBackend(call)
}

// 取得下一个可用的后端服务，并排除某个后端
func (this *ServerConfig) NextBackendIgnore(call *shared.RequestCall, backendIds []string) *BackendConfig {
	if this.hasRequestGroupFilters {
		group := this.MatchRequestGroup(call.Formatter)
		if group != nil {
			// request
			if group.HasRequestHeaders() {
				for _, h := range group.RequestHeaders {
					if h.HasVariables() {
						call.Request.Header.Set(h.Name, call.Formatter(h.Value))
					} else {
						call.Request.Header.Set(h.Name, h.Value)
					}
				}
			}

			// response
			if group.HasResponseHeaders() {
				call.AddResponseCall(func(resp http.ResponseWriter) {
					// TODO 应用ignore headers
					for _, h := range group.ResponseHeaders {
						resp.Header().Set(h.Name, call.Formatter(h.Value))
					}
				})
			}

			backendList := group.BackendList.CloneBackendList()
			backendList.DeleteBackends(backendIds)
			backendList.SetupScheduling(false)
			return backendList.NextBackend(call)
		}
	}

	// 默认分组
	if this.defaultRequestGroup != nil {
		// request
		if this.defaultRequestGroup.HasRequestHeaders() {
			for _, h := range this.defaultRequestGroup.RequestHeaders {
				if h.HasVariables() {
					call.Request.Header.Set(h.Name, call.Formatter(h.Value))
				} else {
					call.Request.Header.Set(h.Name, h.Value)
				}
			}
		}

		// response
		if this.defaultRequestGroup.HasResponseHeaders() {
			call.AddResponseCall(func(resp http.ResponseWriter) {
				for _, h := range this.defaultRequestGroup.ResponseHeaders {
					// TODO 应用ignore headers
					resp.Header().Set(h.Name, call.Formatter(h.Value))
				}
			})
		}

		backendList := this.defaultRequestGroup.CloneBackendList()
		backendList.DeleteBackends(backendIds)
		backendList.SetupScheduling(false)
		return backendList.NextBackend(call)
	}

	backendList := this.BackendList.CloneBackendList()
	backendList.DeleteBackends(backendIds)
	backendList.SetupScheduling(false)
	return backendList.NextBackend(call)
}

// 设置调度算法
func (this *ServerConfig) SetupScheduling(isBackup bool) {
	for _, group := range this.RequestGroups {
		group.SetupScheduling(isBackup)
	}
	this.BackendList.SetupScheduling(isBackup)
}

// 添加Page
func (this *ServerConfig) AddPage(page *PageConfig) {
	this.Pages = append(this.Pages, page)
}

// 装载事件
func (this *ServerConfig) OnAttach() {
	// 开启后端健康检查
	this.checkBackends(this, nil, nil)
	for _, location := range this.Locations {
		location.OnAttach()
		this.checkBackends(this, location, nil)

		if location.Websocket != nil {
			this.checkBackends(this, location, location.Websocket)
		}
	}

	// 开启WAF
	if this.waf != nil {
		this.waf.Start()
	}
}

// 卸载事件
func (this *ServerConfig) OnDetach() {
	// 停止后端健康检查
	backends := []*BackendConfig{}
	for _, backend := range this.Backends {
		if !lists.Contains(backends, backend) {
			backends = append(backends, backend)
		}
	}
	for _, location := range this.Locations {
		location.OnDetach()

		for _, backend := range location.Backends {
			if !lists.Contains(backends, backend) {
				backends = append(backends, backend)
			}
		}
	}
	for _, backend := range backends {
		backend.OnDetach()
	}

	// 停止WAF
	if this.waf != nil {
		this.waf.Stop()
		this.waf = nil
	}
}

// 判断是否为TCP
func (this *ServerConfig) IsTCP() bool {
	return this.TCP != nil
}

// 判断是否为HTTP
func (this *ServerConfig) IsHTTP() bool {
	return this.TCP == nil
}

// 是否为正向代理
func (this *ServerConfig) IsForwardHTTP() bool {
	return this.ForwardHTTP != nil
}

// 克隆运行时状态
func (this *ServerConfig) CloneState(oldServer *ServerConfig) {
	if oldServer == nil {
		return
	}

	// backends
	for _, backend := range this.Backends {
		oldBackend := oldServer.FindBackend(backend.Id)
		if oldBackend == nil {
			continue
		}
		backend.CloneState(oldBackend)
	}

	// 路径规则
	for _, location := range this.Locations {
		oldLocation := oldServer.FindLocation(location.Id)
		if oldLocation == nil {
			continue
		}
		location.CloneState(oldLocation)
	}
}

// 查找证书
func (this *ServerConfig) FindCerts(certId string) (result []*SSLCertConfig) {
	if this.SSL == nil || len(this.SSL.Certs) == 0 {
		return
	}
	for _, cert := range this.SSL.Certs {
		if cert.Id == certId {
			result = append(result, cert)
		}
	}
	return
}

// 检查是否匹配关键词
func (this *ServerConfig) MatchKeyword(keyword string) (matched bool, name string, tags []string) {
	// 描述
	if teautils.MatchKeyword(this.Description, keyword) {
		matched = true

		name = this.Description
		if len(this.Name) > 0 {
			tags = this.Name
		}
		return
	}

	// 域名
	for _, n := range this.Name {
		if teautils.MatchKeyword(n, keyword) {
			matched = true
			name = this.Description
			tags = this.Name
			break
		}
	}

	// 地址
	for _, addr := range this.Listen {
		if strings.Index(addr, keyword) > -1 {
			matched = true
			name = this.Description
			tags = []string{addr}
			return
		}
	}

	if this.SSL != nil {
		for _, addr := range this.SSL.Listen {
			if strings.Index(addr, keyword) > -1 {
				matched = true
				name = this.Description
				tags = []string{addr}
				return
			}
		}
	}

	return
}

// 添加通知接收者
func (this *ServerConfig) AddNoticeReceiver(level notices.NoticeLevel, receiver *notices.NoticeReceiver) {
	if this.NoticeSetting == nil {
		this.NoticeSetting = map[notices.NoticeLevel][]*notices.NoticeReceiver{}
	}
	receivers, found := this.NoticeSetting[level]
	if !found {
		receivers = []*notices.NoticeReceiver{}
	}
	receivers = append(receivers, receiver)
	this.NoticeSetting[level] = receivers
}

// 删除通知接收者
func (this *ServerConfig) RemoveNoticeReceiver(level notices.NoticeLevel, receiverId string) {
	if this.NoticeSetting == nil {
		return
	}
	receivers, found := this.NoticeSetting[level]
	if !found {
		return
	}

	result := []*notices.NoticeReceiver{}
	for _, r := range receivers {
		if r.Id == receiverId {
			continue
		}
		result = append(result, r)
	}
	this.NoticeSetting[level] = result
}

// 查找一个或多个级别对应的接收者，并合并相同的接收者
func (this *ServerConfig) FindAllNoticeReceivers(level ...notices.NoticeLevel) []*notices.NoticeReceiver {
	if len(level) == 0 {
		return []*notices.NoticeReceiver{}
	}

	m := maps.Map{} // mediaId_user => bool
	result := []*notices.NoticeReceiver{}
	for _, l := range level {
		receivers, ok := this.NoticeSetting[l]
		if !ok {
			continue
		}
		for _, receiver := range receivers {
			if !receiver.On {
				continue
			}
			key := receiver.Key()
			if m.Has(key) {
				continue
			}
			m[key] = true
			result = append(result, receiver)
		}
	}
	return result
}

// 设置通知
func (this *ServerConfig) SetupNoticeItems() {
	if this.NoticeItems.BackendDown == nil {
		this.NoticeItems.BackendDown = &notices.Item{
			On:      true,
			Level:   notices.NoticeLevelWarning,
			Subject: "后端服务器\"${backend.address}\"自动下线通知",
			Body:    "后端服务器\"${backend.address}\"已经因\"${cause}\"自动下线\n位置：${position}",
		}
	}

	if this.NoticeItems.BackendUp == nil {
		this.NoticeItems.BackendUp = &notices.Item{
			On:      true,
			Level:   notices.NoticeLevelSuccess,
			Subject: "后端服务器\"${backend.address}\"自动上线通知",
			Body:    "后端服务器\"${backend.address}\"已经因\"${cause}\"自动上线\n位置：${position}",
		}
	}
}

// 从请求中构建通知设置
func (this *ServerConfig) SetupNoticeItemsFromRequest(req *http.Request) {
	if req == nil {
		return
	}
	this.NoticeItems.BackendDown = notices.NewItemFromRequest(req, "backendDown")
	this.NoticeItems.BackendUp = notices.NewItemFromRequest(req, "backendUp")
}

// 检查后端服务器
func (this *ServerConfig) checkBackends(server *ServerConfig, location *LocationConfig, websocket *WebsocketConfig) {
	backends := []*BackendConfig{}
	var list BackendListInterface
	if websocket != nil {
		backends = websocket.Backends
		list = websocket
	} else if location != nil {
		backends = location.Backends
		list = location
	} else {
		backends = server.Backends
		list = server
	}

	for _, backend := range backends {
		backend.OnAttach()
		backend.DownCallback(func(backend *BackendConfig) {
			teaevents.Post(&BackendDownEvent{
				Server:    this,
				Backend:   backend,
				Location:  location,
				Websocket: websocket,
			})
			if backend.IsBackup {
				list.SetupScheduling(true)
			} else {
				list.SetupScheduling(false)
			}
		})
		backend.UpCallback(func(backend *BackendConfig) {
			teaevents.Post(&BackendUpEvent{
				Server:    this,
				Backend:   backend,
				Location:  location,
				Websocket: websocket,
			})
			if backend.IsBackup {
				list.SetupScheduling(true)
			} else {
				list.SetupScheduling(false)
			}
		})
	}
}

// 统计指标
func (this *ServerConfig) AddStatItem(item string) {
	if lists.ContainsString(this.StatItems, item) {
		return
	}
	this.StatItems = append(this.StatItems, item)
}

// 删除统计指标
func (this *ServerConfig) RemoveStatItem(item string) {
	result := []string{}
	for _, i := range this.StatItems {
		if i == item {
			continue
		}
		result = append(result, i)
	}
	this.StatItems = result
}
