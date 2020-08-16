package teaconfigs

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestWildcardStatus_Match(t *testing.T) {
	a := assert.NewAssertion(t)

	{
		s := NewWildcardStatus("123")
		a.IsTrue(s.Match(123))
	}

	{
		s := NewWildcardStatus("5x3")
		a.IsTrue(s.Match(523))
		a.IsTrue(s.Match(553))
	}

	{
		s := NewWildcardStatus("5xx")
		a.IsTrue(s.Match(523))
		a.IsTrue(s.Match(553))
		a.IsTrue(s.Match(500))
		a.IsFalse(s.Match(5000))
		a.IsFalse(s.Match(50))
		a.IsFalse(s.Match(400))
	}
}
