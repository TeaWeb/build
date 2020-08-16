package widgets

// 线定义
type Line struct {
	Param    string `yaml:"param" json:"param"`
	IsFilled bool   `yaml:"isFilled" json:"isFilled"`
	Color    string `yaml:"color" json:"color"`
	Name     string `yaml:"name" json:"name"`
}

// 获取新对象
func NewLine() *Line {
	return &Line{}
}
