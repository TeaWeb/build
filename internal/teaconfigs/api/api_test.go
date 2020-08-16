package api

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestAPIMatch(t *testing.T) {
	api := NewAPI()
	api.Path = "/user/:id"
	err := api.Validate()
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(api.Match("/hello"))
	t.Log(api.Match("/user"))
	t.Log(api.Match("/user/"))
	t.Log(api.Match("/user/123"))

	api.Path = "/user/:id/:name"
	api.Validate()
	t.Log(api.Match("/user/123/liu"))
}

func TestAPI_IsWatching(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()

	api := NewAPI()
	api.Path = "/hello"

	api.StartWatching()

	a.IsTrue(api.IsWatching())
	api.StopWatching()
	a.IsFalse(api.IsWatching())
}
