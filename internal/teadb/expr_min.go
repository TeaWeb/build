package teadb

// 最小值表达式
type MinExpr struct {
	Field string
}

func NewMinExpr(field string) *MinExpr {
	return &MinExpr{
		Field: field,
	}
}
