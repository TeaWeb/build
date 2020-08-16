package agents

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teatesting"
	"github.com/iwind/TeaGo/logs"
	"net/http"
	"testing"
)

func TestURLConnectivitySource_Execute(t *testing.T) {
	source := NewURLConnectivitySource()
	source.URL = "https://baidu.com/"
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(value, t)
}

func TestURLConnectivitySource_ExecutePost(t *testing.T) {
	if !teatesting.RequireHTTPServer() {
		return
	}

	source := NewURLConnectivitySource()
	source.Method = http.MethodPost
	source.URL = "http://127.0.0.1:9991/webhook?hell=world"
	source.DataFormat = SourceDataFormatSingeLine
	source.Headers = []*shared.Variable{
		/**{
			Name:  "Content-Type",
			Value: "application/json",
		},**/
		{
			Name:  "Hello",
			Value: "World",
		},
	}
	source.Params = []*shared.Variable{
		{
			Name:  "name",
			Value: "lu",
		},
		{
			Name:  "age",
			Value: "20",
		},
	}
	source.TextBody = "Hello, World" // will be ignored because params is not empty
	err := source.Validate()
	if err != nil {
		t.Fatal(err)
	}
	result, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(result, t)
}

func TestURLConnectivitySource_ExecutePut(t *testing.T) {
	if !teatesting.RequireHTTPServer() {
		return
	}

	source := NewURLConnectivitySource()
	source.URL = "http://127.0.0.1:9991/webhook"
	source.Method = http.MethodPut
	source.DataFormat = SourceDataFormatSingeLine
	source.Headers = []*shared.Variable{
		{
			Name:  "Content-Type",
			Value: "application/json",
		},
	}
	source.TextBody = "HELLO, WORLD"
	result, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}

	logs.PrintAsJSON(result, t)
}
