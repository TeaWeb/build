package teadb

import (
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"testing"
)

func TestMongoDriver_buildFilter(t *testing.T) {
	q := new(Query)
	q.Init()
	q.Attr("name", "lu")
	q.Op("age", OperandGt, 1024)
	q.Op("age", OperandLt, 2048)
	q.Op("count", OperandEq, 3)

	driver := new(MongoDriver)
	filter, err := driver.buildFilter(q)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(filter, t)
}

func TestMongoDriver_buildFilter_Or(t *testing.T) {
	q := new(Query)
	q.Init()
	q.Attr("a", 1)
	q.Or([]*OperandList{
		NewOperandList().Add("timestamp", NewOperand(OperandEq, "123")),
		NewOperandList().Add("timestamp",
			NewOperand(OperandGt, "456"),
			NewOperand(OperandNotIn, []int{1, 2, 3}),
		),
		NewOperandList().Add("timestamp", NewOperand(OperandLt, 1024)),
	})
	driver := new(MongoDriver)
	filter, err := driver.buildFilter(q)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(filter, t)
}

func TestMongoDriver_setMapValue(t *testing.T) {
	m := maps.Map{}

	driver := new(MongoDriver)
	driver.setMapValue(m, []string{"a", "b", "c", "d", "e"}, 123)
	logs.PrintAsJSON(m, t)
}

func TestMongoDriver_connect(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	driver := new(MongoDriver)
	client, err := driver.connect()
	if err != nil {
		t.Log("ERROR:", err.Error())
		return
	}
	t.Log("client:", client)
}

func TestMongoDriver_Test(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	driver := new(MongoDriver)
	err := driver.Test()
	if err != nil {
		t.Log("ERROR:", err.Error())
		return
	}
	t.Log("client:", driver)
}

func TestMongoDriver_convertArrayElement(t *testing.T) {
	driver := new(MongoDriver)
	t.Log(driver.convertArrayElement("value.usage.avg"))
	t.Log(driver.convertArrayElement("value.usage.all.0"))
	t.Log(driver.convertArrayElement("value.0"))
}

func TestMongoDriver_ListTables(t *testing.T) {
	if !teatesting.RequireDBAvailable() || !teatesting.RequireMongoDB() {
		return
	}

	driver := new(MongoDriver)
	driver.isAvailable = true
	names, err := driver.ListTables()
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(names)
}
