package checkpoints

import (
	"bytes"
	"github.com/TeaWeb/build/internal/teawaf/requests"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestRequestJSONArgCheckpoint_RequestValue_Map(t *testing.T) {
	rawReq, err := http.NewRequest(http.MethodPost, "http://teaos.cn", bytes.NewBuffer([]byte(`
{
	"name": "lu",
	"age": 20,
	"books": [ "PHP", "Golang", "Python" ]
}
`)))
	if err != nil {
		t.Fatal(err)
	}

	req := requests.NewRequest(rawReq)
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	checkpoint := new(RequestJSONArgCheckpoint)
	t.Log(checkpoint.RequestValue(req, "name", nil))
	t.Log(checkpoint.RequestValue(req, "age", nil))
	t.Log(checkpoint.RequestValue(req, "Hello", nil))
	t.Log(checkpoint.RequestValue(req, "", nil))
	t.Log(checkpoint.RequestValue(req, "books", nil))
	t.Log(checkpoint.RequestValue(req, "books.1", nil))

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}

func TestRequestJSONArgCheckpoint_RequestValue_Array(t *testing.T) {
	rawReq, err := http.NewRequest(http.MethodPost, "http://teaos.cn", bytes.NewBuffer([]byte(`
[{
	"name": "lu",
	"age": 20,
	"books": [ "PHP", "Golang", "Python" ]
}]
`)))
	if err != nil {
		t.Fatal(err)
	}

	req := requests.NewRequest(rawReq)
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	checkpoint := new(RequestJSONArgCheckpoint)
	t.Log(checkpoint.RequestValue(req, "0.name", nil))
	t.Log(checkpoint.RequestValue(req, "0.age", nil))
	t.Log(checkpoint.RequestValue(req, "0.Hello", nil))
	t.Log(checkpoint.RequestValue(req, "", nil))
	t.Log(checkpoint.RequestValue(req, "0.books", nil))
	t.Log(checkpoint.RequestValue(req, "0.books.1", nil))

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}

func TestRequestJSONArgCheckpoint_RequestValue_Error(t *testing.T) {
	rawReq, err := http.NewRequest(http.MethodPost, "http://teaos.cn", bytes.NewBuffer([]byte(`
[{
	"name": "lu",
	"age": 20,
	"books": [ "PHP", "Golang", "Python" ]
}]
`)))
	if err != nil {
		t.Fatal(err)
	}

	req := requests.NewRequest(rawReq)
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	checkpoint := new(RequestJSONArgCheckpoint)
	t.Log(checkpoint.RequestValue(req, "0.name", nil))
	t.Log(checkpoint.RequestValue(req, "0.age", nil))
	t.Log(checkpoint.RequestValue(req, "0.Hello", nil))
	t.Log(checkpoint.RequestValue(req, "", nil))
	t.Log(checkpoint.RequestValue(req, "0.books", nil))
	t.Log(checkpoint.RequestValue(req, "0.books.1", nil))

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}
