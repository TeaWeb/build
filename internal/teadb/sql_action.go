package teadb

// SQL动作类型
type SQLAction = int

// SQL动作列表
const (
	SQLInsert SQLAction = 1
	SQLSelect SQLAction = 2
	SQLDelete SQLAction = 3
	SQLUpdate SQLAction = 4
)
