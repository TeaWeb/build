package checkpoints

import "github.com/iwind/TeaGo/maps"

type OptionsOption struct {
	Name       string
	Code       string
	Value      string // default value
	IsRequired bool
	Size       int
	Comment    string
	RightLabel string
	Validate   func(value string) (ok bool, message string)
	Options    []maps.Map
}

func NewOptionsOption(name string, code string) *OptionsOption {
	return &OptionsOption{
		Name: name,
		Code: code,
	}
}

func (this *OptionsOption) Type() string {
	return "options"
}

func (this *OptionsOption) SetOptions(options []maps.Map) {
	this.Options = options
}
