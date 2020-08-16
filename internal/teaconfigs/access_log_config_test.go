package teaconfigs

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestAccessLogConfig_Match(t *testing.T) {
	a := assert.NewAssertion(t)

	{
		accessLog := NewAccessLogConfig()
		a.IsNil(accessLog.Validate())
		a.IsTrue(accessLog.Match(100))
		a.IsTrue(accessLog.Match(200))
		a.IsTrue(accessLog.Match(300))
		a.IsTrue(accessLog.Match(400))
		a.IsTrue(accessLog.Match(500))
	}

	{
		accessLog := NewAccessLogConfig()
		accessLog.Status1 = false
		accessLog.Status2 = false
		a.IsNil(accessLog.Validate())
		a.IsFalse(accessLog.Match(100))
		a.IsFalse(accessLog.Match(200))
		a.IsTrue(accessLog.Match(300))
		a.IsTrue(accessLog.Match(400))
		a.IsTrue(accessLog.Match(500))
	}

	{
		accessLog := NewAccessLogConfig()
		accessLog.Status3 = false
		accessLog.Status4 = false
		accessLog.Status5 = false
		a.IsNil(accessLog.Validate())
		a.IsTrue(accessLog.Match(100))
		a.IsTrue(accessLog.Match(200))
		a.IsFalse(accessLog.Match(300))
		a.IsFalse(accessLog.Match(400))
		a.IsFalse(accessLog.Match(500))
	}
}
