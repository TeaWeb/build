package teautils

import (
	"testing"
)

func TestObjectToMapJSON(t *testing.T) {
	var a interface{} = nil
	m := map[string]interface{}{}
	t.Log(ObjectToMapJSON(a, &m))
	t.Log(m)
}
