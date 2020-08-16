package shared

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestRegexp(t *testing.T) {
	a := assert.NewAssertion(t)

	a.IsTrue(RegexpFloatNumber.MatchString("123"))
	a.IsTrue(RegexpFloatNumber.MatchString("123.456"))
	a.IsFalse(RegexpFloatNumber.MatchString(".456"))
	a.IsFalse(RegexpFloatNumber.MatchString("abc"))
	a.IsFalse(RegexpFloatNumber.MatchString("123."))
	a.IsFalse(RegexpFloatNumber.MatchString("123.456e7"))
}
