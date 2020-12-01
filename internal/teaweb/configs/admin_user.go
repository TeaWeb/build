package configs

import (
	"github.com/iwind/TeaGo/rands"
	"sync"
)

// 管理员用户
type AdminUser struct {
	Username string   `yaml:"username" json:"username"` // 用户名
	Password string   `yaml:"password" json:"password"` // 密码
	Role     []string `yaml:"role" json:"role"`         // 角色
	Key      string   `yaml:"key" json:"key"`           // Key，用来请求API等

	Name      string `yaml:"name" json:"name"`           // 姓名
	Avatar    string `yaml:"avatar" json:"avatar"`       // 头像
	Tel       string `yaml:"tel" json:"tel"`             // 联系电话
	CreatedAt int64  `yaml:"createdAt" json:"createdAt"` // 创建时间
	LoggedAt  int64  `yaml:"loggedAt" json:"loggedAt"`   // 最后登录时间
	LoggedIP  string `yaml:"loggedIP" json:"loggedIP"`   // 最后登录IP

	Grant []string `yaml:"grant" json:"grant"` // 权限，会细化到项目，比如：apis:example.com

	IsDisabled bool `yaml:"isDisabled" json:"isDisabled"` // 是否禁用

	countLoginTries uint // 错误登录次数
	locker          sync.Mutex
}

// 获取新对象
func NewAdminUser() *AdminUser {
	user := &AdminUser{}
	user.Key = user.GenerateKey()
	return user
}

// 判断用户是否已被授权
func (this *AdminUser) Granted(grant string) bool {
	// 角色设置
	if len(this.Role) == 0 {
		return false
	}
	for _, roleCode := range this.Role {
		role := SharedAdminConfig().FindActiveRole(roleCode)
		if role == nil {
			continue
		}
		if role.Granted(grant) {
			return true
		}
	}

	// 特殊设置
	for _, grantCode := range this.Grant {
		grant := SharedAdminConfig().FindActiveGrant(grantCode)
		if grant != nil {
			return true
		}
	}

	return false
}

func (this *AdminUser) IncreaseLoginTries() {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.countLoginTries++
}

func (this *AdminUser) CountLoginTries() uint {
	this.locker.Lock()
	defer this.locker.Unlock()
	return this.countLoginTries
}

func (this *AdminUser) ResetLoginTries() {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.countLoginTries = 0
}

// 重置状态
func (this *AdminUser) Reset() {
	this.ResetLoginTries()
}

// 生成Key
func (this *AdminUser) GenerateKey() string {
	return rands.HexString(32)
}
