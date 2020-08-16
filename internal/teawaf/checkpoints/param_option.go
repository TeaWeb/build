package checkpoints

type KeyValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ParamOptions struct {
	Options []*KeyValue `json:"options"`
}

func NewParamOptions() *ParamOptions {
	return &ParamOptions{}
}

func (this *ParamOptions) AddParam(name string, value string) {
	this.Options = append(this.Options, &KeyValue{
		Name:  name,
		Value: value,
	})
}
