package configs

const (
	// 内置权限
	AdminGrantAll    = "all"
	AdminGrantProxy  = "proxy"
	AdminGrantQ      = "q"
	AdminGrantApi    = "api"
	AdminGrantLog    = "log"
	AdminGrantAgent  = "agent"
	AdminGrantPlugin = "plugin"
	AdminGrantTeam   = "team"
)

// 权限定义
type AdminGrant struct {
	Name       string `yaml:"name" json:"name"`
	Code       string `yaml:"code" json:"code"`
	IsDisabled bool   `yaml:"isDisabled" json:"isDisabled"`
}

// 构造新权限
func NewAdminGrant(name string, code string) *AdminGrant {
	return &AdminGrant{
		Name: name,
		Code: code,
	}
}
