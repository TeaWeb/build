package agents

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teatesting"
	"net/http"
	"testing"
)

func TestWebHookSource_ExecuteGet(t *testing.T) {
	if !teatesting.RequireHTTPServer() {
		return
	}

	webHook := NewWebHookSource()
	webHook.Method = http.MethodGet
	webHook.URL = "http://127.0.0.1:9991/webhook?hell=world"
	webHook.DataFormat = SourceDataFormatSingeLine
	err := webHook.Validate()
	if err != nil {
		t.Fatal(err)
	}
	result, err := webHook.Execute(map[string]string{
		"host": "127.0.0.1",
		"port": "3306",
	})
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(result)
	}
}

func TestWebHookSource_ExecutePost(t *testing.T) {
	if !teatesting.RequireHTTPServer() {
		return
	}

	webHook := NewWebHookSource()
	webHook.Method = http.MethodPost
	webHook.URL = "http://127.0.0.1:9991/webhook?hell=world"
	webHook.DataFormat = SourceDataFormatSingeLine
	webHook.Headers = []*shared.Variable{
		/**{
			Name:  "Content-Type",
			Value: "application/json",
		},**/
		{
			Name:  "Hello",
			Value: "World",
		},
	}
	webHook.Params = []*shared.Variable{
		{
			Name:  "name",
			Value: "lu",
		},
		{
			Name:  "age",
			Value: "20",
		},
	}
	webHook.TextBody = "Hello, World" // will be ignored because params is not empty
	err := webHook.Validate()
	if err != nil {
		t.Fatal(err)
	}
	result, err := webHook.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestWebHookSource_ExecutePost2(t *testing.T) {
	if !teatesting.RequireHTTPServer() {
		return
	}

	webHook := NewWebHookSource()
	webHook.Method = http.MethodPost
	webHook.URL = "http://127.0.0.1:9991/webhook?hell=world"
	webHook.DataFormat = SourceDataFormatSingeLine
	webHook.Headers = []*shared.Variable{
		/**{
			Name:  "Content-Type",
			Value: "application/json",
		},**/
		{
			Name:  "Hello",
			Value: "World",
		},
	}
	webHook.TextBody = "Hello, World"
	err := webHook.Validate()
	if err != nil {
		t.Fatal(err)
	}
	result, err := webHook.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestWebHookSource_ExecutePut(t *testing.T) {
	if !teatesting.RequireHTTPServer() {
		return
	}

	webHook := NewWebHookSource()
	webHook.URL = "http://127.0.0.1:9991/webhook"
	webHook.Method = http.MethodPut
	webHook.DataFormat = SourceDataFormatSingeLine
	webHook.Headers = []*shared.Variable{
		{
			Name:  "Content-Type",
			Value: "application/json",
		},
	}
	webHook.TextBody = "HELLO, WORLD"
	result, err := webHook.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
}
