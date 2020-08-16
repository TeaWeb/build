package agents

import (
	"github.com/iwind/TeaGo/assert"
	"reflect"
	"testing"
)

func TestFuncFloat(t *testing.T) {
	a := assert.NewAssertion(t)

	a.IsTrue(FuncFloat() == 0)
	a.IsTrue(FuncFloat("123.456789") == 123.456789)
	a.IsTrue(FuncFloat("123.456789", "%.2f") == "123.46")
}

func TestFuncFormat(t *testing.T) {
	a := assert.NewAssertion(t)

	t.Log(FuncFormat(123.456789, "%.2f"))
	a.IsTrue(FuncFormat(123.456789, "%.2f") == "123.46")
	a.IsTrue(FuncFormat(123.456123, "%.3f%s", "HELLO") == "123.456HELLO")
}

func TestFuncAppend(t *testing.T) {
	a := assert.NewAssertion(t)

	a.IsTrue(FuncAppend("a", "b", "c") == "abc")
}

func TestFuncCall(t *testing.T) {
	{
		funcType := reflect.ValueOf(FuncFloat)
		values := funcType.Call([]reflect.Value{reflect.ValueOf("123.4567890123")})
		for _, v := range values {
			t.Log(v.Interface())
		}
	}

	{
		funcType := reflect.ValueOf(FuncAppend)
		values := funcType.Call([]reflect.Value{reflect.ValueOf("a"), reflect.ValueOf("b"), reflect.ValueOf("c")})
		for _, v := range values {
			t.Log(v.Interface())
		}
	}
}

func TestRunFuncExpr(t *testing.T) {
	a := assert.NewAssertion(t)

	t.Log(RunFuncExpr(123.456789, []byte("float|format('%.3f%s,%s', 'a', 'b')|append('a1','b2','c3')")))
	t.Log(RunFuncExpr(123.456, []byte("append('a', 'b2', 'c345\"', \"6789'10'\", '\\'Hello\\'\\\"')")))
	t.Log(RunFuncExpr(123.456, []byte("append('78910')|float|format('%.6f')")))
	t.Log(RunFuncExpr(123.456, []byte("format('%.2f') | append('a', 'b', 'c\td')")))
	t.Log(RunFuncExpr(123.456, []byte("append(123.456, true, 'a')")))
	t.Log(RunFuncExpr(123.456, []byte("")))
	t.Log(RunFuncExpr(123.456, []byte(" ")))

	{
		v, err := RunFuncExpr(123.456, []byte(" format('%.2f')"))
		if err != nil {
			t.Fatal(err)
		}
		a.IsTrue(v == "123.46")
	}

	t.Log(RunFuncExpr(nil, []byte("float")))

	{
		v, err := RunFuncExpr("123.456", []byte("round"))
		if err != nil {
			t.Fatal(err)
		}
		t.Log(v)
		a.IsTrue(v == int64(123))
	}

	{
		v, err := RunFuncExpr("123.567", []byte("round"))
		if err != nil {
			t.Fatal(err)
		}
		t.Log(v)
		a.IsTrue(v == int64(124))
	}

	{
		v, err := RunFuncExpr("123.567", []byte("round(2)"))
		if err != nil {
			t.Fatal(err)
		}
		t.Log(v)
		a.IsTrue(v == "123.57")
	}

	{
		v, err := RunFuncExpr("123.4567123", []byte("round(4)"))
		if err != nil {
			t.Fatal(err)
		}
		t.Log(v)
		a.IsTrue(v == "123.4567")
	}

	{
		v, err := RunFuncExpr("123.567", []byte("ceil"))
		if err != nil {
			t.Fatal(err)
		}
		t.Log(v)
		a.IsTrue(v == int64(124))
	}

	{
		v, err := RunFuncExpr("123.567", []byte("floor"))
		if err != nil {
			t.Fatal(err)
		}
		t.Log(v)
		a.IsTrue(v == int64(123))
	}
}

func TestCheckLiteral(t *testing.T) {
	a := assert.NewAssertion(t)
	{
		_, err := checkLiteral("abc")
		a.IsNotNil(err)
		t.Log(err)
	}
	{
		result, err := checkLiteral("true")
		a.IsNil(err)
		a.IsTrue(result == true)
	}
	{
		result, err := checkLiteral("false")
		a.IsNil(err)
		a.IsTrue(result == false)
	}

	{
		result, err := checkLiteral("null")
		a.IsNil(err)
		a.IsTrue(result == nil)
	}

	{
		result, err := checkLiteral("nil")
		a.IsNil(err)
		a.IsTrue(result == nil)
	}

	{
		result, err := checkLiteral("123")
		a.IsNil(err)
		t.Log(result)
		a.IsTrue(result == int64(123))
	}

	{
		result, err := checkLiteral("+123")
		a.IsNil(err)
		t.Log(result)
		a.IsTrue(result == int64(123))
	}

	{
		result, err := checkLiteral("-123")
		a.IsNil(err)
		t.Log(result)
		a.IsTrue(result == int64(-123))
	}

	{
		result, err := checkLiteral("123.456")
		a.IsNil(err)
		t.Log(result)
		a.IsTrue(result == 123.456)
	}
	{
		result, err := checkLiteral("-123.456")
		a.IsNil(err)
		t.Log(result)
		a.IsTrue(result == -123.456)
	}
}
