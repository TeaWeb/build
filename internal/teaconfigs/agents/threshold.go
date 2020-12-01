package agents

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/rands"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/string"
	"net"
	"regexp"
	"strings"
)

// 参数值正则
var RegexpParamNamedVariable = regexp.MustCompile(`(?U)\${[$\w\s.-]+(\s*\|.+)?}`)

// 阈值定义
type Threshold struct {
	Id            string                   `yaml:"id" json:"id"`                       // ID
	Param         string                   `yaml:"param" json:"param"`                 // 参数
	Operator      ThresholdOperator        `yaml:"operator" json:"operator"`           // 运算符
	Value         string                   `yaml:"value" json:"value"`                 // 对比值
	NoticeLevel   notices.NoticeLevel      `yaml:"noticeLevel" json:"noticeLevel"`     // 通知级别
	NoticeMessage string                   `yaml:"noticeMessage" json:"noticeMessage"` // 通知消息
	Actions       []map[string]interface{} `yaml:"actions" json:"actions"`             // 动作配置
	MaxFails      int                      `yaml:"maxFails" json:"maxFails"`           // 连续失败次数

	isIP bool

	regValue   *regexp.Regexp
	floatValue float64
	ipValue    net.IP

	supportsMath bool

	shouldLoop bool   // 是否应该循环测试，如果包含名为$（dollar符号）的变量，则表示是循环测试
	loopVar    string // 要循环的变量
}

// 新阈值对象
func NewThreshold() *Threshold {
	return &Threshold{
		Id: rands.HexString(16),
	}
}

// 校验
func (this *Threshold) Validate() error {
	this.supportsMath = false
	if this.Operator == ThresholdOperatorRegexp || this.Operator == ThresholdOperatorNotRegexp {
		reg, err := regexp.Compile(this.Value)
		if err != nil {
			return err
		}
		this.regValue = reg
	} else if this.Operator == ThresholdOperatorGt ||
		this.Operator == ThresholdOperatorGte ||
		this.Operator == ThresholdOperatorLt ||
		this.Operator == ThresholdOperatorLte ||
		this.Operator == ThresholdOperatorNumberEq {
		this.floatValue = types.Float64(this.Value)
		this.supportsMath = true
	} else if this.Operator == ThresholdOperatorEq {
		this.supportsMath = true // 为了兼容以前版本的此处必须为true
	} else if lists.ContainsString([]string{
		ThresholdOperatorEqIP,
		ThresholdOperatorGtIP,
		ThresholdOperatorGteIP,
		ThresholdOperatorLtIP,
		ThresholdOperatorLteIP,
	}, this.Operator) {
		this.ipValue = net.ParseIP(this.Value)
		this.isIP = this.ipValue != nil

		if !this.isIP {
			return errors.New("value should be a valid ip")
		}
	} else if lists.ContainsString([]string{
		ThresholdOperatorIPRange,
	}, this.Operator) {
		if strings.Contains(this.Value, ",") {
			ipList := strings.SplitN(this.Value, ",", 2)
			ipString1 := strings.TrimSpace(ipList[0])
			ipString2 := strings.TrimSpace(ipList[1])

			if len(ipString1) > 0 {
				ip1 := net.ParseIP(ipString1)
				if ip1 == nil {
					return errors.New("start ip is invalid")
				}
			}

			if len(ipString2) > 0 {
				ip2 := net.ParseIP(ipString2)
				if ip2 == nil {
					return errors.New("end ip is invalid")
				}
			}
		} else if strings.Contains(this.Value, "/") {
			_, _, err := net.ParseCIDR(this.Value)
			if err != nil {
				return err
			}
		} else {
			return errors.New("invalid ip range")
		}
	}

	// 检查参数值
	for _, v := range RegexpParamNamedVariable.FindAllStringSubmatch(this.Param, -1) {
		varName := v[0][2 : len(v[0])-1]
		pieces := strings.Split(varName, ".")
		if lists.ContainsString(pieces, "$") {
			this.shouldLoop = true
			this.loopVar = varName
			break
		}
	}

	return nil
}

// 将此条件应用于阈值，检查是否匹配
func (this *Threshold) Test(value interface{}, oldValue interface{}) (ok bool, err error) {
	ok, _, err = this.testParam(this.Param, this.shouldLoop, value, oldValue)
	return
}

// 将此条件应用于阈值，检查是否匹配，如果匹配同时也返回$匹配的行数据
func (this *Threshold) TestRow(value interface{}, oldValue interface{}) (ok bool, row interface{}, err error) {
	return this.testParam(this.Param, this.shouldLoop, value, oldValue)
}

// 检查阈值，但指定更多的参数
func (this *Threshold) testParam(param string, shouldLoop bool, value interface{}, oldValue interface{}) (ok bool, row interface{}, err error) {
	// 处理$（dollar符号）
	if shouldLoop {
		pieces := strings.Split(this.loopVar, ".")
		dollarIndex := 0
		for index, piece := range pieces {
			if piece == "$" {
				dollarIndex = index
				break
			}
		}

		if dollarIndex == 0 {
			if types.IsSlice(value) {
				lists.Each(value, func(k int, v interface{}) {
					indexParam := fmt.Sprintf("%d", k)
					if len(pieces) > 1 {
						indexParam += "." + strings.Join(pieces[dollarIndex+1:], ".")
					}
					newParam := strings.Replace(param, "${"+this.loopVar+"}", "${"+indexParam+"}", -1)
					ok1, _, err1 := this.testParam(newParam, false, value, oldValue)
					if ok1 {
						ok = ok1
						err = err1
					}
				})
				return
			}
		} else {
			newValue := teautils.Get(value, pieces[:dollarIndex])
			if types.IsSlice(newValue) {
				lists.Each(newValue, func(k int, v interface{}) {
					indexParam := strings.Join(pieces[:dollarIndex], ".") + "." + fmt.Sprintf("%d.", k) + strings.Join(pieces[dollarIndex+1:], ".")

					newParam := strings.Replace(param, "${"+this.loopVar+"}", "${"+indexParam+"}", -1)
					ok1, _, err1 := this.testParam(newParam, false, value, oldValue)
					if ok1 {
						ok = ok1
						err = err1
						row = v
					}
				})
				return
			}
		}

		return false, nil, nil
	}

	paramValue, err := EvalParam(param, value, oldValue, nil, this.supportsMath)
	if err != nil {
		return false, nil, err
	}

	switch this.Operator {
	case ThresholdOperatorRegexp:
		if this.regValue == nil {
			return false, nil, nil
		}
		return this.regValue.MatchString(types.String(paramValue)), nil, nil
	case ThresholdOperatorNotRegexp:
		if this.regValue == nil {
			return false, nil, nil
		}
		return !this.regValue.MatchString(types.String(paramValue)), nil, nil
	case ThresholdOperatorGt:
		return types.Float64(paramValue) > this.floatValue, nil, nil
	case ThresholdOperatorGte:
		return types.Float64(paramValue) >= this.floatValue, nil, nil
	case ThresholdOperatorLt:
		return types.Float64(paramValue) < this.floatValue, nil, nil
	case ThresholdOperatorLte:
		return types.Float64(paramValue) <= this.floatValue, nil, nil
	case ThresholdOperatorEq:
		return paramValue == this.Value, nil, nil
	case ThresholdOperatorNumberEq:
		return types.Float64(paramValue) == this.floatValue, nil, nil
	case ThresholdOperatorNot:
		return paramValue != this.Value, nil, nil
	case ThresholdOperatorPrefix:
		return strings.HasPrefix(types.String(paramValue), this.Value), nil, nil
	case ThresholdOperatorSuffix:
		return strings.HasSuffix(types.String(paramValue), this.Value), nil, nil
	case ThresholdOperatorContains:
		return strings.Contains(types.String(paramValue), this.Value), nil, nil
	case ThresholdOperatorNotContains:
		return !strings.Contains(types.String(paramValue), this.Value), nil, nil
	case ThresholdOperatorVersionRange:
		if strings.Contains(this.Value, ",") {
			versions := strings.SplitN(this.Value, ",", 2)
			version1 := strings.TrimSpace(versions[0])
			version2 := strings.TrimSpace(versions[1])
			if len(version1) > 0 && stringutil.VersionCompare(paramValue, version1) < 0 {
				return false, nil, nil
			}
			if len(version2) > 0 && stringutil.VersionCompare(paramValue, version2) > 0 {
				return false, nil, nil
			}
			return true, nil, nil
		} else {
			return stringutil.VersionCompare(paramValue, this.Value) >= 0, nil, nil
		}
	case ThresholdOperatorEqIP:
		ip := net.ParseIP(paramValue)
		if ip == nil {
			return false, nil, nil
		}
		return this.isIP && bytes.Compare(this.ipValue, ip) == 0, nil, nil
	case ThresholdOperatorGtIP:
		ip := net.ParseIP(paramValue)
		if ip == nil {
			return false, nil, nil
		}
		return this.isIP && bytes.Compare(ip, this.ipValue) > 0, nil, nil
	case ThresholdOperatorGteIP:
		ip := net.ParseIP(paramValue)
		if ip == nil {
			return false, nil, nil
		}
		return this.isIP && bytes.Compare(ip, this.ipValue) >= 0, nil, nil
	case ThresholdOperatorLtIP:
		ip := net.ParseIP(paramValue)
		if ip == nil {
			return false, nil, nil
		}
		return this.isIP && bytes.Compare(ip, this.ipValue) < 0, nil, nil
	case ThresholdOperatorLteIP:
		ip := net.ParseIP(paramValue)
		if ip == nil {
			return false, nil, nil
		}
		return this.isIP && bytes.Compare(ip, this.ipValue) <= 0, nil, nil
	case ThresholdOperatorIPRange:
		ip := net.ParseIP(paramValue)
		if ip == nil {
			return false, nil, nil
		}

		// 检查IP范围格式
		if strings.Contains(this.Value, ",") {
			ipList := strings.SplitN(this.Value, ",", 2)
			ipString1 := strings.TrimSpace(ipList[0])
			ipString2 := strings.TrimSpace(ipList[1])

			if len(ipString1) > 0 {
				ip1 := net.ParseIP(ipString1)
				if ip1 == nil {
					return false, nil, nil
				}

				if bytes.Compare(ip, ip1) < 0 {
					return false, nil, nil
				}
			}

			if len(ipString2) > 0 {
				ip2 := net.ParseIP(ipString2)
				if ip2 == nil {
					return false, nil, nil
				}

				if bytes.Compare(ip, ip2) > 0 {
					return false, nil, nil
				}
			}

			return true, nil, nil
		} else if strings.Contains(this.Value, "/") {
			_, ipNet, err := net.ParseCIDR(this.Value)
			if err != nil {
				return false, nil, nil
			}
			return ipNet.Contains(ip), nil, nil
		} else {
			return false, nil, nil
		}
	}
	return false, nil, nil
}

// 执行数值运算，使用Javascript语法
func (this *Threshold) Eval(value interface{}, old interface{}) (string, error) {
	return EvalParam(this.Param, value, old, nil, this.supportsMath)
}

// 执行动作
func (this *Threshold) RunActions(params map[string]string) error {
	if len(this.Actions) == 0 {
		return nil
	}

	for _, a := range this.Actions {
		code, found := a["code"]
		if !found {
			return errors.New("action 'code' not found")
		}

		options, found := a["options"]
		if !found {
			return errors.New("action 'options' not found")
		}
		optionsMap, ok := options.(map[string]interface{})
		if !ok {
			return errors.New("action 'options' should be a valid map")
		}

		action := FindAction(types.String(code))
		if action == nil {
			return errors.New("action for '" + types.String(code) + "' not found")
		}

		instance := action["instance"]
		err := teautils.MapToObjectJSON(optionsMap, &instance)
		if err != nil {
			return err
		}

		output, err := instance.(ActionInterface).Run(params)
		if err != nil {
			return err
		}
		if len(output) > 0 {
			logs.Println("[threshold]run actions:", output)
		}
	}

	return nil
}

// 取得描述文本
func (this *Threshold) Expression() string {
	return this.Param + " " + this.Operator + " " + this.Value
}
