package teadb

// 最大值表达式
type MaxExpr struct {
	Field string
}

func NewMaxExpr(field string) *MaxExpr {
	return &MaxExpr{
		Field: field,
	}
}
