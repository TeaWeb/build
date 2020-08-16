package shared

// 索引字段定义
type IndexField struct {
	Name string
	Asc  bool
}

// 创建新索引定义
func NewIndexField(name string, asc bool) *IndexField {
	return &IndexField{
		Name: name,
		Asc:  asc,
	}
}
