package teautils

import (
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	json "github.com/json-iterator/go"
	"github.com/robertkrimen/otto"
	"math"
	"reflect"
	"testing"
)

func TestConvertJSONObjectSafely_NaN1(t *testing.T) {
	v := ConvertJSONObjectSafely(math.NaN())
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))
}

func TestConvertJSONObjectSafely_NaN2(t *testing.T) {
	v, err := otto.NaNValue().Export()
	if err != nil {
		t.Fatal(err)
	}
	logs.Println(v, reflect.TypeOf(v).String())
	v = ConvertJSONObjectSafely(v)
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))
}

func TestConvertJSONObjectSafely_NaN3(t *testing.T) {
	_, jsValue, err := otto.Run(`(function () { 
	return { "nan": NaN, "v": "lu", "nan2": parseInt("abc") };
})()`)
	if err != nil {
		t.Fatal(err)
	}
	v, err := jsValue.Export()
	if err != nil {
		t.Fatal(err)
	}
	logs.Println(v, reflect.TypeOf(v).String())
	v = ConvertJSONObjectSafely(v)
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))
}

func TestConvertJSONObjectSafely_Interface(t *testing.T) {
	{
		m := map[interface{}]interface{}{
			"name": "lu",
			"nested": map[interface{}]interface{}{
				"a": "b",
				"c": maps.Map{
					"d": "e",
				},
			},
		}
		m1 := ConvertJSONObjectSafely(m)
		data, err := json.Marshal(m1)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(data))
	}
}
