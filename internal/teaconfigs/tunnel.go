package teaconfigs

import "github.com/iwind/TeaGo/utils/string"

// 网络隧道定义
type TunnelConfig struct {
	Id       string           `yaml:"id" json:"id"`             // ID
	On       bool             `yaml:"on" json:"on"`             // 是否启用
	Endpoint string           `yaml:"endpoint" json:"endpoint"` // 终端地址
	Secret   string           `yaml:"secret" json:"secret"`     // 连接用的密钥
	TLS      bool             `yaml:"tls" json:"tls"`           // 是否支持TLS TODO 暂时没有实现
	Certs    []*SSLCertConfig `yaml:"certs" json:"certs"`       // TLS证书 TODO 暂时没有实现

	isActive bool
	errors   []string
}

// 隧道设置
func NewTunnelConfig() *TunnelConfig {
	return &TunnelConfig{
		On: true,
		Id: stringutil.Rand(16),
	}
}

// 校验
func (this *TunnelConfig) Validate() error {
	// certificates
	if len(this.Certs) > 0 {
		for _, cert := range this.Certs {
			err := cert.Validate()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// 设置错误信息
func (this *TunnelConfig) AddError(err string) {
	this.errors = append(this.errors, err)
}

// 获取错误信息
func (this *TunnelConfig) Errors() []string {
	return this.errors
}

// 设置是否已启动
func (this *TunnelConfig) SetIsActive(isActive bool) {
	this.isActive = isActive
}

// 判断是否已启动
func (this *TunnelConfig) IsActive() bool {
	return this.isActive
}
