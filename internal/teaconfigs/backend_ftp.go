package teaconfigs

// backend
type FTPBackendConfig struct {
	Username string `yaml:"username" json:"username"` // 用户名
	Password string `yaml:"password" json:"password"` // 密码
	Dir      string `yaml:"dir" json:"dir"`           // 目录
}
