package configs

// 角色定义
type AdminRole struct {
	Name       string   `yaml:"name" json:"name"`             // 名称
	Code       string   `yaml:"code" json:"code"`             // 代号
	Grant      []string `yaml:"grant" json:"grant"`           // 授权
	IsDisabled bool     `yaml:"isDisabled" json:"isDisabled"` // 是否禁用
}

// 构造新角色
func NewAdminRole() *AdminRole {
	return &AdminRole{}
}

// 判断是否授权
func (this *AdminRole) Granted(grant string) bool {
	for _, grantCode := range this.Grant {
		if grantCode == "all" {
			return true
		}

		if grantCode == grant {
			return true
		}
	}
	return false
}
