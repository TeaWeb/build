package teadb

// 总和表达式
type SumExpr struct {
	Field string
}

func NewSumExpr(field string) *SumExpr {
	return &SumExpr{
		Field: field,
	}
}
