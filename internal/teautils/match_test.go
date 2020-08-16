package teautils

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestMatchKeyword(t *testing.T) {
	a := assert.NewAssertion(t)
	a.IsTrue(MatchKeyword("a b c", "a"))
	a.IsFalse(MatchKeyword("a b c", ""))
	a.IsTrue(MatchKeyword("abc", "BC"))
}
