package shared

// 变量
type Variable struct {
	Name  string `yaml:"name" json:"name"`   // 变量名
	Value string `yaml:"value" json:"value"` // 变量值
}

// 创建新变量
func NewVariable(name string, value string) *Variable {
	return &Variable{
		Name:  name,
		Value: value,
	}
}
