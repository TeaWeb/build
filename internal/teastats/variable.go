package teastats

// 变量说明
type Variable struct {
	Code        string `yaml:"code" json:"code"`               // 代号
	Description string `yaml:"description" json:"description"` // 描述
}

// 获取新变量
func NewVariable(code, description string) *Variable {
	return &Variable{
		Code:        code,
		Description: description,
	}
}
