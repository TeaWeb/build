package api

// 单个API测试结果
type APITestResult struct {
	API      string                 `yaml:"api" json:"api"`           // API
	Scripts  []*APITestScriptResult `yaml:"scripts" json:"scripts"`   // 脚本
	IsPassed bool                   `yaml:"isPassed" json:"isPassed"` // 是否通过测试
}

// 获取新对象
func NewAPITestResult() *APITestResult {
	return &APITestResult{
		IsPassed: true,
	}
}

// 添加脚本执行结果
func (this *APITestResult) AddScriptResult(scriptResult *APITestScriptResult) {
	this.Scripts = append(this.Scripts, scriptResult)

	if !scriptResult.IsPassed {
		this.IsPassed = false
	}
}
