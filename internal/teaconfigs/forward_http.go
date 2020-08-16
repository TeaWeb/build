package teaconfigs

// 正向代理设置
type ForwardHTTPConfig struct {
	EnableMITM bool `yaml:"enableMITM" json:"enableMITM"`
}

// 获取新对象
func NewForwardHTTPConfig() *ForwardHTTPConfig {
	return &ForwardHTTPConfig{}
}
