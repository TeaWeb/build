package teadb

// 操作符类型
type OperandCode = string

// 常用操作符
const (
	OperandEq      OperandCode = "eq"
	OperandLt      OperandCode = "lt"
	OperandLte     OperandCode = "lte"
	OperandGt      OperandCode = "gt"
	OperandGte     OperandCode = "gte"
	OperandIn      OperandCode = "in"
	OperandNotIn   OperandCode = "nin"
	OperandNeq     OperandCode = "ne"
	OperandOr      OperandCode = "or"
	operandSQLCond OperandCode = "sqlCond" // SQL专用
)

type Operand struct {
	Code  OperandCode
	Value interface{}
}

func NewOperand(code OperandCode, value interface{}) *Operand {
	return &Operand{
		Code:  code,
		Value: value,
	}
}
