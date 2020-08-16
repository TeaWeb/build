package agents

import (
	"github.com/iwind/TeaGo/assert"
	"github.com/iwind/TeaGo/maps"
	"testing"
)

func TestEvalParam(t *testing.T) {
	t.Log(EvalParam("This is is message, ${ITEM.name}, ${ITEM}, ${ITE}", nil, nil, maps.Map{
		"ITEM": maps.Map{
			"name": "MySQL",
		},
	}, false))
}

// 测试空格
func TestEvalParam_Spaces(t *testing.T) {
	a := assert.NewAssertion(t)

	threshold := NewThreshold()
	threshold.Param = "[${data    .	version}]"
	threshold.Operator = ThresholdOperatorEq
	threshold.Value = "1.0.25"
	err := threshold.Validate()
	if err != nil {
		t.Error(err)
	}
	{
		result, err := threshold.Eval(map[string]interface{}{
			"data": maps.Map{
				"version": "1.0.26",
			},
		}, nil)
		a.IsNil(err)
		a.IsTrue(result == "[1.0.26]")
	}
	{
		result, err := EvalParam(threshold.Param, nil, nil, maps.Map{
			"data": maps.Map{
				"version": "1.1",
			},
		}, true)
		a.IsNil(err)
		a.IsTrue(result == "[1.1]")
	}
}

// 格式化
func TestEvalParam_Format(t *testing.T) {
	t.Log(EvalParam("this is ${0 | float('%.2f')} and ${0}", "123.456789", "", maps.Map{}, false))
}
