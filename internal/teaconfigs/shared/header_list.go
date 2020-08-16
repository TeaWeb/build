package shared

import (
	"github.com/iwind/TeaGo/lists"
	"strings"
)

// HeaderList相关操作接口
type HeaderListInterface interface {
	// 校验
	ValidateHeaders() error

	// 取得所有的IgnoreHeader
	AllIgnoreResponseHeaders() []string

	// 添加IgnoreHeader
	AddIgnoreResponseHeader(name string)

	// 判断是否包含IgnoreHeader
	ContainsIgnoreResponseHeader(name string) bool

	// 移除IgnoreHeader
	RemoveIgnoreResponseHeader(name string)

	// 修改IgnoreHeader
	UpdateIgnoreResponseHeader(oldName string, newName string)

	// 取得所有的Header
	AllResponseHeaders() []*HeaderConfig

	// 添加Header
	AddResponseHeader(header *HeaderConfig)

	// 判断是否包含Header
	ContainsResponseHeader(name string) bool

	// 查找Header
	FindResponseHeader(headerId string) *HeaderConfig

	// 移除Header
	RemoveResponseHeader(headerId string)

	// 取得所有的请求Header
	AllRequestHeaders() []*HeaderConfig

	// 添加请求Header
	AddRequestHeader(header *HeaderConfig)

	// 查找请求Header
	FindRequestHeader(headerId string) *HeaderConfig

	// 移除请求Header
	RemoveRequestHeader(headerId string)
}

// HeaderList定义
type HeaderList struct {
	// 添加的响应Headers
	Headers []*HeaderConfig `yaml:"headers" json:"headers"`

	// 忽略的响应Headers
	IgnoreHeaders []string `yaml:"ignoreHeaders" json:"ignoreHeaders"`

	// 自定义请求Headers
	RequestHeaders []*HeaderConfig `yaml:"requestHeaders" json:"requestHeaders"`

	hasResponseHeaders bool
	hasRequestHeaders  bool

	hasIgnoreHeaders       bool
	uppercaseIgnoreHeaders []string
}

// 校验
func (this *HeaderList) ValidateHeaders() error {
	this.hasResponseHeaders = len(this.Headers) > 0
	this.hasRequestHeaders = len(this.RequestHeaders) > 0

	for _, h := range this.Headers {
		err := h.Validate()
		if err != nil {
			return err
		}
	}

	for _, h := range this.RequestHeaders {
		err := h.Validate()
		if err != nil {
			return err
		}
	}

	this.hasIgnoreHeaders = len(this.IgnoreHeaders) > 0
	this.uppercaseIgnoreHeaders = []string{}
	for _, headerKey := range this.IgnoreHeaders {
		this.uppercaseIgnoreHeaders = append(this.uppercaseIgnoreHeaders, strings.ToUpper(headerKey))
	}

	return nil
}

// 是否有Headers
func (this *HeaderList) HasResponseHeaders() bool {
	return this.hasResponseHeaders
}

// 取得所有的IgnoreHeader
func (this *HeaderList) AllIgnoreResponseHeaders() []string {
	if this.IgnoreHeaders == nil {
		return []string{}
	}
	return this.IgnoreHeaders
}

// 添加IgnoreHeader
func (this *HeaderList) AddIgnoreResponseHeader(name string) {
	if !lists.ContainsString(this.IgnoreHeaders, name) {
		this.IgnoreHeaders = append(this.IgnoreHeaders, name)
	}
}

// 判断是否包含IgnoreHeader
func (this *HeaderList) ContainsIgnoreResponseHeader(name string) bool {
	if len(this.IgnoreHeaders) == 0 {
		return false
	}
	return lists.ContainsString(this.IgnoreHeaders, name)
}

// 修改IgnoreHeader
func (this *HeaderList) UpdateIgnoreResponseHeader(oldName string, newName string) {
	result := []string{}
	for _, h := range this.IgnoreHeaders {
		if h == oldName {
			result = append(result, newName)
		} else {
			result = append(result, h)
		}
	}
	this.IgnoreHeaders = result
}

// 移除IgnoreHeader
func (this *HeaderList) RemoveIgnoreResponseHeader(name string) {
	result := []string{}
	for _, n := range this.IgnoreHeaders {
		if n == name {
			continue
		}
		result = append(result, n)
	}
	this.IgnoreHeaders = result
}

// 取得所有的Header
func (this *HeaderList) AllResponseHeaders() []*HeaderConfig {
	if this.Headers == nil {
		return []*HeaderConfig{}
	}
	return this.Headers
}

// 添加Header
func (this *HeaderList) AddResponseHeader(header *HeaderConfig) {
	this.Headers = append(this.Headers, header)
}

// 判断是否包含Header
func (this *HeaderList) ContainsResponseHeader(name string) bool {
	for _, h := range this.Headers {
		if h.Name == name {
			return true
		}
	}
	return false
}

// 查找Header
func (this *HeaderList) FindResponseHeader(headerId string) *HeaderConfig {
	for _, h := range this.Headers {
		if h.Id == headerId {
			return h
		}
	}
	return nil
}

// 移除Header
func (this *HeaderList) RemoveResponseHeader(headerId string) {
	result := []*HeaderConfig{}
	for _, h := range this.Headers {
		if h.Id == headerId {
			continue
		}
		result = append(result, h)
	}
	this.Headers = result
}

// 添加请求Header
func (this *HeaderList) AddRequestHeader(header *HeaderConfig) {
	this.RequestHeaders = append(this.RequestHeaders, header)
}

// 判断是否有请求Header
func (this *HeaderList) HasRequestHeaders() bool {
	return this.hasRequestHeaders
}

// 取得所有的请求Header
func (this *HeaderList) AllRequestHeaders() []*HeaderConfig {
	return this.RequestHeaders
}

// 查找请求Header
func (this *HeaderList) FindRequestHeader(headerId string) *HeaderConfig {
	for _, h := range this.RequestHeaders {
		if h.Id == headerId {
			return h
		}
	}
	return nil
}

// 移除请求Header
func (this *HeaderList) RemoveRequestHeader(headerId string) {
	result := []*HeaderConfig{}
	for _, h := range this.RequestHeaders {
		if h.Id == headerId {
			continue
		}
		result = append(result, h)
	}
	this.RequestHeaders = result
}

// 判断是否有Ignore Headers
func (this *HeaderList) HasIgnoreHeaders() bool {
	return this.hasIgnoreHeaders
}

// 查找大写的Ignore Headers
func (this *HeaderList) UppercaseIgnoreHeaders() []string {
	return this.uppercaseIgnoreHeaders
}
