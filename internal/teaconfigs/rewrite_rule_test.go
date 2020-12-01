package teaconfigs

import (
	"github.com/iwind/TeaGo/assert"
	"github.com/iwind/TeaGo/rands"
	"sync"
	"testing"
	"time"
)

func TestRewriteRule(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()

	rule := RewriteRule{
		Pattern: "/(hello)/(world)",
		Replace: "/${1}/${2}",
	}
	a.IsNil(rule.Validate())

	{
		_, _, b := rule.Match("/hello/worl", func(source string) string {
			return source
		})
		a.IsFalse(b)
		a.Log("proxy:", rule.TargetProxy())
		a.Log("url:", rule.TargetURL())
	}

	{
		_, _, b := rule.Match("/hello/world", func(source string) string {
			return source
		})
		a.IsTrue(b)
		a.Log("proxy:", rule.TargetProxy())
		a.Log("url:", rule.TargetURL())
	}

	{
		r := NewRewriteRule()
		r.Replace = "http://127.0.0.1${0}"
		r.Pattern = ".*"
		r.Validate()
		u, _, b := r.Match("/hello", func(source string) string {
			return source
		})
		a.Log(b)
		a.Log(u)
	}
}

func TestRewriteRule_NamedMatch(t *testing.T) {
	r := NewRewriteRule()
	r.Replace = "http://127.0.0.1/${1}/${last}/${ni}"
	r.Pattern = "/(\\w+)/(?P<last>\\w+)/(?P<ni>\\w+)"
	r.Validate()

	before := time.Now()
	count := 100
	for i := 0; i < count; i++ {
		s, _, b := r.Match("/hello/world/ni", func(source string) string {
			return source
		})
		if i == 0 {
			if b {
				t.Log("matched:", s)
			} else {
				t.Log("not matched")
			}
		}
	}
	t.Log(float64(count) / (time.Since(before).Seconds()))
}

func TestRewriteRule_NamedMatchConcurrent(t *testing.T) {
	r := NewRewriteRule()
	r.Replace = "http://127.0.0.1/${1}/${last}/${ni}"
	r.Pattern = "/(\\w+)/(?P<last>\\w+)/(?P<ni>\\w+)"
	r.Validate()

	threads := 1000
	count := 1000
	wg := sync.WaitGroup{}
	wg.Add(threads * count)
	fails := 0
	var locker sync.Mutex
	for i := 0; i < threads; i++ {
		go func() {
			for j := 0; j < count; j++ {
				func() {
					defer wg.Done()

					var randomString = rands.HexString(16)

					replace, _, b := r.Match("/hello/world/"+randomString, func(source string) string {
						return source
					})
					if !b {
						locker.Lock()
						fails++
						locker.Unlock()
					} else if replace != "http://127.0.0.1/hello/world/"+randomString {
						locker.Lock()
						fails++
						locker.Unlock()
					}
				}()
			}
		}()
	}
	wg.Wait()
	if fails > 0 {
		t.Log("fail")
	} else {
		t.Log("success")
	}
}

func TestRewriteRule_CaseInsensitive(t *testing.T) {
	a := assert.NewAssertion(t)

	r := NewRewriteRule()
	r.Replace = "http://127.0.0.1${0}"
	r.Pattern = "(?i)/index.php"
	r.Validate()
	_, _, ok := r.Match("/index.php", func(source string) string {
		return source
	})
	a.IsTrue(ok)

	_, _, ok = r.Match("/INDEX.php", func(source string) string {
		return source
	})
	a.IsTrue(ok)
}

func TestRewriteRule_Slashes(t *testing.T) {
	a := assert.NewAssertion(t)

	r := NewRewriteRule()
	r.Replace = "http://127.0.0.1/${0}"
	r.Pattern = "(?i)/index.php"
	r.Validate()
	replace, _, ok := r.Match("/index.php", func(source string) string {
		return source
	})
	a.IsTrue(ok)
	t.Log(replace)
}

func TestRewriteRuleProxy(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()

	rule := RewriteRule{
		Pattern: "/(hello)/(world)",
		Replace: "proxy://lb001/${1}/${2}",
	}
	a.IsNil(rule.Validate())

	replace, _, b := rule.Match("/hello/world", func(source string) string {
		return source
	})
	a.IsTrue(b)
	a.IsTrue(rule.TargetProxy() == "lb001")
	a.IsTrue(replace == "/hello/world")
}
