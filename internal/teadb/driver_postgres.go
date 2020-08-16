package teadb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/TeaWeb/build/internal/teaconfigs/db"
	"github.com/TeaWeb/build/internal/teadb/shared"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/logs"
	timeutil "github.com/iwind/TeaGo/utils/time"
	_ "github.com/lib/pq"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
)

type PostgresDriver struct {
	SQLDriver
}

// 初始化
func (this *PostgresDriver) Init() {
	this.driver = "postgres"

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

// 初始化数据库
func (this *PostgresDriver) initDB() error {
	config, err := db.LoadPostgresConfig()
	if err != nil {
		return err
	}
	dbInstance, err := sql.Open("postgres", config.DSN)
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
	return nil
}

// 检查表是否存在
func (this *PostgresDriver) CheckTableExists(table string) (bool, error) {
	currentDB, err := this.checkDB()
	if err != nil {
		return false, err
	}

	stmt, err := currentDB.PrepareContext(context.Background(), "SELECT table_name FROM information_schema.tables  WHERE table_name=$1")
	if err != nil {
		return false, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	row := stmt.QueryRow(table)
	i := interface{}(nil)
	err = row.Scan(&i)
	return err == nil, nil
}

// 创建表
func (this *PostgresDriver) CreateTable(table string, definitionSQL string) error {
	currentDB, err := this.checkDB()
	if err != nil {
		return err
	}

	exists, err := this.CheckTableExists(table)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	this.dbLocker.Lock()
	defer this.dbLocker.Unlock()

	_, err = currentDB.ExecContext(context.Background(), definitionSQL)
	return err
}

// 测试DSN
func (this *PostgresDriver) TestDSN(dsn string, autoCreateDB bool) (message string, ok bool) {
	dbInstance, err := sql.Open("postgres", dsn)
	if err != nil {
		message = "DSN解析错误：" + err.Error()
		return
	}
	defer func() {
		_ = dbInstance.Close()
	}()

	if autoCreateDB {
		u, err := url.Parse(dsn)
		if err != nil {
			message = err.Error()
			return
		}
		if len(u.Path) <= 1 {
			message = "database name should not be empty"
			return
		}
		database := u.Path[1:]
		u.Path = "/"

		newDBInstance, err := sql.Open("postgres", u.String())
		if err != nil {
			message = err.Error()
			return
		}
		_, err = newDBInstance.ExecContext(context.Background(), `CREATE DATABASE "`+database+`"`)
		if err != nil {
			if !strings.Contains(err.Error(), "exists") {
				message = err.Error()
				_ = newDBInstance.Close()
				return
			}
		}

		_ = newDBInstance.Close()
	}

	// 测试创建数据表
	_, err = dbInstance.ExecContext(context.Background(), `CREATE TABLE "public"."teaweb_test" (
		"id" serial8 primary key,
		"_id" varchar(24)
		);
		
		CREATE UNIQUE INDEX "teaweb_test_id" ON "public"."teaweb_test" ("_id");`)
	if err != nil {
		message = "尝试创建数据表失败：" + err.Error()
		return
	}

	// 测试写入数据表
	_, err = dbInstance.ExecContext(context.Background(), "INSERT INTO \"teaweb_test\" (\"_id\") VALUES ('"+shared.NewObjectId().Hex()+"')")
	if err != nil {
		message = "尝试写入数据表失败：" + err.Error()
		return
	}

	// 测试删除数据表
	_, err = dbInstance.ExecContext(context.Background(), "DROP TABLE \"teaweb_test\"")
	if err != nil {
		message = "尝试删除数据表失败：" + err.Error()
		return
	}

	// 检查函数
	rows, err := dbInstance.Query(`SELECT JSON_EXTRACT_PATH('{"a":1}', 'a')`)
	if err != nil {
		message = "检查JSON_EXTRACT_PATH()函数失败：" + err.Error() + "。请尝试使用PostgreSQL v9.3以上版本。"
		return
	}
	_ = rows.Close()

	ok = true
	return
}

// 列出所有表
func (this *PostgresDriver) ListTables() ([]string, error) {
	currentDB, err := this.checkDB()
	if err != nil {
		return nil, err
	}

	rows, err := currentDB.Query("SELECT \"table_name\" FROM \"information_schema\".\"tables\"")
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
func (this *PostgresDriver) StatTables(tables []string) (map[string]*TableStat, error) {
	if len(tables) == 0 {
		return nil, nil
	}

	currentDB, err := this.checkDB()
	if err != nil {
		return nil, err
	}

	result := map[string]*TableStat{}
	for _, table := range tables {
		row := currentDB.QueryRow("SELECT COUNT(*) FROM " + this.quoteKeyword(table))
		countRows := int64(0)
		if row == nil {
			countRows = 0
		} else {
			err = row.Scan(&countRows)
			if err != nil {
				return nil, err
			}
		}

		tableSize := int64(0)
		formattedSize := ""

		row = currentDB.QueryRow("SELECT \"oid\" FROM \"pg_class\" WHERE \"relname\"=$1", table)
		if row == nil {
			tableSize = 0
		} else {
			oid := ""
			err = row.Scan(&oid)
			if err != nil {
				logs.Println("err:", err)
				return nil, err
			}

			row = currentDB.QueryRow("SELECT pg_total_relation_size($1)", oid)
			if row == nil {
				tableSize = 0
			} else {
				err = row.Scan(&tableSize)
				if err != nil {
					return nil, err
				}
			}

			if tableSize < 1024 {
				formattedSize = fmt.Sprintf("%.2fB", float64(tableSize))
			} else if tableSize < 1024*1024 {
				formattedSize = fmt.Sprintf("%.2fKB", float64(tableSize)/1024)
			} else if tableSize < 1024*1024*1024 {
				formattedSize = fmt.Sprintf("%.2fMB", float64(tableSize)/1024/1024)
			} else {
				formattedSize = fmt.Sprintf("%.2fGB", float64(tableSize)/1024/1024/1024)
			}
		}

		result[table] = &TableStat{
			Count:         countRows,
			Size:          tableSize,
			FormattedSize: formattedSize,
		}
	}

	return result, nil
}

// 清理访问日志任务
func (this *PostgresDriver) cleanAccessLogs() {
	reg := regexp.MustCompile(`^teaweb_logs_(\d{8})$`)

	teautils.Every(1*time.Minute, func(ticker *teautils.Ticker) {
		config, _ := db.LoadPostgresConfig()
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
		logs.Println("[postgres]clean access logs before '" + compareDay + "'")

		tables, err := this.ListTables()
		if err != nil {
			logs.Println("[postgres]" + err.Error())
			return
		}

		for _, table := range tables {
			if !reg.MatchString(table) {
				continue
			}

			if table < compareDay {
				logs.Println("[postgres]clean table '" + table + "'")
				err = this.DropTable(table)
				if err != nil {
					logs.Println("[postgres]" + err.Error())
				}
			}
		}
	})
}
