package teaconfigs

// 关闭页面配置
type ShutdownConfig struct {
	On     bool   `yaml:"on" json:"on"`
	URL    string `yaml:"url" json:"url"`
	Status int    `yaml:"status" json:"status"`
}

// 获取新对象
func NewShutdownConfig() *ShutdownConfig {
	return &ShutdownConfig{}
}

// 校验
func (this *ShutdownConfig) Validate() error {
	return nil
}
