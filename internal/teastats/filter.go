package teastats

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
)

// 筛选器接口
type FilterInterface interface {
	// 名称
	Name() string

	// 描述
	Description() string

	// 代号
	Codes() []string

	// 参数说明
	ParamVariables() []*Variable

	// 统计数据说明
	ValueVariables() []*Variable

	// 索引参数
	Indexes() []string

	// 启动
	Start(queue *Queue, code string)

	// 筛选某个访问日志
	Filter(accessLog *accesslogs.AccessLog)

	// 提交数据
	Commit()

	// 停止
	Stop()
}
