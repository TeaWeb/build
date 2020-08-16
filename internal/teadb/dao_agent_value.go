package teadb

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
)

// Agent数值记录DAO
type AgentValueDAOInterface interface {
	// 设置驱动
	SetDriver(driver DriverInterface)

	// 初始化
	Init()

	// 获取表格
	TableName(agentId string) string

	// 插入数据
	Insert(agentId string, value *agents.Value) error

	// 清除数值
	ClearItemValues(agentId string, appId string, itemId string, level notices.NoticeLevel) error

	// 查找最近的一条记录
	FindLatestItemValue(agentId string, appId string, itemId string) (*agents.Value, error)

	// 查找最近的一条非错误的记录
	FindLatestItemValueNoError(agentId string, appId string, itemId string) (*agents.Value, error)

	// 取得最近的数值记录
	FindLatestItemValues(agentId string, appId string, itemId string, noticeLevel notices.NoticeLevel, lastId string, size int) ([]*agents.Value, error)

	// 列出数值
	ListItemValues(agentId string, appId string, itemId string, noticeLevel notices.NoticeLevel, lastId string, offset int, size int) ([]*agents.Value, error)

	// 分组查询
	QueryValues(query *Query) ([]*agents.Value, error)

	// 根据时间对值进行分组查询
	GroupValuesByTime(query *Query, timeField string, result map[string]Expr) ([]*agents.Value, error)

	// 删除Agent相关表
	DropAgentTable(agentId string) error
}
