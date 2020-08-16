package teadb

// SQL相关接口
type SQLDriverInterface interface {
	// 创建表格
	CreateTable(table string, definitionSQL string) error

	// 修改多条记录
	UpdateOnes(query *Query, values map[string]interface{}) error

	// 读取JSON字段
	JSONExtract(field string, path string) string
}
