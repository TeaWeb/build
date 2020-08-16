package teadb

import (
	"github.com/TeaWeb/build/internal/teaconfigs/audits"
	"github.com/iwind/TeaGo/logs"
)

type SQLAuditLogDAO struct {
	BaseDAO
}

// 初始化
func (this *SQLAuditLogDAO) Init() {

}

func (this *SQLAuditLogDAO) TableName() string {
	table := "teaweb_logs_audit"
	this.initTable(table)
	return table
}

// 计算审计日志数量
func (this *SQLAuditLogDAO) CountAllAuditLogs() (int64, error) {
	return NewQuery(this.TableName()).Count()
}

// 列出审计日志
func (this *SQLAuditLogDAO) ListAuditLogs(offset int, size int) ([]*audits.Log, error) {
	ones, err := NewQuery(this.TableName()).
		Offset(offset).
		Limit(size).
		Desc("_id").
		FindOnes(new(audits.Log))
	if err != nil {
		return nil, err
	}
	result := []*audits.Log{}
	for _, one := range ones {
		result = append(result, one.(*audits.Log))
	}
	return result, nil
}

// 插入一条审计日志
func (this *SQLAuditLogDAO) InsertOne(auditLog *audits.Log) error {
	return this.driver.InsertOne(this.TableName(), auditLog)
}

// 初始化表格
func (this *SQLAuditLogDAO) initTable(table string) {
	if isInitializedTable(table) {
		return
	}

	logs.Println("[db]check table '" + table + "'")

	switch sharedDBType {
	case "mysql":
		err := this.driver.(SQLDriverInterface).CreateTable(table, "CREATE TABLE `"+table+"` ("+
			"`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',"+
			"`_id` varchar(24) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'global id',"+
			"`action` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,"+
			"`username` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,"+
			"`description` varchar(1024) COLLATE utf8mb4_bin DEFAULT NULL,"+
			"`options` json DEFAULT NULL,"+
			"`timestamp` int(11) DEFAULT NULL,"+
			"PRIMARY KEY (`id`),"+
			"UNIQUE KEY `_id` (`_id`)"+
			") ENGINE=InnoDB AUTO_INCREMENT=100013 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;")
		if err != nil {
			logs.Error(err)
			removeInitializedTable(table)
		}
	case "postgres":
		err := this.driver.(SQLDriverInterface).CreateTable(table, `CREATE TABLE "public"."`+table+`" (
		"id" serial8 primary key,
		"_id" varchar(24),
		"action" varchar(255),
		"username" varchar(255),
		"description" varchar(1024),
		"options" json,
		"timestamp" int4 DEFAULT 0
		)
		;

		CREATE UNIQUE INDEX "`+table+`_id" ON "public"."`+table+`" ("_id");`)
		if err != nil {
			logs.Error(err)
			removeInitializedTable(table)
		}
	}
}
