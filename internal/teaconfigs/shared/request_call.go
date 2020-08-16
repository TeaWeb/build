package shared

import (
	"github.com/iwind/TeaGo/maps"
	"net/http"
)

// 请求调用
type RequestCall struct {
	Formatter         func(source string) string
	Request           *http.Request
	ResponseCallbacks []func(resp http.ResponseWriter)
	Options           maps.Map
}

// 获取新对象
func NewRequestCall() *RequestCall {
	return &RequestCall{
		Options: maps.Map{},
	}
}

// 重置
func (this *RequestCall) Reset() {
	this.Formatter = nil
	this.Request = nil
	this.ResponseCallbacks = nil
	this.Options = maps.Map{}
}

// 添加响应回调
func (this *RequestCall) AddResponseCall(callback func(resp http.ResponseWriter)) {
	this.ResponseCallbacks = append(this.ResponseCallbacks, callback)
}

// 执行响应回调
func (this *RequestCall) CallResponseCallbacks(resp http.ResponseWriter) {
	for _, callback := range this.ResponseCallbacks {
		callback(resp)
	}
}
