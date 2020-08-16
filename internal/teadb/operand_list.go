package teadb

import "github.com/iwind/TeaGo/lists"

// 字段操作符列表
type OperandList struct {
	operandMap map[string][]*Operand
	fields     []string // 字段列表，用来保证查询的排序
}

// 获取新对象
func NewOperandList() *OperandList {
	return &OperandList{
		operandMap: map[string][]*Operand{},
	}
}

// 添加字段操作符
func (this *OperandList) Add(field string, operand ...*Operand) *OperandList {
	if len(operand) == 0 {
		return this
	}
	operands, ok := this.operandMap[field]
	if ok {
		operands = append(operands, operand...)
	} else {
		operands = operand
	}
	this.operandMap[field] = operands
	if !lists.ContainsString(this.fields, field) {
		this.fields = append(this.fields, field)
	}
	return this
}

// 所有字段数量
func (this *OperandList) Len() int {
	return len(this.operandMap)
}

// 循环所有字段
func (this *OperandList) Range(f func(field string, operands []*Operand)) {
	for _, field := range this.fields {
		f(field, this.operandMap[field])
	}
}
