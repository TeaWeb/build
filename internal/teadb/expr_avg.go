package teadb

// 平均值表达式
type AvgExpr struct {
	Field string
}

func NewAvgExpr(field string) *AvgExpr {
	return &AvgExpr{
		Field: field,
	}
}
