package teaconfigs

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/rands"
)

// 请求条件分组
type RequestGroup struct {
	BackendList `yaml:",inline"`

	Id              string                  `yaml:"id" json:"id"`                           // ID
	Name            string                  `yaml:"name" json:"name"`                       // 名称
	Cond            []*shared.RequestCond   `yaml:"conds" json:"conds"`                     // 匹配条件
	IPRanges        []*shared.IPRangeConfig `yaml:"ipRanges" json:"ipRanges"`               // IP范围
	IsDefault       bool                    `yaml:"isDefault" json:"isDefault"`             // 是否为默认分组
	RequestHeaders  []*shared.HeaderConfig  `yaml:"requestHeaders" json:"requestHeaders"`   // 请求Header
	ResponseHeaders []*shared.HeaderConfig  `yaml:"responseHeaders" json:"responseHeaders"` // 响应Header

	hasConds           bool
	hasIPRanges        bool
	hasRequestHeaders  bool
	hasResponseHeaders bool
}

// 获取新对象
func NewRequestGroup() *RequestGroup {
	return &RequestGroup{
		Id: rands.HexString(16),
	}
}

// 校验
func (this *RequestGroup) Validate() error {
	// cond
	this.hasConds = len(this.Cond) > 0
	for _, cond := range this.Cond {
		err := cond.Validate()
		if err != nil {
			return err
		}
	}

	// ip range
	this.hasIPRanges = len(this.IPRanges) > 0
	for _, ipRange := range this.IPRanges {
		err := ipRange.Validate()
		if err != nil {
			return err
		}
	}

	// request header
	this.hasRequestHeaders = len(this.RequestHeaders) > 0
	for _, header := range this.RequestHeaders {
		err := header.Validate()
		if err != nil {
			return err
		}
	}

	// response header
	this.hasResponseHeaders = len(this.ResponseHeaders) > 0
	for _, header := range this.ResponseHeaders {
		err := header.Validate()
		if err != nil {
			return err
		}
	}

	// backend
	this.ValidateBackends()

	return nil
}

// 添加匹配条件
func (this *RequestGroup) AddCond(cond *shared.RequestCond) {
	this.Cond = append(this.Cond, cond)
}

// 添加IP范围
func (this *RequestGroup) AddIPRange(ipRange *shared.IPRangeConfig) {
	this.IPRanges = append(this.IPRanges, ipRange)
}

// 添加请求Header
func (this *RequestGroup) AddRequestHeader(header *shared.HeaderConfig) {
	this.RequestHeaders = append(this.RequestHeaders, header)
}

// 添加响应Header
func (this *RequestGroup) AddResponseHeader(header *shared.HeaderConfig) {
	this.ResponseHeaders = append(this.ResponseHeaders, header)
}

// 判断匹配
func (this *RequestGroup) Match(formatter func(source string) string) bool {
	if this.hasConds {
		for _, cond := range this.Cond {
			if !cond.Match(formatter) {
				return false
			}
		}
	}

	if this.hasIPRanges {
		found := false
		for _, ipRange := range this.IPRanges {
			// TODO 优化 formatter 同样的参数只format一次
			if ipRange.Contains(formatter(ipRange.Param)) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// 判断是否有过滤条件
func (this *RequestGroup) HasFilters() bool {
	return (this.hasConds || this.hasIPRanges) && this.HasBackends()
}

// 判断是否有请求Header
func (this *RequestGroup) HasRequestHeaders() bool {
	return this.hasRequestHeaders
}

// 判断是否有响应Header
func (this *RequestGroup) HasResponseHeaders() bool {
	return this.hasResponseHeaders
}

// 复制
func (this *RequestGroup) Copy() *RequestGroup {
	return &RequestGroup{
		Id:              this.Id,
		Name:            this.Name,
		Cond:            this.Cond,
		IPRanges:        this.IPRanges,
		RequestHeaders:  this.RequestHeaders,
		ResponseHeaders: this.ResponseHeaders,
		IsDefault:       this.IsDefault,
	}
}
