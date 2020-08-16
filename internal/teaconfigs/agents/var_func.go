package agents

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/types"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// func列表，code => type
var funcMap = map[string]reflect.Value{
	"float":  reflect.ValueOf(FuncFloat),
	"round":  reflect.ValueOf(FuncRound),
	"ceil":   reflect.ValueOf(FuncCeil),
	"floor":  reflect.ValueOf(FuncFloor),
	"format": reflect.ValueOf(FuncFormat),
	"append": reflect.ValueOf(FuncAppend),
}

// 浮点数处理
// float | float('%.2f')
func FuncFloat(args ...interface{}) interface{} {
	if len(args) == 0 {
		return 0
	}
	if len(args) == 1 {
		return types.Float64(args[0])
	}
	args[0] = types.Float64(args[0])
	return FuncFormat(args...)
}

// 对数字四舍五入
// round | round(2)
func FuncRound(args ...interface{}) interface{} {
	if len(args) == 0 {
		return 0
	}
	if len(args) == 1 {
		return int64(math.Round(types.Float64(args[0])))
	}

	precision := types.Int64(args[1])
	if precision <= 0 {
		return int64(math.Round(types.Float64(args[0])))
	}
	return fmt.Sprintf("%."+strconv.FormatInt(precision, 10)+"f", types.Float64(args[0]))
}

// 对数字进行取不小于它的整数
// ceil
func FuncCeil(args ...interface{}) interface{} {
	if len(args) == 0 {
		return 0
	}
	return int64(math.Ceil(types.Float64(args[0])))
}

// 对数字进行取不大于它的整数
// floor
func FuncFloor(args ...interface{}) interface{} {
	if len(args) == 0 {
		return 0
	}
	return int64(math.Floor(types.Float64(args[0])))
}

// 格式化
// format('%.2f')
func FuncFormat(args ...interface{}) interface{} {
	if len(args) == 0 {
		return ""
	}
	if len(args) == 1 {
		return types.String(args[0])
	}
	format := types.String(args[1])
	if len(args) == 2 {
		return fmt.Sprintf(format, args[0])
	} else {
		return fmt.Sprintf(format, append([]interface{}{args[0]}, args[2:]...)...)
	}
}

// 附加字符串
// append('a', 'b')
func FuncAppend(args ...interface{}) interface{} {
	s := strings.Builder{}
	for _, arg := range args {
		s.WriteString(types.String(arg))
	}
	return s.String()
}

// 执行函数表达式
func RunFuncExpr(value interface{}, expr []byte) (interface{}, error) {
	// 防止因nil参数导致panic
	if value == nil {
		value = ""
	}

	// 空表达式直接返回
	if len(expr) == 0 || len(bytes.TrimSpace(expr)) == 0 {
		return value, nil
	}

	identifier := []byte{}

	isInQuote := false
	isQuoted := false
	quoteByte := byte(0)
	funcCode := ""
	args := []interface{}{value}

	for index, b := range expr {
		switch b {
		case '|':
			if isInQuote {
				identifier = append(identifier, b)
			} else { // end function
				if len(funcCode) == 0 {
					funcCode = string(identifier)
				} else if len(identifier) > 0 {
					return value, errors.New("invalid identifier '" + string(identifier) + "'")
				}
				value, err := callFunc(funcCode, args)
				if err != nil {
					return value, err
				}
				return RunFuncExpr(value, expr[index+1:])
			}
		case '(':
			if isInQuote {
				identifier = append(identifier, b)
			} else {
				funcCode = string(identifier)
				identifier = []byte{}
			}
		case ',', ')':
			if isInQuote {
				identifier = append(identifier, b)
			} else {
				if isQuoted {
					isQuoted = false
					args = append(args, string(identifier))
				} else {
					arg, err := checkLiteral(string(identifier))
					if err != nil {
						return value, err
					}
					args = append(args, arg)
				}
				identifier = []byte{}
			}
		case '\\':
			if isInQuote && (index == len(expr)-1 || expr[index+1] != quoteByte) {
				identifier = append(identifier, b)
			} else {
				continue
			}
		case '\'', '"':
			if isInQuote {
				if quoteByte == b && expr[index-1] != '\\' {
					isInQuote = false
					break
				}
			} else if index == 0 || expr[index-1] != '\\' {
				isInQuote = true
				isQuoted = true
				quoteByte = b
				break
			}
			identifier = append(identifier, b)
		case ' ', '\t':
			if isInQuote {
				identifier = append(identifier, b)
			}
		default:
			identifier = append(identifier, b)
		}
	}

	if len(funcCode) == 0 {
		funcCode = string(identifier)
	} else if len(identifier) > 0 {
		return value, errors.New("invalid identifier '" + string(identifier) + "'")
	}

	return callFunc(funcCode, args)
}

// 注册一个函数
func RegisterFunc(code string, f interface{}) {
	funcMap[code] = reflect.ValueOf(f)
}

// 调用一个函数
func callFunc(funcCode string, args []interface{}) (value interface{}, err error) {
	fn, ok := funcMap[funcCode]
	if !ok {
		return value, errors.New("function '" + funcCode + "' not found")
	}
	argValues := []reflect.Value{}
	for _, arg := range args {
		argValues = append(argValues, reflect.ValueOf(arg))
	}
	result := fn.Call(argValues)
	if len(result) == 0 {
		value = nil
	} else {
		value = result[0].Interface()
	}
	return
}

// 检查字面量，支持true, false, null, nil和整数数字、浮点数数字（不支持e）
func checkLiteral(identifier string) (result interface{}, err error) {
	if len(identifier) == 0 {
		return "", nil
	}
	switch strings.ToLower(identifier) {
	case "true":
		result = true
		return
	case "false":
		result = false
		return
	case "null", "nil":
		result = nil
		return
	default:
		if shared.RegexpAllDigitNumber.MatchString(identifier) {
			result = types.Int64(identifier)
			return
		}
		if shared.RegexpAllFloatNumber.MatchString(identifier) {
			result = types.Float64(identifier)
			return
		}
	}

	err = errors.New("undefined identifier '" + identifier + "'")
	return
}
