package teadb

import "errors"

var ErrorDBUnavailable = errors.New("database is not available")

// 驱动基础
type BaseDriver struct {
	isAvailable bool
}

// 是否可用
func (this *BaseDriver) IsAvailable() bool {
	return this.isAvailable
}

// 设置是否可用
func (this *BaseDriver) SetIsAvailable(b bool) {
	this.isAvailable = b
}
