package tealogs

import (
	"database/sql"
	"errors"
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iwind/TeaGo/logs"
	"strconv"
	"strings"
)

// MySQL存储
type MySQLStorage struct {
	Storage `yaml:", inline"`

	Host            string `yaml:"host" json:"host"`
	Port            int    `yaml:"port" json:"port"`
	Username        string `yaml:"username" json:"username"`
	Password        string `yaml:"password" json:"password"`
	Database        string `yaml:"database" json:"database"`
	Table           string `yaml:"table" json:"table"` // can use variables, such as ${year}
	LogField        string `yaml:"logField" json:"logField"`
	AutoCreateTable bool   `yaml:"autoCreateTable" json:"autoCreateTable"`
}

// 启动
func (this *MySQLStorage) Start() error {
	if len(this.Host) == 0 {
		return errors.New("'host' should not be empty")
	}
	if this.Port <= 0 {
		this.Port = 3306
	}
	if len(this.Database) == 0 {
		return errors.New("'database' should not be empty")
	}
	if len(this.Table) == 0 {
		return errors.New("'table' should not be empty")
	}
	if len(this.LogField) == 0 {
		return errors.New("'logField' should not be empty")
	}

	return nil
}

// 写入日志
func (this *MySQLStorage) Write(accessLogs []*accesslogs.AccessLog) error {
	if len(accessLogs) == 0 {
		return nil
	}

	dsn := ""
	if len(this.Username) > 0 {
		dsn += this.Username + ":" + this.Password + "@"
	}
	dsn += "tcp(" + this.Host + ":" + strconv.Itoa(this.Port) + ")/" + this.Database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	table := this.FormatVariables(this.Table)

	for _, accessLog := range accessLogs {
		dataString, err := this.FormatAccessLogString(accessLog)
		if err != nil {
			logs.Error(err)
			continue
		}
		_, err = db.Exec("INSERT INTO `"+table+"` (`"+this.LogField+"`) VALUES (?)", dataString)
		if err != nil {
			if strings.Contains(err.Error(), "Error 1146") && this.AutoCreateTable {
				// try to create table
				_, err = db.Exec("CREATE TABLE `" + table + "` ( " +
					"`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT, " +
					"`" + this.LogField + "` longtext COLLATE utf8mb4_bin, " +
					"PRIMARY KEY (`id`) " +
					") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;")
				if err != nil && !strings.Contains(err.Error(), "Error 1050") {
					logs.Error(err)
				} else {
					// 再次尝试
					_, err = db.Exec("INSERT INTO `"+table+"` (`"+this.LogField+"`) VALUES (?)", dataString)
					if err != nil {
						logs.Error(err)
					}
				}
			} else {
				logs.Error(err)
			}
		}
	}

	return nil
}

// 关闭
func (this *MySQLStorage) Close() error {
	return nil
}
