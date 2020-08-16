package teadb

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teadb/shared"
	"github.com/iwind/TeaGo/logs"
	"testing"
	"time"
)

func TestQuery_Asc(t *testing.T) {
	q := new(Query)
	q.Init()

	q.Asc("name")
	q.Desc("age")

	logs.PrintAsJSON(q.sortFields, t)

	q.Desc("name")
	logs.PrintAsJSON(q.sortFields, t)
}

func TestQuery_FindOne(t *testing.T) {
	q := new(Query)
	q.Init()
	q.Table("values.agent.local")
	q.Attr("appId", "system")
	q.Attr("itemId", "cpu.load")
	//q.Op("timestamp", OperandGt, 1553586946)
	q.Desc("createdAt")

	one, err := q.FindOne(new(agents.Value))
	if err != nil {
		t.Fatal(err)
	}
	if one == nil {
		t.Log("not found")
		return
	}
	t.Log(one.(*agents.Value))
	t.Log(one.(*agents.Value).Timestamp > 1553586946)
	logs.PrintAsJSON(one, t)
}

func TestQuery_FindOnes(t *testing.T) {
	q := new(Query)
	q.Init()

	q.Table("values.agent.local")
	q.Attr("appId", "system")
	q.Attr("itemId", "cpu.load")
	q.Desc("createdAt")

	q.Limit(5)
	ones, err := q.FindOnes(new(agents.Value))
	if err != nil {
		t.Fatal(err)
	}
	if len(ones) == 0 {
		t.Log("not found")
		return
	}
	for _, one := range ones {
		v := one.(*agents.Value)
		t.Log(v, one)
	}
}

func TestQuery_InsertOne(t *testing.T) {
	q := new(Query)
	q.Init()

	q.Table("abc")

	v := new(agents.Value)
	v.Id = shared.NewObjectId()
	v.Value = map[string]interface{}{
		"load1":  1,
		"load5":  5,
		"load15": 15,
	}

	v.SetTime(time.Now())

	err := q.InsertOne(v)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestQuery_InsertOnes(t *testing.T) {
	q := new(Query)
	q.Init()

	q.Table("abc")

	s := []interface{}{}
	for i := 0; i < 10; i++ {
		v := new(agents.Value)
		v.Id = shared.NewObjectId()
		v.Value = map[string]interface{}{
			"load1":  i,
			"load5":  5,
			"load15": 15,
		}

		v.SetTime(time.Now())
		s = append(s, v)
	}

	err := q.InsertOnes(s)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestQuery_Count(t *testing.T) {
	q := new(Query)
	q.Init()
	q.Attr("appId", "system")
	q.Table("values.agent.local")

	count, err := q.Count()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("count:", count)
}

func TestQuery_Sum(t *testing.T) {
	q := new(Query)
	q.Init()
	q.Attr("appId", "system").
		Attr("itemId", "cpu.load")
	q.Table("values.agent.local")

	sum, err := q.Sum("value.load1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("sum:", sum)
}

func TestQuery_Min(t *testing.T) {
	q := new(Query)
	q.Init()
	q.Attr("appId", "system").
		Attr("itemId", "cpu.load")
	q.Table("values.agent.local")

	min, err := q.Min("value.load1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("min:", min)
}

func TestQuery_Max(t *testing.T) {
	q := new(Query)
	q.Init()
	q.Attr("appId", "system").
		Attr("itemId", "cpu.load")
	q.Table("values.agent.local")

	max, err := q.Max("value.load1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("max:", max)
}

func TestQuery_Avg(t *testing.T) {
	q := new(Query)
	q.Init()
	q.Attr("appId", "system").
		Attr("itemId", "cpu.load")
	q.Table("values.agent.local")

	avg, err := q.Avg("value.load1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("avg:", avg)
}

func TestQuery_Group(t *testing.T) {
	q := NewQuery("values.agent.local")
	q.
		Debug().
		Attr("appId", "system").
		Attr("itemId", "cpu.load")

	result, err := q.Group("timeFormat.month", map[string]Expr{
		"load1":  NewAvgExpr("value.load1"),
		"load5":  "value.load5",
		"load15": NewMaxExpr("value.load15"),
	})
	if err != nil {
		t.Fatal(err)
	}

	logs.PrintAsJSON(result, t)
}

func TestQuery_Group_Dot(t *testing.T) {
	q := NewQuery("values.agent.local")
	q.
		Debug().
		Attr("appId", "system").
		Attr("itemId", "cpu.usage")

	result, err := q.Group("timeFormat.month", map[string]Expr{
		"usage.avg": NewAvgExpr("value.usage.avg"),
	})
	if err != nil {
		t.Fatal(err)
	}

	logs.PrintAsJSON(result, t)
}
