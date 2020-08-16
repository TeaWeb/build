package agents

import "testing"

func TestConvertSourceToMap(t *testing.T) {
	source := NewDNSSource()
	source.Type = "A"
	source.Domain = "teaos.cn"
	t.Log(ConvertSourceToMap(source))
}
