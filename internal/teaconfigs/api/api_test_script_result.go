package api

// 单个脚本测试结果
type APITestScriptResult struct {
	Code     string   `yaml:"code" json:"code"`         // 脚本代码
	IsPassed bool     `yaml:"isPassed" json:"isPassed"` // 是否通过测试
	Failures []string `yaml:"failures" json:"failures"` // 失败
}

// 获取新对象
func NewAPITestScriptResult() *APITestScriptResult {
	return &APITestScriptResult{}
}
