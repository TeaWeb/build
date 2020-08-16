package checkpoints

type Checkpoint struct {
}

func (this *Checkpoint) Init() {

}

func (this *Checkpoint) IsRequest() bool {
	return true
}

func (this *Checkpoint) ParamOptions() *ParamOptions {
	return nil
}

func (this *Checkpoint) Options() []OptionInterface {
	return nil
}

func (this *Checkpoint) Start() {

}

func (this *Checkpoint) Stop() {

}
