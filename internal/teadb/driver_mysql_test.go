package teadb

import (
	"database/sql"
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestMySQLDriver_Open(t *testing.T) {
	dbInstance, err := sql.Open("mysql", "root:abcdef@tcp(127.0.0.1:3306)/teaweb123?charset=utf8mb4&timeout=30s")
	if err != nil {
		t.Log("error:", err.Error())
		return
	}
	_ = dbInstance.Close()
	t.Log("ok")
}

func TestMySQLDriver_buildWhere(t *testing.T) {
	q := NewQuery("myTable")
	q.Attr("name", "lu")
	q.Attr("age", 10)
	q.Gt("timestamp", "gt")
	q.Lt("timestamp", "lt")
	q.Gte("timestamp", "gte")
	q.Lte("timestamp", "lte")
	q.Not("timestamp", "not")
	q.Attr("a", []string{"a", "b", "c"})
	q.Attr("timestamp", nil)
	q.Or([]*OperandList{
		NewOperandList().Add("timestamp", NewOperand(OperandEq, "123")),
		NewOperandList().Add("timestamp",
			NewOperand(OperandGt, "456"),
			NewOperand(OperandNotIn, []int{1, 2, 3}),
		),
	})

	driver := new(MySQLDriver)
	driver.Init()
	paramsHolder := NewSQLParamsHolder(driver.driver)
	where, err := driver.buildWhere(q.operandList, nil, paramsHolder)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("where:", where)
	logs.PrintAsJSON(paramsHolder.Params, t)
}

func TestMySQLDriver_buildWhere_Or(t *testing.T) {
	q := NewQuery("myTable")
	q.Or([]*OperandList{
		NewOperandList().Add("timestamp", NewOperand(OperandEq, "123")),
		NewOperandList().Add("timestamp",
			NewOperand(OperandGt, "456"),
			NewOperand(OperandNotIn, []int{1, 2, 3}),
		),
		NewOperandList().Add("timestamp", NewOperand(OperandLt, 1024)),
	})

	driver := new(MySQLDriver)
	driver.Init()
	paramsHolder := NewSQLParamsHolder(driver.driver)
	where, err := driver.buildWhere(q.operandList, nil, paramsHolder)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("where:", where)
	logs.PrintAsJSON(paramsHolder.Params, t)
}

func TestMySQLDriver_asSQL_SELECT(t *testing.T) {
	{
		q := NewQuery("myTable")
		q.Attr("name", "lu")
		q.Attr("age", 10)

		driver := new(MySQLDriver)
		driver.Init()
		s, err := driver.asSQL(SQLSelect, q, NewSQLParamsHolder(driver.driver), "", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(s)
	}

	{
		q := NewQuery("myTable")
		q.Result("a", "b")
		q.Attr("name", "lu")
		q.Attr("age", 10)
		q.Offset(20)
		q.Limit(10)
		q.Desc("_id")
		q.Asc("createdAt")

		driver := new(MySQLDriver)
		driver.driver = "mysql"
		s, err := driver.asSQL(SQLSelect, q, NewSQLParamsHolder(driver.driver), "", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(s)
	}
}

func TestMySQLDriver_asSQL_DELETE(t *testing.T) {
	{
		q := NewQuery("myTable")
		q.Attr("name", "lu")
		q.Attr("age", 10)

		driver := new(MySQLDriver)
		driver.driver = "mysql"
		s, err := driver.asSQL(SQLDelete, q, NewSQLParamsHolder(driver.driver), "", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(s)
	}

	{
		q := NewQuery("myTable")
		q.Result("a", "b")
		q.Attr("name", "lu")
		q.Attr("age", 10)
		q.Offset(20)
		q.Limit(10)
		q.Desc("_id")
		q.Asc("createdAt")

		driver := new(MySQLDriver)
		s, err := driver.asSQL(SQLDelete, q, NewSQLParamsHolder(driver.driver), "", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(s)
	}
}

func TestMySQLDriver_asSQL_Update(t *testing.T) {
	{
		q := NewQuery("myTable")
		q.Attr("name", "lu")
		q.Attr("age", 10)

		driver := new(MySQLDriver)
		s, err := driver.asSQL(SQLUpdate, q, NewSQLParamsHolder(driver.driver), "", map[string]interface{}{
			"name":   1,
			"age":    2,
			"count":  3,
			"book":   4,
			"number": 5,
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(s)
	}

	{
		q := NewQuery("myTable")
		q.Result("a", "b")
		q.Attr("name", "lu")
		q.Attr("age", 10)
		q.Offset(20)
		q.Limit(10)
		q.Desc("_id")
		q.Asc("createdAt")

		driver := new(MySQLDriver)
		s, err := driver.asSQL(SQLUpdate, q, NewSQLParamsHolder(driver.driver), "", map[string]interface{}{
			"a": 1,
			"b": 2,
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(s)
	}
}

func TestMySQLDriver_TestDSN(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	driver := new(MySQLDriver)
	{
		message, ok := driver.TestDSN("root:abcdef@tcp(127.0.0.1:3306)/teaweb123?charset=utf8mb4&timeout=30s", false)
		t.Log(message, ok)
	}
	{
		message, ok := driver.TestDSN("root:123456@tcp(127.0.0.1:3306)/teaweb123?charset=utf8mb4&timeout=30s", false)
		t.Log(message, ok)
	}

	{
		message, ok := driver.TestDSN("root:123456@tcp(127.0.0.1:3306)/teaweb?charset=utf8mb4&timeout=30s", false)
		t.Log(message, ok)
	}
}

func TestMySQLDriver_TestDSN_Create(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	driver := new(MySQLDriver)
	message, ok := driver.TestDSN("root:123456@tcp(127.0.0.1:3306)/teaweb?charset=utf8mb4&timeout=30s", true)
	t.Log(message, ok)
}

func TestMySQLDriver_Ping(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	driver := new(MySQLDriver)
	driver.driver = "mysql"
	err := driver.initDB()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(driver.Test())
}

func TestMySQLDriver_ListTables(t *testing.T) {
	if !teatesting.RequireDBAvailable() || !teatesting.RequireMySQL() {
		return
	}

	driver := new(MySQLDriver)
	driver.isAvailable = true
	err := driver.initDB()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(driver.ListTables())
}

func TestMySQLDriver_StatTables(t *testing.T) {
	if !teatesting.RequireDBAvailable() || !teatesting.RequireMySQL() {
		return
	}

	driver := new(MySQLDriver)
	driver.isAvailable = true
	err := driver.initDB()
	if err != nil {
		t.Fatal(err)
	}

	tables, err := driver.ListTables()
	if err != nil {
		t.Fatal(err)
	}
	result, err := driver.StatTables(tables)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(result, t)
}
