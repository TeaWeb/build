package agents

import (
	"errors"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/robertkrimen/otto"
	"strings"
)

// 使用某个参数执行数值运算，使用Javascript语法
func EvalParam(param string, value interface{}, old interface{}, varMapping maps.Map, supportsMath bool) (string, error) {
	if old == nil {
		old = value
	}
	var resultErr error = nil
	paramValue := RegexpParamNamedVariable.ReplaceAllStringFunc(param, func(s string) string {
		varName := s[2 : len(s)-1]
		funcExpr := ""
		index := strings.Index(varName, "|")
		if index > -1 {
			funcExpr = varName[index+1:]
			varName = strings.TrimSpace(varName[:index])
		}

		result, err := evalVarName(varName, value, old, varMapping, supportsMath)
		if err != nil {
			resultErr = err
		} else if len(funcExpr) > 0 {
			result, err = RunFuncExpr(result, []byte(funcExpr))
		}

		return types.String(result)
	})

	// 支持加、减、乘、除、余
	if len(paramValue) > 0 {
		if supportsMath && strings.ContainsAny(param, "+-*/%") {
			vm := otto.New()
			v, err := vm.Run(paramValue)
			if err != nil {
				return "", errors.New("\"" + param + "\": eval \"" + paramValue + "\":" + err.Error())
			} else {
				paramValue = v.String()
			}
		}

		// javascript
		if strings.HasPrefix(paramValue, "javascript:") {
			vm := otto.New()
			v, err := vm.Run(paramValue[len("javascript:"):])
			if err != nil {
				return "", errors.New("\"" + param + "\": eval \"" + paramValue + "\":" + err.Error())
			} else {
				paramValue = v.String()
			}
		}
	}

	return paramValue, resultErr
}

func evalVarName(varName string, value interface{}, old interface{}, varMapping maps.Map, supportsMath bool) (interface{}, error) {
	// 从varMapping中查找
	if varMapping != nil {
		index := strings.Index(varName, ".")
		var firstKey = varName
		if index > 0 {
			firstKey = varName[:index]
		}
		firstKey = strings.TrimSpace(firstKey)
		if varMapping.Has(firstKey) {
			keys := strings.Split(varName, ".")
			for index, key := range keys {
				keys[index] = strings.TrimSpace(key)
			}
			result := teautils.Get(varMapping, keys)
			if result == nil {
				return "", nil
			}

			return result, nil
		}
	}

	if value == nil {
		return nil, nil
	}

	// 支持${OLD}和${OLD.xxx}
	if varName == "OLD" {
		return evalVarName("0", old, nil, nil, supportsMath)
	} else if strings.HasPrefix(varName, "OLD.") {
		return evalVarName(varName[4:], old, nil, nil, supportsMath)
	}

	switch v := value.(type) {
	case string:
		if varName == "0" {
			return v, nil
		}
		return "", nil
	case int8, int16, int, int32, int64, uint8, uint16, uint, uint32, uint64:
		if varName == "0" {
			return v, nil
		}
		return 0, nil
	case float32, float64:
		if varName == "0" {
			return v, nil
		}
		return 0, nil
	case bool:
		if varName == "0" {
			return v, nil
		}
		return false, nil
	default:
		if types.IsSlice(value) || types.IsMap(value) {
			keys := strings.Split(varName, ".")
			for index, key := range keys {
				keys[index] = strings.TrimSpace(key)
			}
			return teautils.Get(v, keys), nil
		}
	}
	return "${" + varName + "}", nil
}
