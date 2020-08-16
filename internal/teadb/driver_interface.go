package teadb

import "github.com/iwind/TeaGo/maps"

// 数据库驱动
type DriverInterface interface {
	// 初始化
	Init()

	// 设置是否可用
	SetIsAvailable(b bool)

	// 取得是否可用
	IsAvailable() bool

	// 查找单条记录
	FindOne(query *Query, modelPtr interface{}) (interface{}, error)

	// 查找多条记录
	FindOnes(query *Query, modelPtr interface{}) ([]interface{}, error)

	// 插入一条记录
	InsertOne(table string, modelPtr interface{}) error

	// 插入多条记录
	InsertOnes(table string, modelPtrSlice interface{}) error

	// 删除多条记录
	DeleteOnes(query *Query) error

	// 计算总数量
	Count(query *Query) (int64, error)

	// 计算总和
	Sum(query *Query, field string) (float64, error)

	// 计算平均值
	Avg(query *Query, field string) (float64, error)

	// 计算最小值
	Min(query *Query, field string) (float64, error)

	// 计算最大值
	Max(query *Query, field string) (float64, error)

	// 对数据进行分组统计
	Group(query *Query, field string, result map[string]Expr) ([]maps.Map, error)

	// 测试数据库连接
	Test() error

	// 关闭
	Shutdown() error

	// 列出所有表
	ListTables() ([]string, error)

	// 统计数据表信息
	StatTables(tables []string) (map[string]*TableStat, error)

	// 删除表
	DropTable(table string) error
}
