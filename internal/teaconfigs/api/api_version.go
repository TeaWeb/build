package api

// API版本定义
type APIVersion struct {
	Name string `yaml:"name" json:"name"`
	Code string `yaml:"code" json:"code"`
	On   bool   `yaml:"on" json:"on"`
}
