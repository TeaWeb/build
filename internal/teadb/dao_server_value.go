package teadb

import (
	"github.com/TeaWeb/build/internal/teaconfigs/stats"
	"github.com/TeaWeb/build/internal/teadb/shared"
)

// 代理统计DAO
type ServerValueDAOInterface interface {
	// 设置驱动
	SetDriver(driver DriverInterface)
	
	// 初始化
	Init()

	// 表名
	TableName(serverId string) string

	// 插入新数据
	InsertOne(serverId string, value *stats.Value) error

	// 删除过期的数据
	DeleteExpiredValues(serverId string, period stats.ValuePeriod, life int) error

	// 查询相同的数值记录
	FindSameItemValue(serverId string, item *stats.Value) (*stats.Value, error)

	// 修改值和时间戳
	UpdateItemValueAndTimestamp(serverId string, valueId string, value map[string]interface{}, timestamp int64) error

	// 创建索引
	CreateIndex(serverId string, fields []*shared.IndexField) error

	// 查询数据
	QueryValues(query *Query) ([]*stats.Value, error)

	// 根据item查找一条数据
	FindOneWithItem(serverId string, item string) (*stats.Value, error)

	// 删除代理服务相关表
	DropServerTable(serverId string) error
}
