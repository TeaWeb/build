package teadb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs/db"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var jsonArrayIndexReg = regexp.MustCompile(`\.(\d+)`)

type SQLDriver struct {
	BaseDriver

	driver string

	db       *sql.DB
	dbLocker sync.Mutex

	stmtMap    map[string]*sql.Stmt // query => stmt
	stmtLocker sync.Mutex

	sqlMode string
}

// 查找单条记录
func (this *SQLDriver) FindOne(query *Query, modelPtr interface{}) (interface{}, error) {
	ones, err := this.FindOnes(query.Limit(1), modelPtr)
	if err != nil {
		return nil, err
	}
	if len(ones) == 0 {
		return nil, nil
	}
	return ones[0], nil
}

// 查找多条记录
func (this *SQLDriver) FindOnes(query *Query, modelPtr interface{}) ([]interface{}, error) {
	currentDB, err := this.checkDB()
	if err != nil {
		return nil, err
	}

	holder := NewSQLParamsHolder(this.driver)
	sqlString, err := this.asSQL(SQLSelect, query, holder, "", nil)
	if err != nil {
		return nil, err
	}

	if query.debug {
		logs.Println("sql:", sqlString)
	}

	stmt, ok := this.findStmt(sqlString)
	if !ok {
		stmt, err = currentDB.PrepareContext(context.Background(), sqlString)
		if err != nil {
			return nil, this.processError(err)
		}
		this.putStmt(sqlString, stmt)
	}

	rows, err := stmt.Query(holder.Args...)
	if err != nil {
		return nil, this.processError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	modelType := reflect.TypeOf(modelPtr)
	modelElem := modelType.Elem()
	method, methodExists := modelType.MethodByName("SetDBColumns")
	result := []interface{}{}
	for rows.Next() {
		cols, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		colPtrList := []interface{}{}
		for range cols {
			p := interface{}(nil)
			colPtrList = append(colPtrList, &p)
		}
		err = rows.Scan(colPtrList...)
		if err != nil {
			return nil, err
		}
		values := maps.Map{}
		for index, col := range cols {
			v := reflect.Indirect(reflect.ValueOf(colPtrList[index])).Interface()
			if v != nil {
				if _, ok := v.([]byte); ok {
					v = string(v.([]byte))
				}
			}
			values[col] = v
		}
		one := reflect.New(modelElem)
		if methodExists {
			method.Func.Call([]reflect.Value{one, reflect.ValueOf(values)})
		}
		result = append(result, one.Interface())
	}

	return result, nil
}

// 插入一条记录
func (this *SQLDriver) InsertOne(table string, modelPtr interface{}) error {
	table = strings.Replace(table, ".", "_", -1)

	currentDB, err := this.checkDB()
	if err != nil {
		return err
	}

	if modelPtr == nil {
		return errors.New("modelPtr should not be nil")
	}

	modelType := reflect.TypeOf(modelPtr)
	method, methodExists := modelType.MethodByName("DBColumns")
	if !methodExists {
		return errors.New("'DBColumns() maps.Map' method not exist in '" + modelType.String() + "'")
	}
	result := method.Func.Call([]reflect.Value{reflect.ValueOf(modelPtr)})
	if len(result) != 1 {
		return errors.New("'DBColumns() maps.Map' method not exist in '" + modelType.String() + "'")
	}
	v := result[0].Interface()
	m, ok := v.(maps.Map)
	if !ok {
		return errors.New("'DBColumns() maps.Map' method not exist in '" + modelType.String() + "'")
	}

	// 对字段进行排序
	keys := teautils.MapKeys(m)
	sort.Strings(keys)

	b := strings.Builder{}
	b.WriteString("INSERT INTO " + this.quoteKeyword(table) + " (")
	index := 0
	args := []interface{}{}
	for _, k := range keys {
		if index > 0 {
			b.WriteString(", ")
		}
		b.WriteString(this.quoteKeyword(k))
		args = append(args, m.Get(k))
		index++
	}
	b.WriteString(") ")
	b.WriteString("VALUES (")
	for index := range args {
		if index > 0 {
			switch this.driver {
			case "mysql":
				b.WriteString(", ?")
			case "postgres":
				b.WriteString(", $" + strconv.Itoa(index+1))
			default:
				b.WriteString(", ?")
			}
		} else {
			switch this.driver {
			case "mysql":
				b.WriteString("?")
			case "postgres":
				b.WriteString("$" + strconv.Itoa(index+1))
			default:
				b.WriteString("?")
			}
		}
	}
	b.WriteString(")")
	sqlString := b.String()
	stmt, ok := this.findStmt(b.String())
	if !ok {
		stmt, err = currentDB.PrepareContext(context.Background(), sqlString)
		if err != nil {
			return this.processError(err)
		}
		this.putStmt(sqlString, stmt)
	}

	_, err = stmt.Exec(args...)

	return this.processError(err)
}

// 插入多条记录
func (this *SQLDriver) InsertOnes(table string, modelPtrSlice interface{}) error {
	currentDB, err := this.checkDB()
	if err != nil {
		return err
	}
	if modelPtrSlice == nil {
		return nil
	}

	sliceType := reflect.TypeOf(modelPtrSlice)
	if sliceType.Kind() != reflect.Slice {
		return errors.New("only slice can be accepted in 'InsertOnes' method")
	}

	modelValues := reflect.ValueOf(modelPtrSlice)
	countValues := modelValues.Len()
	if modelValues.Len() == 0 {
		return nil
	}

	modelPtr := modelValues.Index(0).Interface()
	modelType := reflect.TypeOf(modelPtr)
	method, methodExists := modelType.MethodByName("DBColumns")
	if !methodExists {
		return errors.New("'DBColumns() maps.Map' method not exist in '" + modelType.String() + "'")
	}

	b := strings.Builder{}
	result := method.Func.Call([]reflect.Value{reflect.ValueOf(modelPtr)})
	if len(result) != 1 {
		return errors.New("'DBColumns() maps.Map' method not exist in '" + modelType.String() + "'")
	}
	v := result[0].Interface()
	m, ok := v.(maps.Map)
	if !ok {
		return errors.New("'DBColumns() maps.Map' method not exist in '" + modelType.String() + "'")
	}
	b.WriteString("INSERT INTO " + this.quoteKeyword(table) + " (")
	keys := []string{}
	index := 0
	for k := range m {
		if index > 0 {
			b.WriteString(", ")
		}
		b.WriteString(this.quoteKeyword(k))
		keys = append(keys, k)
		index++
	}
	b.WriteString(") ")
	b.WriteString("VALUES ")

	args := []interface{}{}
	paramIndex := 0
	for i := 0; i < countValues; i++ {
		if i > 0 {
			b.WriteString(", ")
		}
		modelValue := modelValues.Index(i)
		result := method.Func.Call([]reflect.Value{reflect.ValueOf(modelValue.Interface())})
		if len(result) != 1 {
			return errors.New("'DBColumns() maps.Map' method not exist in '" + modelType.String() + "'")
		}
		v := result[0].Interface()
		m, ok := v.(maps.Map)
		if !ok {
			return errors.New("'DBColumns() maps.Map' method not exist in '" + modelType.String() + "'")
		}

		b.WriteString("(")
		for index, key := range keys {
			if index > 0 {
				switch this.driver {
				case "mysql":
					b.WriteString(", ?")
				case "postgres":
					b.WriteString(", $" + strconv.Itoa(paramIndex+1))
					paramIndex++
				default:
					b.WriteString(", ?")
				}
			} else {
				switch this.driver {
				case "mysql":
					b.WriteString("?")
				case "postgres":
					b.WriteString("$" + strconv.Itoa(paramIndex+1))
					paramIndex++
				default:
					b.WriteString("?")
				}
			}
			args = append(args, m.Get(key))
		}
		b.WriteString(")")
	}

	stmt, err := currentDB.PrepareContext(context.Background(), b.String())
	if err != nil {
		return this.processError(err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec(args...)

	return this.processError(err)
}

// 删除多条记录
func (this *SQLDriver) DeleteOnes(query *Query) error {
	currentDB, err := this.checkDB()
	if err != nil {
		return err
	}

	holder := NewSQLParamsHolder(this.driver)
	sqlString, err := this.asSQL(SQLDelete, query, holder, "", nil)
	if err != nil {
		return err
	}

	stmt, ok := this.findStmt(sqlString)
	if !ok {
		stmt, err = currentDB.PrepareContext(context.Background(), sqlString)
		if err != nil {
			return this.processError(err)
		}
		this.putStmt(sqlString, stmt)
	}

	_, err = stmt.Exec(holder.Args...)
	return this.processError(err)
}

// 修改多条记录
func (this *SQLDriver) UpdateOnes(query *Query, values map[string]interface{}) error {
	currentDB, err := this.checkDB()
	if err != nil {
		return err
	}

	holder := NewSQLParamsHolder(this.driver)
	sqlString, err := this.asSQL(SQLUpdate, query, holder, "", values)
	if err != nil {
		return err
	}

	stmt, ok := this.findStmt(sqlString)
	if !ok {
		stmt, err = currentDB.PrepareContext(context.Background(), sqlString)
		if err != nil {
			return this.processError(err)
		}
		this.putStmt(sqlString, stmt)
	}

	_, err = stmt.Exec(holder.Args...)

	return this.processError(err)
}

// 计算总数量
func (this *SQLDriver) Count(query *Query) (int64, error) {
	currentDB, err := this.checkDB()
	if err != nil {
		return 0, err
	}

	holder := NewSQLParamsHolder(this.driver)
	query.Result("COUNT(*)")
	sqlString, err := this.asSQL(SQLSelect, query, holder, "", nil)

	if query.debug {
		logs.Println("sql:", sqlString)
	}

	if err != nil {
		return 0, err
	}
	stmt, ok := this.findStmt(sqlString)
	if !ok {
		stmt, err = currentDB.PrepareContext(context.Background(), sqlString)
		if err != nil {
			return 0, this.processError(err)
		}
		this.putStmt(sqlString, stmt)
	}
	row := stmt.QueryRow(holder.Args...)
	if row == nil {
		return 0, nil
	}

	result := interface{}(nil)
	err = row.Scan(&result)
	if err != nil {
		return 0, this.processError(err)
	}

	return types.Int64(result), nil
}

// 计算总和
func (this *SQLDriver) Sum(query *Query, field string) (float64, error) {
	currentDB, err := this.checkDB()
	if err != nil {
		return 0, err
	}

	holder := NewSQLParamsHolder(this.driver)
	query.Result("SUM(" + this.quoteKeyword(field) + ")")
	sqlString, err := this.asSQL(SQLSelect, query, holder, "", nil)

	if query.debug {
		logs.Println("sql:", sqlString)
	}

	if err != nil {
		return 0, err
	}
	stmt, ok := this.findStmt(sqlString)
	if !ok {
		stmt, err = currentDB.PrepareContext(context.Background(), sqlString)
		if err != nil {
			return 0, this.processError(err)
		}
		this.putStmt(sqlString, stmt)
	}
	row := stmt.QueryRow(holder.Args...)
	if row == nil {
		return 0, nil
	}

	result := interface{}(nil)
	err = row.Scan(&result)
	if err != nil {
		return 0, this.processError(err)
	}

	return types.Float64(result), nil
}

// 计算平均值
func (this *SQLDriver) Avg(query *Query, field string) (float64, error) {
	currentDB, err := this.checkDB()
	if err != nil {
		return 0, err
	}

	holder := NewSQLParamsHolder(this.driver)
	query.Result("AVG(" + this.quoteKeyword(field) + ")")
	sqlString, err := this.asSQL(SQLSelect, query, holder, "", nil)

	if query.debug {
		logs.Println("sql:", sqlString)
	}

	if err != nil {
		return 0, err
	}
	stmt, ok := this.findStmt(sqlString)
	if !ok {
		stmt, err = currentDB.PrepareContext(context.Background(), sqlString)
		if err != nil {
			return 0, this.processError(err)
		}
		this.putStmt(sqlString, stmt)
	}

	row := stmt.QueryRow(holder.Args...)
	if row == nil {
		return 0, nil
	}

	result := interface{}(nil)
	err = row.Scan(&result)
	if err != nil {
		return 0, this.processError(err)
	}

	return types.Float64(result), nil
}

// 计算最小值
func (this *SQLDriver) Min(query *Query, field string) (float64, error) {
	currentDB, err := this.checkDB()
	if err != nil {
		return 0, err
	}

	holder := NewSQLParamsHolder(this.driver)
	query.Result("MIN(" + this.quoteKeyword(field) + ")")
	sqlString, err := this.asSQL(SQLSelect, query, holder, "", nil)

	if query.debug {
		logs.Println("sql:", sqlString)
	}

	if err != nil {
		return 0, err
	}

	stmt, ok := this.findStmt(sqlString)
	if !ok {
		stmt, err = currentDB.PrepareContext(context.Background(), sqlString)
		if err != nil {
			return 0, this.processError(err)
		}
		this.putStmt(sqlString, stmt)
	}

	row := stmt.QueryRow(holder.Args...)
	if row == nil {
		return 0, nil
	}

	result := interface{}(nil)
	err = row.Scan(&result)
	if err != nil {
		return 0, this.processError(err)
	}

	return types.Float64(result), nil
}

// 计算最大值
func (this *SQLDriver) Max(query *Query, field string) (float64, error) {
	currentDB, err := this.checkDB()
	if err != nil {
		return 0, err
	}

	holder := NewSQLParamsHolder(this.driver)
	query.Result("MAX(" + this.quoteKeyword(field) + ")")
	sqlString, err := this.asSQL(SQLSelect, query, holder, "", nil)

	if query.debug {
		logs.Println("sql:", sqlString)
	}

	if err != nil {
		return 0, err
	}

	stmt, ok := this.findStmt(sqlString)
	if !ok {
		stmt, err = currentDB.PrepareContext(context.Background(), sqlString)
		if err != nil {
			return 0, this.processError(err)
		}
		this.putStmt(sqlString, stmt)
	}

	row := stmt.QueryRow(holder.Args...)
	if row == nil {
		return 0, nil
	}

	result := interface{}(nil)
	err = row.Scan(&result)
	if err != nil {
		return 0, this.processError(err)
	}

	return types.Float64(result), nil
}

// 对数据进行分组统计
func (this *SQLDriver) Group(query *Query, groupField string, result map[string]Expr) ([]maps.Map, error) {
	currentDB, err := this.checkDB()
	if err != nil {
		return nil, err
	}

	for field, expr := range result {
		switch e := expr.(type) {
		case *SumExpr:
			index := strings.Index(e.Field, ".")
			if index > -1 {
				e.Field = this.JSONExtractNumeric(e.Field[:index], e.Field[index+1:])
			}
			query.Result("SUM(" + this.quoteKeyword(e.Field) + ") AS " + this.quoteKeyword(field))
		case *AvgExpr:
			index := strings.Index(e.Field, ".")
			if index > -1 {
				e.Field = this.JSONExtractNumeric(e.Field[:index], e.Field[index+1:])
			}
			query.Result("AVG(" + this.quoteKeyword(e.Field) + ") AS " + this.quoteKeyword(field))
		case *MaxExpr:
			index := strings.Index(e.Field, ".")
			if index > -1 {
				e.Field = this.JSONExtractNumeric(e.Field[:index], e.Field[index+1:])
			}
			query.Result("MAX(" + this.quoteKeyword(e.Field) + ") AS " + this.quoteKeyword(field))
		case *MinExpr:
			index := strings.Index(e.Field, ".")
			if index > -1 {
				e.Field = this.JSONExtractNumeric(e.Field[:index], e.Field[index+1:])
			}
			query.Result("MIN(" + this.quoteKeyword(e.Field) + ") AS " + this.quoteKeyword(field))
		case string:
			index := strings.Index(e, ".")
			isAgg := false
			if index > -1 {
				if this.driver == "postgres" {
					isAgg = true
					e = this.JSONExtract("(array_agg("+this.quoteKeyword(e[:index])+"))[1]", e[index+1:])
				} else {
					e = this.JSONExtract(e[:index], e[index+1:])
				}
			}
			if this.driver == "postgres" && !isAgg {
				query.Result("(array_agg(" + this.quoteKeyword(e) + "))[1] AS " + this.quoteKeyword(field))
			} else {
				query.Result(this.quoteKeyword(e) + " AS " + this.quoteKeyword(field))
			}
		}
	}

	holder := NewSQLParamsHolder(this.driver)
	sqlString, err := this.asSQL(SQLSelect, query, holder, groupField, nil)
	if err != nil {
		return nil, err
	}

	if query.debug {
		logs.Println("sql:", sqlString)
	}

	var stmt *sql.Stmt = nil
	if this.driver == db.DBTypeMySQL && strings.Contains(this.sqlMode, "ONLY_FULL_GROUP_BY") {
		tx, err := currentDB.Begin()
		if err != nil {
			return nil, err
		}
		defer func() {
			_ = tx.Commit()
		}()

		// 屏蔽MySQL的ONLY_FULL_GROUP_BY选项
		_, err = tx.ExecContext(context.Background(), "SET SESSION sql_mode=(SELECT REPLACE(@@sql_mode,'ONLY_FULL_GROUP_BY',''));")
		if err != nil {
			logs.Error(err)
		}

		// 事务不能缓存SQL
		stmt, err = tx.PrepareContext(context.Background(), sqlString)
		if err != nil {
			return nil, this.processError(err)
		}
	} else {
		var ok = false
		stmt, ok = this.findStmt(sqlString)
		if !ok {
			stmt, err = currentDB.PrepareContext(context.Background(), sqlString)
			if err != nil {
				return nil, this.processError(err)
			}
			this.putStmt(sqlString, stmt)
		}
	}

	rows, err := stmt.Query(holder.Args...)
	if err != nil {
		return nil, this.processError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	ones := []maps.Map{}
	for rows.Next() {
		columns, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		result := []interface{}{}
		for range columns {
			v := interface{}(nil)
			result = append(result, &v)
		}
		err = rows.Scan(result...)
		if err != nil {
			return nil, this.processError(err)
		}
		m := maps.Map{}
		for index, column := range columns {
			v := reflect.Indirect(reflect.ValueOf(result[index])).Interface()
			if v != nil {
				switch v1 := v.(type) {
				case []byte:
					v = string(v1)
				}
			}

			keys := strings.Split(column, ".")
			if len(keys) > 1 {
				this.setMapValue(m, keys, v)
			} else {
				m[column] = v
			}
		}

		ones = append(ones, m)
	}

	return ones, nil
}

// 测试数据库连接
func (this *SQLDriver) Test() error {
	_, err := this.checkDB()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	done := make(chan bool, 1)
	isDone := false

	go func() {
		select {
		case <-ctx.Done():
			cancel()
			err = errors.New("unable connect to database: timeout")

			if !isDone {
				isDone = true
				done <- true
			}
		}
	}()

	go func() {
		err = this.db.PingContext(ctx)
		if !isDone {
			isDone = true
			done <- true
		}
	}()

	<-done
	isDone = true
	close(done)

	return err
}

// 重启服务
func (this *SQLDriver) Shutdown() error {
	this.dbLocker.Lock()
	defer this.dbLocker.Unlock()
	if this.db != nil {
		_ = this.db.Close()
	}
	return nil
}

// 删除表
func (this *SQLDriver) DropTable(table string) error {
	currentDB, err := this.checkDB()
	if err != nil {
		return err
	}
	_, err = currentDB.ExecContext(context.Background(), "DROP TABLE "+this.quoteKeyword(table))
	return this.processError(err)
}

// 读取JSON字段
func (this *SQLDriver) JSONExtract(field string, path string) string {
	switch this.driver {
	case "mysql":
		return "JSON_EXTRACT(" + this.quoteKeyword(field) + ", \"$." + jsonArrayIndexReg.ReplaceAllString(path, "[$1]") + "\")"
	case "postgres":
		return "JSON_EXTRACT_PATH_TEXT(" + this.quoteKeyword(field) + ", '" + strings.Replace(path, ".", "', '", -1) + "')"
	}
	return ""
}

// 读取JSON字段
func (this *SQLDriver) JSONExtractNumeric(field string, path string) string {
	switch this.driver {
	case "mysql":
		return "JSON_EXTRACT(" + this.quoteKeyword(field) + ", \"$." + jsonArrayIndexReg.ReplaceAllString(path, "[$1]") + "\")"
	case "postgres":
		return "JSON_EXTRACT_PATH_TEXT(" + this.quoteKeyword(field) + ", '" + strings.Replace(path, ".", "', '", -1) + "')::\"float8\""
	}
	return ""
}

// 处理错误
func (this *SQLDriver) processError(err error) error {
	if err == nil {
		return nil
	}
	if err == sql.ErrConnDone || err == driver.ErrBadConn {
		// 断开连接时处理
	}
	return err
}

func (this *SQLDriver) setMapValue(m maps.Map, keys []string, v interface{}) {
	l := len(keys)
	if l == 0 {
		return
	}
	if l == 1 {
		m[keys[0]] = v
		return
	}
	subM, ok := m[keys[0]]
	if ok {
		subV, ok := subM.(maps.Map)
		if ok {
			this.setMapValue(subV, keys[1:], v)
		} else {
			m[keys[0]] = maps.Map{}
			this.setMapValue(m[keys[0]].(maps.Map), keys[1:], v)
		}
	} else {
		m[keys[0]] = maps.Map{}
		this.setMapValue(m[keys[0]].(maps.Map), keys[1:], v)
	}
}

func (this *SQLDriver) asSQL(action SQLAction, query *Query, paramsHolder *SQLParamsHolder, groupField string, updateValues map[string]interface{}) (string, error) {
	b := strings.Builder{}

	switch action {
	case SQLSelect:
		b.WriteString("SELECT ")

		// result
		if len(query.resultFields) == 0 {
			b.WriteString("* ")
		} else {
			for index, field := range query.resultFields {
				if index > 0 {
					b.WriteString(", ")
				}
				b.WriteString(this.quoteKeyword(field))
			}
			b.WriteString(" ")
		}
	case SQLDelete:
		b.WriteString("DELETE ")
	case SQLUpdate:
		b.WriteString("UPDATE ")
	}

	// table
	if action == SQLSelect || action == SQLDelete {
		b.WriteString("FROM ")
	}
	b.WriteString(this.quoteKeyword(query.table))
	b.WriteString(" ")

	// set
	if action == SQLUpdate {
		b.WriteString("SET ")
		index := 0

		keys := teautils.MapKeys(updateValues)
		sort.Strings(keys)

		for _, k := range keys {
			if index > 0 {
				b.WriteString(", ")
			}
			b.WriteString(this.quoteKeyword(k))
			b.WriteString("=")
			b.WriteString(paramsHolder.Add(updateValues[k]))
			index++
		}
		b.WriteString(" ")
	}

	// where
	if query.operandList.Len() > 0 {
		where, err := this.buildWhere(query.operandList, query.fieldMapping, paramsHolder)
		if err != nil {
			return "", err
		}
		if len(where) > 0 {
			b.WriteString("WHERE ")
			b.WriteString(where)
			b.WriteString(" ")
		}
	}

	// group
	hasGroups := false
	if action == SQLSelect && len(groupField) > 0 {
		if query.fieldMapping != nil {
			groupField = query.fieldMapping(groupField)
		}
		b.WriteString("GROUP BY " + this.quoteKeyword(groupField))
		b.WriteString(" ")
		hasGroups = true
	}

	// order
	if action == SQLSelect && len(query.sortFields) > 0 {
		b.WriteString("ORDER BY ")
		for index, field := range query.sortFields {
			if index > 0 {
				b.WriteString(", ")
			}
			if query.fieldMapping != nil {
				field.Name = query.fieldMapping(field.Name)
			}

			// 支持点符号
			if strings.IndexAny(field.Name, "( ") == -1 {
				dotIndex := strings.Index(field.Name, ".")
				if dotIndex > -1 {
					field.Name = this.JSONExtract(field.Name[:dotIndex], field.Name[dotIndex+1:])
					if this.driver == db.DBTypePostgres {
						field.Name += "::\"float8\""
					}
				}
			}
			if hasGroups && this.driver == db.DBTypePostgres && !strings.ContainsAny(field.Name, "({:") {
				b.WriteString("(array_agg(" + this.quoteKeyword(field.Name) + "))[1]")
			} else {
				b.WriteString(this.quoteKeyword(field.Name))
			}
			if field.IsAsc() {
				b.WriteString(" ASC ")
			} else {
				b.WriteString(" DESC ")
			}
		}
	}

	// limit
	if query.size > 0 {
		b.WriteString("LIMIT " + strconv.Itoa(query.size) + " ")
	}
	if query.offset > 0 {
		b.WriteString("OFFSET " + strconv.Itoa(query.offset) + " ")
	}

	if len(paramsHolder.Params) > 0 {
		return paramsHolder.Parse(b.String()), nil
	}

	return b.String(), nil
}

func (this *SQLDriver) quoteKeyword(s string) string {
	if strings.IndexAny(s, "( :{") > -1 {
		return s
	}
	switch this.driver {
	case "mysql":
		return "`" + s + "`"
	case "postgres":
		return "\"" + s + "\""
	}
	return "\"" + s + "\""
}

// 构造where
func (this *SQLDriver) buildWhere(operandList *OperandList, fieldMapping func(field string) string, paramsHolder *SQLParamsHolder) (string, error) {
	b := strings.Builder{}
	hasPrefix := false

	var resultErr error = nil

	operandList.Range(func(field string, operands []*Operand) {
		if fieldMapping != nil {
			field = fieldMapping(field)
		}
		for _, op := range operands {
			if !hasPrefix {
				hasPrefix = true
			} else {
				b.WriteString(" AND ")
			}
			switch op.Code {
			case OperandEq:
				b.WriteString(this.quoteKeyword(field) + "=" + paramsHolder.Add(op.Value))
			case OperandLt:
				b.WriteString(this.quoteKeyword(field) + "<" + paramsHolder.Add(op.Value))
			case OperandLte:
				b.WriteString(this.quoteKeyword(field) + "<=" + paramsHolder.Add(op.Value))
			case OperandGt:
				b.WriteString(this.quoteKeyword(field) + ">" + paramsHolder.Add(op.Value))
			case OperandGte:
				b.WriteString(this.quoteKeyword(field) + ">=" + paramsHolder.Add(op.Value))
			case OperandIn:
				b.WriteString(this.quoteKeyword(field) + " IN " + paramsHolder.AddSlice(op.Value))
			case OperandNotIn:
				b.WriteString(this.quoteKeyword(field) + " NOT IN " + paramsHolder.AddSlice(op.Value))
			case OperandNeq:
				b.WriteString(this.quoteKeyword(field) + "!=" + paramsHolder.AddSlice(op.Value))
			case operandSQLCond:
				if op.Value != nil {
					cond, ok := op.Value.(*SQLCond)
					if ok {
						b.WriteString(cond.Expr)
						for k, v := range cond.Params {
							paramsHolder.AddHolder(k, v)
						}
					} else {
						resultErr = errors.New("operand 'operandSQLCond' value must be '*SQLCond'")
						return
					}
				}
			case OperandOr:
				if op.Value != nil {
					operandLists, ok := op.Value.([]*OperandList)
					if ok {
						if len(operandLists) > 1 {
							b.WriteString("(")
						}
						for index, operandList := range operandLists {
							f, err := this.buildWhere(operandList, fieldMapping, paramsHolder)
							if err != nil {
								resultErr = err
								return
							}
							if index > 0 {
								b.WriteString("OR ")
							}
							b.WriteString("(")
							b.WriteString(f)
							b.WriteString(") ")
						}
						if operandList.Len() > 1 {
							b.WriteString(") ")
						}
					} else {
						resultErr = errors.New("or: should be a valid []OperandMap")
						return
					}
				} else {
					resultErr = errors.New("or: should be a valid []OperandMap")
					return
				}
			default:
				resultErr = errors.New("invalid operand '" + op.Code + "'")
				return
			}
		}
	})
	if resultErr != nil {
		return "", resultErr
	}
	return b.String(), nil
}

func (this *SQLDriver) checkDB() (*sql.DB, error) {
	if this.db == nil {
		return nil, errors.New("db open failed")
	}
	return this.db, nil
}

func (this *SQLDriver) findStmt(query string) (stmt *sql.Stmt, ok bool) {
	this.stmtLocker.Lock()
	stmt, ok = this.stmtMap[query]
	this.stmtLocker.Unlock()

	return stmt, ok
}

func (this *SQLDriver) putStmt(query string, stmt *sql.Stmt) {
	this.stmtLocker.Lock()

	if this.stmtMap == nil {
		this.stmtMap = map[string]*sql.Stmt{}
	}

	// 限制最多只能缓存1024个SQL
	if len(this.stmtMap) >= 1024 {
		for _, s := range this.stmtMap {
			_ = s.Close()
		}
		this.stmtMap = map[string]*sql.Stmt{}
	}
	this.stmtMap[query] = stmt
	this.stmtLocker.Unlock()
}
