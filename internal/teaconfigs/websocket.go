package teaconfigs

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/maps"
	"net/http"
	"time"
)

// websocket设置
type WebsocketConfig struct {
	// 后端服务器列表
	BackendList `yaml:",inline"`

	On bool `yaml:"on" json:"on"` // 是否开启

	// 握手超时时间
	HandshakeTimeout string `yaml:"handshakeTimeout" json:"handshakeTimeout"`

	// 允许的域名，支持 www.example.com, example.com, .example.com, *.example.com
	AllowAllOrigins bool     `yaml:"allowAllOrigins" json:"allowAllOrigins"`
	Origins         []string `yaml:"origins" json:"origins"`

	// 转发方式
	ForwardMode WebsocketForwardMode `yaml:"forwardMode" json:"forwardMode"`

	// 请求分组（从server复制而来）
	requestGroups          []*RequestGroup
	defaultRequestGroup    *RequestGroup
	hasRequestGroupFilters bool

	handshakeTimeoutDuration time.Duration
}

// 获取新对象
func NewWebsocketConfig() *WebsocketConfig {
	return &WebsocketConfig{
		On: true,
	}
}

// 校验
func (this *WebsocketConfig) Validate() error {
	// backends
	err := this.ValidateBackends()
	if err != nil {
		return err
	}

	// duration
	this.handshakeTimeoutDuration, _ = time.ParseDuration(this.HandshakeTimeout)

	// request groups
	for _, group := range this.requestGroups {
		if group.IsDefault {
			this.defaultRequestGroup = group
		}

		for _, backend := range this.Backends {
			if len(backend.RequestGroupIds) == 0 && group.Id == "default" {
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

	return nil
}

// 获取握手超时时间
func (this *WebsocketConfig) HandshakeTimeoutDuration() time.Duration {
	return this.handshakeTimeoutDuration
}

// 转发模式名称
func (this *WebsocketConfig) ForwardModeSummary() maps.Map {
	for _, mode := range AllWebsocketForwardModes() {
		if mode["mode"] == this.ForwardMode {
			return mode
		}
	}
	return nil
}

// 匹配域名
func (this *WebsocketConfig) MatchOrigin(origin string) bool {
	if this.AllowAllOrigins {
		return true
	}
	return teautils.MatchDomains(this.Origins, origin)
}

// 添加请求分组
func (this *WebsocketConfig) AddRequestGroup(group *RequestGroup) {
	this.requestGroups = append(this.requestGroups, group)
}

// 使用请求匹配分组
func (this *WebsocketConfig) MatchRequestGroup(formatter func(source string) string) *RequestGroup {
	if !this.hasRequestGroupFilters {
		return nil
	}
	for _, group := range this.requestGroups {
		if group.HasFilters() && group.Match(formatter) {
			return group
		}
	}
	return nil
}

// 取得下一个可用的后端服务
func (this *WebsocketConfig) NextBackend(call *shared.RequestCall) *BackendConfig {
	if this.hasRequestGroupFilters {
		group := this.MatchRequestGroup(call.Formatter)
		if group != nil {
			// request
			if group.HasRequestHeaders() {
				for _, h := range group.RequestHeaders {
					call.Request.Header.Set(h.Name, call.Formatter(h.Value))
				}
			}

			// response
			if group.HasResponseHeaders() {
				call.AddResponseCall(func(resp http.ResponseWriter) {
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
				call.Request.Header.Set(h.Name, call.Formatter(h.Value))
			}
		}

		// response
		if this.defaultRequestGroup.HasResponseHeaders() {
			call.AddResponseCall(func(resp http.ResponseWriter) {
				for _, h := range this.defaultRequestGroup.ResponseHeaders {
					resp.Header().Set(h.Name, call.Formatter(h.Value))
				}
			})
		}

		return this.defaultRequestGroup.NextBackend(call)
	}

	return this.BackendList.NextBackend(call)
}

// 设置调度算法
func (this *WebsocketConfig) SetupScheduling(isBackup bool) {
	for _, group := range this.requestGroups {
		group.SetupScheduling(isBackup)
	}
	this.BackendList.SetupScheduling(isBackup)
}

// 克隆运行时状态
func (this *WebsocketConfig) CloneState(oldWebsocket *WebsocketConfig) {
	if oldWebsocket == nil {
		return
	}

	// backends
	for _, backend := range this.Backends {
		oldBackend := oldWebsocket.FindBackend(backend.Id)
		if oldBackend == nil {
			continue
		}
		backend.CloneState(oldBackend)
	}
}

// 卸载方法
func (this *WebsocketConfig) OnDetach() {
	for _, backend := range this.Backends {
		backend.OnDetach()
	}
}
