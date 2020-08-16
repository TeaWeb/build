package configs

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"net"
)

// 安全设置定义
type AdminSecurity struct {
	TeaVersion string `yaml:"teaVersion" json:"teaVersion"`

	Allow               []string `yaml:"allow" json:"allow"`                             //支持的IP
	Deny                []string `yaml:"deny" json:"deny"`                               // 拒绝的IP
	Secret              string   `yaml:"secret" json:"secret"`                           // 密钥
	IsDisabled          bool     `yaml:"isDisabled" json:"isDisabled"`                   // 是否禁用
	DirAutoComplete     bool     `yaml:"dirAutoComplete" json:"dirAutoComplete"`         // 是否支持目录自动补全
	LoginURL            string   `yaml:"loginURL" json:"loginURL"`                       // 登录页面的URL
	PasswordEncryptType string   `yaml:"passwordEncryptType" json:"passwordEncryptType"` // 密码加密方式

	allowIPRanges []*shared.IPRangeConfig
	denyIPRanges  []*shared.IPRangeConfig
}

// 获取新对象
func NewAdminSecurity() *AdminSecurity {
	return &AdminSecurity{
		Allow: []string{},
		Deny:  []string{},
	}
}

// 密码加密方式列表
func PasswordEncryptTypes() []maps.Map {
	return []maps.Map{
		{
			"name": "明文",
			"code": "clear",
		},
		{
			"name": "MD5",
			"code": "md5",
		},
	}
}

// 校验
func (this *AdminSecurity) Validate() error {
	// 版本兼容性
	if len(this.TeaVersion) == 0 {
		this.TeaVersion = teaconst.TeaVersion
		this.DirAutoComplete = true
		this.PasswordEncryptType = "clear"
	}

	this.allowIPRanges = []*shared.IPRangeConfig{}
	for _, s := range this.Allow {
		r, err := shared.ParseIPRange(s)
		if err != nil {
			logs.Error(err)
		} else {
			this.allowIPRanges = append(this.allowIPRanges, r)
		}
	}

	this.denyIPRanges = []*shared.IPRangeConfig{}
	for _, s := range this.Deny {
		r, err := shared.ParseIPRange(s)
		if err != nil {
			logs.Error(err)
		} else {
			this.denyIPRanges = append(this.denyIPRanges, r)
		}
	}

	return nil
}

// 判断某个IP是否允许访问
func (this *AdminSecurity) AllowIP(ip string) bool {
	netIP := net.ParseIP(ip)
	if netIP == nil {
		return true
	}

	// deny
	if len(this.denyIPRanges) > 0 {
		for _, r := range this.denyIPRanges {
			if r.Contains(ip) {
				return false
			}
		}
	}

	// allow
	if len(this.Allow) > 0 {
		for _, r := range this.allowIPRanges {
			if r.Contains(ip) {
				return true
			}
		}
		return false
	}

	return true
}

// 获取登录URL
func (this *AdminSecurity) NewLoginURL() string {
	url := "/login"
	if len(this.LoginURL) > 0 {
		url = this.LoginURL
	}
	if url[0] != '/' {
		url = "/" + url
	}
	return url
}

// 获取加密方式文字说明
func (this *AdminSecurity) PasswordEncryptTypeText() string {
	for _, m := range PasswordEncryptTypes() {
		if m.GetString("code") == this.PasswordEncryptType {
			return m.GetString("name")
		}
	}
	return "明文"
}
