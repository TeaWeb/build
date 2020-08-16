package teaconfigs

// 特殊页面配置
type PageConfig struct {
	On        bool     `yaml:"on" json:"on"`               // TODO
	Status    []string `yaml:"status" json:"status"`       // 响应支持40x, 50x, 3x2
	URL       string   `yaml:"url" json:"url"`             // URL
	NewStatus int      `yaml:"newStatus" json:"newStatus"` // 新状态码

	statusList    []*WildcardStatus
	hasStatusList bool
}

// 获取新对象
func NewPageConfig() *PageConfig {
	return &PageConfig{
		On: true,
	}
}

// 校验
func (this *PageConfig) Validate() error {
	this.statusList = []*WildcardStatus{}
	for _, s := range this.Status {
		this.statusList = append(this.statusList, NewWildcardStatus(s))
	}
	this.hasStatusList = len(this.statusList) > 0
	return nil
}

// 检查是否匹配
func (this *PageConfig) Match(status int) bool {
	if !this.hasStatusList {
		return false
	}
	for _, s := range this.statusList {
		if s.Match(status) {
			return true
		}
	}
	return false
}
