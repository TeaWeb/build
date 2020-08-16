package forms

var allElementTypes = map[string]ElementInterface{
	"TextField":       new(TextField),
	"TextBox":         new(TextBox),
	"Options":         new(Options),
	"ScriptBox":       new(ScriptBox),
	"CheckBox":        new(CheckBox),
	"EnvBox":          new(EnvBox),
	"HTTPBox":         new(HTTPBox),
	"SingleValueList": new(SingleValueList),
}
