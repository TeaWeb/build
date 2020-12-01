package teaconfigs

import (
	"github.com/iwind/TeaGo/rands"
)

// 本地用户信息
type ACMELocalUser struct {
	Id    string `yaml:"id" json:"id"`
	Name  string `yaml:"name" json:"name"`
	On    bool   `yaml:"on" json:"on"`
	Email string `yaml:"email" json:"email"`
	Key   string `yaml:"key" json:"key"` // base64
	URI   string `yaml:"uri" json:"uri"`
}

// 获取新对象
func NewACMELocalUser() *ACMELocalUser {
	return &ACMELocalUser{
		Id: rands.HexString(16),
		On: true,
	}
}
