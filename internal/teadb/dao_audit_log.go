package teadb

import "github.com/TeaWeb/build/internal/teaconfigs/audits"

// 审计日志DAO
type AuditLogDAOInterface interface {
	// 设置驱动
	SetDriver(driver DriverInterface)

	// 初始化
	Init()

	// 计算审计日志数量
	CountAllAuditLogs() (int64, error)

	// 列出审计日志
	ListAuditLogs(offset int, size int) ([]*audits.Log, error)

	// 插入一条审计日志
	InsertOne(auditLog *audits.Log) error
}
