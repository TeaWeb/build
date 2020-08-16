package checkpoints

// attach option
type FieldOption struct {
	Name        string
	Code        string
	Value       string // default value
	IsRequired  bool
	Size        int
	Comment     string
	Placeholder string
	RightLabel  string
	MaxLength   int
	Validate    func(value string) (ok bool, message string)
}

func NewFieldOption(name string, code string) *FieldOption {
	return &FieldOption{
		Name: name,
		Code: code,
	}
}

func (this *FieldOption) Type() string {
	return "field"
}
