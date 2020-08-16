package agents

import (
	"github.com/iwind/TeaGo/assert"
	"github.com/iwind/TeaGo/maps"
	"testing"
	"time"
)

func TestThreshold_Test(t *testing.T) {
	a := assert.NewAssertion(t)

	threshold := NewThreshold()
	threshold.Param = "${0}"
	threshold.Operator = ThresholdOperatorGt
	threshold.Value = "12"
	err := threshold.Validate()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(threshold.Test("123", nil))

	// v0.1.1之前的Bug，内容中不能含有\n
	{
		threshold = NewThreshold()
		threshold.Param = `${0}.replace("\n", "")`
		threshold.Operator = ThresholdOperatorContains
		threshold.Value = "qy-api"
		threshold.supportsMath = true
		a.IsTrue(threshold.Test(`"31399 qy-api\n5409"`, nil))

		threshold = NewThreshold()
		threshold.Param = `${0}`
		threshold.Operator = ThresholdOperatorContains
		threshold.Value = "qy-api"
		threshold.supportsMath = true
		a.IsTrue(threshold.Test(`"31399 qy-api\n5409"`, nil))

		threshold = NewThreshold()
		threshold.Param = `${0}`
		threshold.Operator = ThresholdOperatorContains
		threshold.Value = `qy-api\n5409`
		err = threshold.Validate()
		if err != nil {
			t.Fatal(err)
		}
		a.IsTrue(threshold.Test(`"31399 qy-api\n5409"`, nil))
	}

	threshold.Param = "${1}"
	threshold.Operator = ThresholdOperatorGt
	err = threshold.Validate()
	if err != nil {
		t.Fatal(err)
	}
	a.IsTrue(threshold.Test([]interface{}{1, 200, 3}, nil))

	threshold.Param = "${host}"
	threshold.Operator = ThresholdOperatorPrefix
	threshold.Value = "127."
	err = threshold.Validate()
	if err != nil {
		t.Fatal(err)
	}
	a.IsTrue(threshold.Test(map[string]interface{}{
		"host": "127.0.0.1",
	}, nil))

	threshold.Param = "${data.version}"
	threshold.Operator = ThresholdOperatorEq
	threshold.Value = "1.0.25"
	err = threshold.Validate()
	if err != nil {
		t.Fatal(err)
	}
	a.IsTrue(threshold.Test(map[string]interface{}{
		"data": maps.Map{
			"version": "1.0.25",
		},
	}, nil))

	threshold.Param = "${data.version1}"
	threshold.Operator = ThresholdOperatorNumberEq
	threshold.Value = "0"
	err = threshold.Validate()
	if err != nil {
		t.Fatal(err)
	}
	a.IsTrue(threshold.Test(map[string]interface{}{
		"data": maps.Map{
			"version": "1.25",
		},
	}, nil))

	threshold.Param = "${data.hello.world.0}"
	threshold.Operator = ThresholdOperatorEq
	threshold.Value = "1"
	a.IsTrue(threshold.Test(map[string]interface{}{
		"data": maps.Map{
			"version": "1.0.25",
			"hello": maps.Map{
				"world": []string{"1", "2", "3", "4", "5"},
			},
		},
	}, nil))
}

// 测试修改
func TestThreshold_Test2(t *testing.T) {
	a := assert.NewAssertion(t)

	threshold := NewThreshold()
	threshold.Param = "${changes}"
	threshold.Operator = ThresholdOperatorEq
	threshold.Value = "true"
	err := threshold.Validate()
	if err != nil {
		t.Fatal(err)
	}
	a.IsTrue(threshold.Test(maps.Map{
		"changes": true,
	}, nil))
}

// 测试多级获取数据
func TestThreshold_Eval(t *testing.T) {
	a := assert.NewAssertion(t)

	threshold := NewThreshold()
	threshold.Param = "${data.hello.world.0} * 100 / ${data.hello.world.1}"
	threshold.Operator = ThresholdOperatorEq
	err := threshold.Validate()
	if err != nil {
		t.Fatal(err)
	}

	{
		result, err := threshold.Eval(map[string]interface{}{
			"data": maps.Map{
				"version": "1.0.25",
				"hello": maps.Map{
					"world": []string{"1", "2", "3", "4", "5"},
				},
			},
		}, nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(result)
		a.IsTrue(result == "50")
	}

	{
		result, err := threshold.Eval(map[string]interface{}{
			"data": maps.Map{
				"version": "1.0.25",
				"hello": maps.Map{
					"world": []string{"3", "2"},
				},
			},
		}, nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(result)
		a.IsTrue(result == "150")
	}
}

func TestThreshold_Array(t *testing.T) {
	a := assert.NewAssertion(t)

	{
		result, err := EvalParam("${0.a.b.0.d}", []maps.Map{
			{
				"a": maps.Map{
					"b": []interface{}{
						maps.Map{
							"d": "123",
						},
					},
				},
			},
		}, nil, nil, true)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(result)
		a.IsTrue(result == "123")
	}

	{
		result, err := EvalParam("${0.a.b.c}", []maps.Map{
			{
				"a": maps.Map{
					"b": []interface{}{
						maps.Map{
							"d": "123",
						},
					},
				},
			},
		}, nil, nil, true)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(result)
		a.IsTrue(result == "")
	}
}

func TestThreshold_Eval_Date(t *testing.T) {
	threshold := NewThreshold()
	threshold.Param = "new Date().getTime() / 1000 - ${timestamp}"
	threshold.Operator = ThresholdOperatorGt
	err := threshold.Validate()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(threshold.Eval(map[string]interface{}{
		"timestamp": time.Now().Unix() - 10,
	}, nil))
}

func TestThreshold_Eval_Javascript(t *testing.T) {
	threshold := NewThreshold()
	threshold.Param = "javascript:new Date().getTime() / 1000 - ${timestamp}"
	t.Log(threshold.Eval(map[string]interface{}{
		"timestamp": time.Now().Unix() - 10,
	}, nil))
}

func TestThreshold_Eval_Dollar(t *testing.T) {
	a := assert.NewAssertion(t)

	threshold := NewThreshold()
	threshold.Param = "${a.$.percent}"
	threshold.Operator = ThresholdOperatorGt
	threshold.Value = "81"
	err := threshold.Validate()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("should loop:", threshold.shouldLoop, threshold.loopVar)
	a.IsTrue(threshold.TestRow(maps.Map{
		"a": []maps.Map{
			{
				"name":    "30",
				"percent": 30,
			},
			{
				"name":    "60",
				"percent": 60,
			},
			{
				"name":    "82",
				"percent": 82,
			},
			{
				"name":    "50",
				"percent": 50,
			},
		},
	}, nil))
	a.IsFalse(threshold.TestRow(maps.Map{
		"a": []maps.Map{
			{
				"name":    "30",
				"percent": 30,
			},
			{
				"name":    "60",
				"percent": 60,
			},
			{
				"name":    "50",
				"percent": 50,
			},
		},
	}, nil))

	a.IsFalse(threshold.TestRow("abc", nil))
}

func TestThreshold_Eval_Dollar2(t *testing.T) {
	a := assert.NewAssertion(t)

	threshold := NewThreshold()
	threshold.Param = "${$.percent}"
	threshold.Operator = ThresholdOperatorGt
	threshold.Value = "81"
	err := threshold.Validate()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("should loop:", threshold.shouldLoop, threshold.loopVar)
	a.IsTrue(threshold.Test([]maps.Map{
		{
			"percent": 30,
		},
		{
			"percent": 60,
		},
		{
			"percent": 82,
		},
		{
			"percent": 50,
		},
	}, nil))
	a.IsFalse(threshold.Test([]maps.Map{
		{
			"percent": 30,
		},
		{
			"percent": 60,
		},
		{
			"percent": 80,
		},
		{
			"percent": 50,
		},
	}, nil))
}

func TestThreshold_Eval_Dollar3(t *testing.T) {
	a := assert.NewAssertion(t)

	threshold := NewThreshold()
	threshold.Param = "${$}"
	threshold.Operator = ThresholdOperatorGt
	threshold.Value = "3"
	err := threshold.Validate()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("should loop:", threshold.shouldLoop, threshold.loopVar)
	a.IsTrue(threshold.Test([]int{1, 2, 3, 4}, nil))
	a.IsFalse(threshold.Test([]int{1, 2, 3}, nil))
}

func TestThreshold_Eval_Nil(t *testing.T) {
	a := assert.NewAssertion(t)

	threshold := NewThreshold()
	threshold.Param = "${0}"
	threshold.Operator = ThresholdOperatorGte
	threshold.Value = "0"
	err := threshold.Validate()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("should loop:", threshold.shouldLoop, threshold.loopVar)
	result, err := threshold.Test(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	a.IsTrue(result == true)
}

func TestThreshold_Old(t *testing.T) {
	a := assert.NewAssertion(t)

	threshold := NewThreshold()
	threshold.Param = "${rows} - ${OLD.rows234}"
	threshold.Operator = ThresholdOperatorEq
	err := threshold.Validate()
	if err != nil {
		t.Fatal(err)
	}
	result, err := threshold.Eval(map[string]interface{}{
		"rows": 1,
	}, map[string]interface{}{
		"rows234": 123,
	})
	if err != nil {
		t.Fatal(err)
	}
	a.IsTrue(result == "-122")
}

func TestThreshold_Old2(t *testing.T) {
	a := assert.NewAssertion(t)

	threshold := NewThreshold()
	threshold.Param = "Math.abs(${0} - ${OLD})"
	threshold.Operator = ThresholdOperatorEq
	threshold.Value = "333"
	err := threshold.Validate()
	if err != nil {
		t.Fatal(err)
	}
	a.IsTrue(threshold.Test(123, 456))
	a.IsFalse(threshold.Test(123, 455))
}

func TestThreshold_RunActions(t *testing.T) {
	threshold := NewThreshold()
	threshold.Actions = []map[string]interface{}{
		{
			"code": "script",
			"options": map[string]interface{}{
				"scriptType": "path",
				"path":       "1",
			},
		},
	}
	t.Log(threshold.RunActions(nil))
}

func TestThresholdIP(t *testing.T) {
	a := assert.NewAssertion(t)

	{
		th := Threshold{
			Param:    "${0}",
			Operator: ThresholdOperatorEqIP,
			Value:    "hello",
		}
		a.IsNotNil(th.Validate())
		a.IsFalse(th.Test("hello", nil))
	}

	{
		th := Threshold{
			Param:    "${0}",
			Operator: ThresholdOperatorEqIP,
			Value:    "hello",
		}
		a.IsNotNil(th.Validate())
		a.IsFalse(th.Test("192.168.1.100", nil))
	}

	{
		th := Threshold{
			Param:    "${0}",
			Operator: ThresholdOperatorEqIP,
			Value:    "192.168.1.100",
		}
		a.IsNil(th.Validate())
		a.IsTrue(th.Test("192.168.1.100", nil))
	}

	{
		th := Threshold{
			Param:    "${0}",
			Operator: ThresholdOperatorGtIP,
			Value:    "192.168.1.90",
		}
		a.IsNil(th.Validate())
		a.IsTrue(th.Test("192.168.1.100", nil))
	}

	{
		th := Threshold{
			Param:    "${0}",
			Operator: ThresholdOperatorGteIP,
			Value:    "192.168.1.90",
		}
		a.IsNil(th.Validate())
		a.IsTrue(th.Test("192.168.1.100", nil))
	}

	{
		th := Threshold{
			Param:    "${0}",
			Operator: ThresholdOperatorLtIP,
			Value:    "192.168.1.90",
		}
		a.IsNil(th.Validate())
		a.IsTrue(th.Test("192.168.1.80", nil))
	}

	{
		th := Threshold{
			Param:    "${0}",
			Operator: ThresholdOperatorLteIP,
			Value:    "192.168.1.90",
		}
		a.IsNil(th.Validate())
		a.IsFalse(th.Test("192.168.1.100", nil))
	}

	{
		th := Threshold{
			Param:    "${0}",
			Operator: ThresholdOperatorIPRange,
			Value:    "192.168.0.90,",
		}
		a.IsNil(th.Validate())
		a.IsTrue(th.Test("192.168.1.100", nil))
	}

	{
		th := Threshold{
			Param:    "${0}",
			Operator: ThresholdOperatorIPRange,
			Value:    "192.168.0.90,192.168.1.100",
		}
		a.IsNil(th.Validate())
		a.IsTrue(th.Test("192.168.1.100", nil))
	}

	{
		th := Threshold{
			Param:    "${0}",
			Operator: ThresholdOperatorIPRange,
			Value:    ",192.168.1.100",
		}
		a.IsNil(th.Validate())
		a.IsTrue(th.Test("192.168.1.100", nil))
	}

	{
		th := Threshold{
			Param:    "${0}",
			Operator: ThresholdOperatorIPRange,
			Value:    "192.168.0.90,192.168.1.99",
		}
		a.IsNil(th.Validate())
		a.IsFalse(th.Test("192.168.1.100", nil))
	}

	{
		th := Threshold{
			Param:    "${0}",
			Operator: ThresholdOperatorIPRange,
			Value:    "192.168.0.90/24",
		}
		a.IsNil(th.Validate())
		a.IsFalse(th.Test("192.168.1.100", nil))
	}

	{
		th := Threshold{
			Param:    "${0}",
			Operator: ThresholdOperatorIPRange,
			Value:    "192.168.0.90/18",
		}
		a.IsNil(th.Validate())
		a.IsTrue(th.Test("192.168.1.100", nil))
	}

	{
		th := Threshold{
			Param:    "${0}",
			Operator: ThresholdOperatorIPRange,
			Value:    "a/18",
		}
		a.IsNotNil(th.Validate())
		a.IsFalse(th.Test("192.168.1.100", nil))
	}
}

func TestThreshold_Version(t *testing.T) {
	a := assert.NewAssertion(t)

	{
		th := Threshold{
			Param:    "1.0",
			Operator: ThresholdOperatorVersionRange,
			Value:    `1.0,1.1`,
		}
		a.IsNil(th.Validate())
		a.IsTrue(th.Test(nil, nil))
	}

	{
		th := Threshold{
			Param:    "1.0",
			Operator: ThresholdOperatorVersionRange,
			Value:    `1.0,`,
		}
		a.IsNil(th.Validate())
		a.IsTrue(th.Test(nil, nil))
	}

	{
		th := Threshold{
			Param:    "1.0",
			Operator: ThresholdOperatorVersionRange,
			Value:    `,1.1`,
		}
		a.IsNil(th.Validate())
		a.IsTrue(th.Test(nil, nil))
	}

	{
		th := Threshold{
			Param:    "0.9",
			Operator: ThresholdOperatorVersionRange,
			Value:    `1.0,1.1`,
		}
		a.IsNil(th.Validate())
		a.IsFalse(th.Test(nil, nil))
	}

	{
		th := Threshold{
			Param:    "0.9",
			Operator: ThresholdOperatorVersionRange,
			Value:    `1.0`,
		}
		a.IsNil(th.Validate())
		a.IsFalse(th.Test(nil, nil))
	}

	{
		th := Threshold{
			Param:    "1.1",
			Operator: ThresholdOperatorVersionRange,
			Value:    `1.0`,
		}
		a.IsNil(th.Validate())
		a.IsTrue(th.Test(nil, nil))
	}
}
