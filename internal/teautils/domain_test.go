package teautils

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestMatchDomain(t *testing.T) {
	a := assert.NewAssertion(t)
	{
		ok := MatchDomains([]string{}, "example.com")
		a.IsFalse(ok)
	}

	{
		ok := MatchDomains([]string{"example.com"}, "example.com")
		a.IsTrue(ok)
	}

	{
		ok := MatchDomains([]string{"www.example.com"}, "example.com")
		a.IsFalse(ok)
	}

	{
		ok := MatchDomains([]string{".example.com"}, "www.example.com")
		a.IsTrue(ok)
	}

	{
		ok := MatchDomains([]string{".example.com"}, "a.www.example.com")
		a.IsTrue(ok)
	}

	{
		ok := MatchDomains([]string{".example.com"}, "a.www.example123.com")
		a.IsFalse(ok)
	}

	{
		ok := MatchDomains([]string{"*.example.com"}, "www.example.com")
		a.IsTrue(ok)
	}

	{
		ok := MatchDomains([]string{"*.*.com"}, "www.example.com")
		a.IsTrue(ok)
	}

	{
		ok := MatchDomains([]string{"www.*.com"}, "www.example.com")
		a.IsTrue(ok)
	}

	{
		ok := MatchDomains([]string{"gallery.*.com"}, "www.example.com")
		a.IsFalse(ok)
	}

	{
		ok := MatchDomains([]string{"~\\w+.example.com"}, "www.example.com")
		a.IsTrue(ok)
	}

	{
		ok := MatchDomains([]string{"~\\w+.example.com"}, "a.www.example.com")
		a.IsTrue(ok)
	}

	{
		ok := MatchDomains([]string{"~^\\d+.example.com$"}, "www.example.com")
		a.IsFalse(ok)
	}

	{
		ok := MatchDomains([]string{"~^\\d+.example.com$"}, "123.example.com")
		a.IsTrue(ok)
	}
}
