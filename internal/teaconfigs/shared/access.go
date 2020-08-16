package shared

// 访问控制
type AccessConfig struct {
	On      bool            `yaml:"on" json:"on"`           // 是否开启
	AllowOn bool            `yaml:"allowOn" json:"allowOn"` // 白名单是否开启
	DenyOn  bool            `yaml:"denyOn" json:"denyOn"`   // 黑名单是否开启
	Allow   []*ClientConfig `yaml:"allow" json:"allow"`     // 允许的IP
	Deny    []*ClientConfig `yaml:"deny" json:"deny"`       // 禁止的IP
}

// 校验
func (this *AccessConfig) Validate() error {
	for _, allow := range this.Allow {
		err := allow.Validate()
		if err != nil {
			return err
		}
	}

	for _, deny := range this.Deny {
		err := deny.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// 添加允许的客户端
func (this *AccessConfig) AddAllow(client *ClientConfig) {
	this.Allow = append(this.Allow, client)
}

// 添加禁用的客户端
func (this *AccessConfig) AddDeny(client *ClientConfig) {
	this.Deny = append(this.Deny, client)
}
