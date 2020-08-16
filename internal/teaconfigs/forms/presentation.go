package forms

// 界面显示信息
type Presentation struct {
	HTML       string `yaml:"html" json:"html"`
	CSS        string `yaml:"css" json:"css"`
	Javascript string `yaml:"javascript" json:"javascript"`
}

func NewPresentation() *Presentation {
	return &Presentation{}
}
