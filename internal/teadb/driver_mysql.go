package teadb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs/db"
	"github.com/TeaWeb/build/internal/teadb/shared"
	"github.com/TeaWeb/build/internal/teautils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iwind/TeaGo/logs"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"regexp"
	"sort"
	"strings"
	"time"
)

type MySQLDriver struct {
	SQLDriver
}

// 初始化
func (this *MySQLDriver) Init() {
	this.driver = "mysql"

	err := this.initDB()
	if err != nil {
		logs.Error(err)
	}

	agentValueDAO = new(SQLAgentValueDAO)
	agentValueDAO.SetDriver(this)
	agentValueDAO.Init()

	agentLogDAO = new(SQLAgentLogDAO)
	agentLogDAO.SetDriver(this)
	agentLogDAO.Init()

	serverValueDAO = new(SQLServerValueDAO)
	serverValueDAO.SetDriver(this)
	serverValueDAO.Init()

	auditLogDAO = new(SQLAuditLogDAO)
	auditLogDAO.SetDriver(this)
	auditLogDAO.Init()

	accessLogDAO = new(SQLAccessLogDAO)
	accessLogDAO.SetDriver(this)
	accessLogDAO.Init()

	noticeDAO = new(SQLNoticeDAO)
	noticeDAO.SetDriver(this)
	noticeDAO.Init()

	// tasks
	go func() {
		this.cleanAccessLogs()
	}()
}

func (this *MySQLDriver) initDB() error {
	config, err := db.LoadMySQLConfig()
	if err != nil {
		return err
	}
	dbInstance, err := sql.Open("mysql", config.DSN)
	if err != nil {
		return err
	}
	if config.PoolSize > 0 {
		half := config.PoolSize / 2
		dbInstance.SetMaxIdleConns(half)
		dbInstance.SetMaxOpenConns(config.PoolSize)
	} else {
		dbInstance.SetMaxIdleConns(32)
		dbInstance.SetMaxOpenConns(64)
	}
	dbInstance.SetConnMaxLifetime(0)
	this.db = dbInstance

	go func() {
		row := this.db.QueryRow("SELECT @@sql_mode")
		if row != nil {
			_ = row.Scan(&this.sqlMode)
		}
	}()

	return nil
}

// 检查表是否存在
func (this *MySQLDriver) CheckTableExists(table string) (bool, error) {
	currentDB, err := this.checkDB()
	if err != nil {
		return false, err
	}

	_, err = currentDB.ExecContext(context.Background(), "SHOW CREATE TABLE `"+table+"`")
	if err != nil {
		return false, nil
	}

	return true, nil
}

// 创建表
func (this *MySQLDriver) CreateTable(table string, definitionSQL string) error {
	currentDB, err := this.checkDB()
	if err != nil {
		return err
	}

	exists, err := this.CheckTableExists(table)
	if err != nil {
		return err
	}
	if !exists {
		_, err = currentDB.ExecContext(context.Background(), definitionSQL)
		if err != nil {
			logs.Error(err)
		}
	}
	return err
}

// 测试DSN
func (this *MySQLDriver) TestDSN(dsn string, autoCreateDB bool) (message string, ok bool) {
	dbInstance, err := sql.Open("mysql", dsn)
	if err != nil {
		message = "DSN解析错误：" + err.Error()
		return
	}
	defer func() {
		_ = dbInstance.Close()
	}()

	// 检查数据库
	if autoCreateDB {
		index := strings.Index(dsn, "/")
		if index == -1 {
			message = "invalid dsn"
			return
		}
		database := dsn[index+1:]
		index = strings.Index(database, "?")
		if index > -1 {
			database = database[:index]
		}
		if len(database) == 0 {
			message = "no database defined"
			return
		}
		newDSN := strings.Replace(dsn, "/"+database, "/", -1)
		newDBInstance, err := sql.Open("mysql", newDSN)
		if err != nil {
			message = err.Error()
			return
		}
		_, err = newDBInstance.ExecContext(context.Background(), "CREATE DATABASE IF NOT EXISTS `"+database+"`")
		if err != nil {
			message = err.Error()
			_ = newDBInstance.Close()
			return
		}
		_ = newDBInstance.Close()
	}

	// 测试创建数据表
	_, err = dbInstance.ExecContext(context.Background(), "CREATE TABLE `teaweb_test` ( "+
		"`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,"+
		"`_id` varchar(24) DEFAULT NULL,"+
		"PRIMARY KEY (`id`),"+
		"UNIQUE KEY `_id` (`_id`)"+
		" ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
	if err != nil {
		message = "尝试创建数据表失败：" + err.Error()
		return
	}

	// 测试写入数据表
	_, err = dbInstance.ExecContext(context.Background(), "INSERT INTO `teaweb_test` (`_id`) VALUES (\""+shared.NewObjectId().Hex()+"\")")
	if err != nil {
		message = "尝试写入数据表失败：" + err.Error()
		return
	}

	// 测试删除数据表
	_, err = dbInstance.ExecContext(context.Background(), "DROP TABLE `teaweb_test`")
	if err != nil {
		message = "尝试删除数据表失败：" + err.Error()
		return
	}

	// 检查函数
	rows, err := dbInstance.Query(`SELECT JSON_EXTRACT('{"a":1}', "$.a")`)
	if err != nil {
		message = "检查JSON_EXTRACT()函数失败：" + err.Error() + "。请尝试使用MySQL v5.7.8以上版本。"
		return
	}
	_ = rows.Close()

	ok = true
	return
}

// 列出所有表
func (this *MySQLDriver) ListTables() ([]string, error) {
	currentDB, err := this.checkDB()
	if err != nil {
		return nil, err
	}

	rows, err := currentDB.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	tables := []string{}
	for rows.Next() {
		s := ""
		err = rows.Scan(&s)
		if err != nil {
			return nil, err
		}
		if !strings.HasPrefix(s, "teaweb_") {
			continue
		}
		tables = append(tables, s)
	}

	sort.Strings(tables)

	return tables, nil
}

// 统计数据表
func (this *MySQLDriver) StatTables(tables []string) (map[string]*TableStat, error) {
	if len(tables) == 0 {
		return nil, nil
	}

	currentDB, err := this.checkDB()
	if err != nil {
		return nil, err
	}

	holders := []string{}
	args := []interface{}{}
	for _, table := range tables {
		holders = append(holders, "?")
		args = append(args, table)
	}
	rows, err := currentDB.Query("SELECT `table_name`, `table_rows`, `data_length` FROM `INFORMATION_SCHEMA`.`TABLES` WHERE `table_schema`=DATABASE() AND `table_name` IN ("+strings.Join(holders, ", ")+")", args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	result := map[string]*TableStat{}
	for rows.Next() {
		tableName := ""
		tableRows := int64(0)
		tableSize := int64(0)
		err = rows.Scan(&tableName, &tableRows, &tableSize)
		if err != nil {
			return nil, err
		}
		formattedSize := ""
		if tableSize < 1024 {
			formattedSize = fmt.Sprintf("%.2fB", float64(tableSize))
		} else if tableSize < 1024*1024 {
			formattedSize = fmt.Sprintf("%.2fKB", float64(tableSize)/1024)
		} else if tableSize < 1024*1024*1024 {
			formattedSize = fmt.Sprintf("%.2fMB", float64(tableSize)/1024/1024)
		} else {
			formattedSize = fmt.Sprintf("%.2fGB", float64(tableSize)/1024/1024/1024)
		}
		result[tableName] = &TableStat{
			Count:         tableRows,
			Size:          tableSize,
			FormattedSize: formattedSize,
		}
	}

	return result, nil
}

// 清理访问日志任务
func (this *MySQLDriver) cleanAccessLogs() {
	reg := regexp.MustCompile(`^teaweb_logs_(\d{8})$`)

	teautils.Every(1*time.Minute, func(ticker *teautils.Ticker) {
		config, _ := db.LoadMySQLConfig()
		if config == nil {
			return
		}

		if config.AccessLog == nil {
			return
		}

		now := time.Now()
		if config.AccessLog.CleanHour != now.Hour() ||
			now.Minute() != 0 ||
			config.AccessLog.KeepDays < 1 {
			return
		}

		compareDay := "teaweb_logs_" + timeutil.Format("Ymd", time.Now().Add(-time.Duration(config.AccessLog.KeepDays*24)*time.Hour))
		logs.Println("[mysql]clean access logs before '" + compareDay + "'")

		tables, err := this.ListTables()
		if err != nil {
			logs.Println("[mysql]" + err.Error())
			return
		}

		for _, table := range tables {
			if !reg.MatchString(table) {
				continue
			}

			if table < compareDay {
				logs.Println("[mysql]clean table '" + table + "'")
				err = this.DropTable(table)
				if err != nil {
					logs.Println("[mysql]" + err.Error())
				}
			}
		}
	})
}
