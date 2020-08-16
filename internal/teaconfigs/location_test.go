package teaconfigs

import (
	"fmt"
	"github.com/iwind/TeaGo/assert"
	"sync"
	"testing"
)

func TestLocationConfig_Match(t *testing.T) {
	location := NewLocation()
	err := location.Validate()
	if err != nil {
		t.Fatal(err)
	}

	a := assert.NewAssertion(t).Quiet()

	location.Pattern = "/hell"
	a.IsNotError(location.Validate())

	_, b := location.Match("/hello", func(source string) string {
		return source
	})
	a.IsTrue(b)

	location.Pattern = "/hello"
	a.IsNotError(location.Validate())

	_, b = location.Match("/hello", func(source string) string {
		return source
	})
	a.IsTrue(b)

	location.Pattern = "~ ^/\\w+$"
	a.IsNotError(location.Validate())
	_, b = location.Match("/hello", func(source string) string {
		return source
	})
	a.IsTrue(b)

	location.Pattern = "!~ ^/HELLO$"
	a.IsNotError(location.Validate())
	_, b = location.Match("/hello", func(source string) string {
		return source
	})
	a.IsTrue(b)

	location.Pattern = "~* ^/HELLO$"
	a.IsNotError(location.Validate())

	_, b = location.Match("/hello", func(source string) string {
		return source
	})
	a.IsTrue(b)

	location.Pattern = "!~* ^/HELLO$"
	a.IsNotError(location.Validate())
	_, b = location.Match("/hello", func(source string) string {
		return source
	})
	a.IsFalse(b)

	location.Pattern = "= /hello"
	a.IsNotError(location.Validate())
	_, b = location.Match("/hello", func(source string) string {
		return source
	})
	a.IsTrue(b)
}

func TestLocationConfig_Match_Concurrent(t *testing.T) {
	location := NewLocation()
	location.Pattern = `~ /(?P<name>\w+)`
	err := location.Validate()
	if err != nil {
		t.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	count := 10000
	wg.Add(count)

	for i := 0; i < count; i ++ {
		go func(i int) {
			defer wg.Done()

			u := "hello" + fmt.Sprintf("%d", i)
			result, b := location.Match("/"+u, func(source string) string {
				return source
			})
			if !b || result["name"] != u {
				t.Fatal("u:", u, result["name"])
			}
		}(i)
	}
	wg.Wait()
}
