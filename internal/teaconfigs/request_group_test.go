package teaconfigs

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestRequestGroup_Match(t *testing.T) {
	a := assert.NewAssertion(t)

	formatter := func(source string) string {
		if source == "${remoteAddr}" {
			return "192.168.1.100"
		}
		if source == "${arg.id}" {
			return "20"
		}
		return ""
	}

	{
		group := NewRequestGroup()
		group.Validate()
		a.IsTrue(group.Match(formatter))
	}

	{
		group := NewRequestGroup()
		{
			ipRange := shared.NewIPRangeConfig()
			ipRange.Type = shared.IPRangeTypeRange
			ipRange.Param = "${remoteAddr}"
			ipRange.IPFrom = "192.168.1.1"
			ipRange.IPTo = "192.168.1.200"
			group.AddIPRange(ipRange)
		}
		err := group.Validate()
		if err != nil {
			t.Fatal(err)
		}
		a.IsTrue(group.Match(formatter))
	}

	{
		group := NewRequestGroup()
		{
			ipRange := shared.NewIPRangeConfig()
			ipRange.Type = shared.IPRangeTypeRange
			ipRange.Param = "${remoteAddr}"
			ipRange.IPFrom = "192.168.1.1"
			ipRange.IPTo = "192.168.1.100"
			group.AddIPRange(ipRange)
		}
		err := group.Validate()
		if err != nil {
			t.Fatal(err)
		}
		a.IsTrue(group.Match(formatter))
	}

	{
		group := NewRequestGroup()
		{
			ipRange := shared.NewIPRangeConfig()
			ipRange.Type = shared.IPRangeTypeRange
			ipRange.Param = "${remoteAddr}"
			ipRange.IPFrom = "192.168.1.1"
			ipRange.IPTo = "192.168.1.99"
			group.AddIPRange(ipRange)
		}
		err := group.Validate()
		if err != nil {
			t.Fatal(err)
		}
		a.IsFalse(group.Match(formatter))
	}

	{
		group := NewRequestGroup()
		{
			ipRange := shared.NewIPRangeConfig()
			ipRange.Type = shared.IPRangeTypeCIDR
			ipRange.Param = "${remoteAddr}"
			ipRange.CIDR = "192.168.1.1/24"
			group.AddIPRange(ipRange)
		}
		err := group.Validate()
		if err != nil {
			t.Fatal(err)
		}
		a.IsTrue(group.Match(formatter))
	}

	{
		group := NewRequestGroup()
		{
			cond := shared.NewRequestCond()
			cond.Param = "${arg.id}"
			cond.Operator = shared.RequestCondOperatorGtFloat
			cond.Value = "19"
			group.AddCond(cond)
		}
		{
			ipRange := shared.NewIPRangeConfig()
			ipRange.Type = shared.IPRangeTypeCIDR
			ipRange.Param = "${remoteAddr}"
			ipRange.CIDR = "192.168.1.1/24"
			group.AddIPRange(ipRange)
		}
		err := group.Validate()
		if err != nil {
			t.Fatal(err)
		}
		a.IsTrue(group.Match(formatter))
	}
}
