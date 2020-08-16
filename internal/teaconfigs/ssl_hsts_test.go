package teaconfigs

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestHSTSConfig(t *testing.T) {
	h := &HSTSConfig{}
	h.Validate()
	t.Log(h.HeaderValue())

	h.IncludeSubDomains = true
	h.Validate()
	t.Log(h.HeaderValue())

	h.Preload = true
	h.Validate()
	t.Log(h.HeaderValue())

	h.IncludeSubDomains = false
	h.Validate()
	t.Log(h.HeaderValue())

	h.MaxAge = 86400
	h.Validate()
	t.Log(h.HeaderValue())

	a := assert.NewAssertion(t)
	a.IsTrue(h.Match("abc.com"))

	h.Domains = []string{"abc.com"}
	h.Validate()
	a.IsTrue(h.Match("abc.com"))

	h.Domains = []string{"1.abc.com"}
	h.Validate()
	a.IsFalse(h.Match("abc.com"))
}
