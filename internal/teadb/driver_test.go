package teadb

import (
	"github.com/TeaWeb/build/internal/teaconfigs/audits"
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/iwind/TeaGo/logs"
	"testing"
	"time"
)

func TestDriver_Test(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	err := SharedDB().Test()
	if err != nil {
		t.Log("error:", err.Error())
	} else {
		t.Log("ok")
	}
}

func TestDriver_FindOne(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	{
		query := NewQuery("teaweb.logs.audit")
		query.Attr("id", 1)
		one, err := SharedDB().FindOne(query, new(audits.Log))
		if err != nil {
			t.Fatal(err)
		}
		logs.PrintAsJSON(one, t)
	}

	{
		query := NewQuery("teaweb.logs.audit")
		query.Attr("id", 2)
		one, err := SharedDB().FindOne(query, new(audits.Log))
		if err != nil {
			t.Fatal(err)
		}
		logs.PrintAsJSON(one, t)
	}
}

func TestDriver_FindOnes(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	query := NewQuery("teaweb.logs.audit")
	query.Attr("id", []int{1, 2})
	ones, err := SharedDB().FindOnes(query, new(audits.Log))
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(ones, t)
}

func TestDriver_InsertOne(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	{
		log := new(audits.Log)
		log.Timestamp = time.Now().Unix()
		log.Action = "LOGIN"
		log.Description = "login from beijing"
		log.Username = "admin"
		log.Options = map[string]string{
			"a": "1",
			"b": "2",
		}
		err := SharedDB().InsertOne("teaweb.logs.audit", log)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("ok")
	}

	{
		log := new(audits.Log)
		log.Timestamp = time.Now().Unix()
		log.Action = "LOGIN"
		log.Description = "login from shanghai"
		log.Username = "admin"
		log.Options = map[string]string{
			"a": "1",
			"b": "2",
		}
		err := SharedDB().InsertOne("teaweb.logs.audit", log)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("ok")
	}
}

func TestDriver_InsertOnes(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	ones := []*audits.Log{
		{
			Action:      "LOGIN",
			Timestamp:   time.Now().Unix(),
			Description: "login",
			Username:    "admin",
		},
		{
			Action:      "LOGOUT",
			Timestamp:   time.Now().Unix(),
			Description: "logout",
			Username:    "user",
		},
	}
	err := SharedDB().InsertOnes("teaweb.logs.audit", ones)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestDriver_DeleteOnes(t *testing.T) {
	q := NewQuery("teaweb.logs.audit")
	q.Gt("id", 3)
	err := SharedDB().DeleteOnes(q)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestDriver_Count(t *testing.T) {
	{
		q := NewQuery("teaweb.logs.audit")
		count, err := SharedDB().Count(q)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("count:", count)
	}

	{
		q := NewQuery("teaweb.logs.audit")
		q.Gt("id", 5)
		count, err := SharedDB().Count(q)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("count:", count)
	}
}

func TestDriver_Avg(t *testing.T) {
	q := NewQuery("teaweb.logs.audit")
	count, err := SharedDB().Avg(q, "id")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("avg:", count)
}

func TestDriver_Sum(t *testing.T) {
	q := NewQuery("teaweb.logs.audit")
	count, err := SharedDB().Sum(q, "id")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("sum:", count)
}

func TestDriver_Min(t *testing.T) {
	q := NewQuery("teaweb.logs.audit")
	count, err := SharedDB().Min(q, "id")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("min:", count)
}

func TestDriver_Max(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	q := NewQuery("teaweb.logs.audit")
	count, err := SharedDB().Max(q, "id")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("max:", count)
}

func TestDriver_Group(t *testing.T) {
	if !teatesting.RequireDBAvailable() {
		return
	}

	q := NewQuery("teaweb.logs.audit")
	q.Gt("id", 0)
	q.Desc("_id")
	result, err := SharedDB().Group(q, "action", map[string]Expr{
		"timestamp": NewAvgExpr("timestamp"),
		"id":        NewSumExpr("id"),
		"action":    "action",
		"k2":        "options.k2",
		"a.b.c":     "options.k4",
	})
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(result, t)
}
