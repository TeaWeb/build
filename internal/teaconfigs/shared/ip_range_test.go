package shared

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestGeoConfig_Contains(t *testing.T) {
	a := assert.NewAssertion(t)

	{
		geo := NewIPRangeConfig()
		geo.Type = IPRangeTypeRange
		geo.IPFrom = "192.168.1.100"
		geo.IPTo = "192.168.1.110"
		a.IsNil(geo.Validate())
		a.IsTrue(geo.Contains("192.168.1.100"))
		a.IsTrue(geo.Contains("192.168.1.101"))
		a.IsTrue(geo.Contains("192.168.1.110"))
		a.IsFalse(geo.Contains("192.168.1.111"))
	}

	{
		geo := NewIPRangeConfig()
		geo.Type = IPRangeTypeCIDR
		geo.CIDR = "192.168.1.1/24"
		a.IsNil(geo.Validate())
		a.IsTrue(geo.Contains("192.168.1.100"))
		a.IsFalse(geo.Contains("192.168.2.100"))
	}

	{
		geo := NewIPRangeConfig()
		geo.Type = IPRangeTypeCIDR
		geo.CIDR = "192.168.1.1/16"
		a.IsNil(geo.Validate())
		a.IsTrue(geo.Contains("192.168.2.100"))
	}
}

func TestParseIPRange(t *testing.T) {
	a := assert.NewAssertion(t)

	{
		_, err := ParseIPRange("")
		a.IsNotNil(err)
	}

	{
		r, err := ParseIPRange("192.168.1.100")
		a.IsNil(err)
		a.IsTrue(r.IPFrom == r.IPTo)
		a.IsTrue(r.IPFrom == "192.168.1.100")
		a.IsTrue(r.Contains("192.168.1.100"))
		a.IsFalse(r.Contains("192.168.1.99"))
	}

	{
		r, err := ParseIPRange("192.168.1.100/24")
		a.IsNil(err)
		a.IsTrue(r.CIDR == "192.168.1.100/24")
		a.IsTrue(r.Contains("192.168.1.100"))
		a.IsTrue(r.Contains("192.168.1.99"))
		a.IsFalse(r.Contains("192.168.2.100"))
	}

	{
		r, err := ParseIPRange("192.168.1.100, 192.168.1.200")
		a.IsNil(err)
		a.IsTrue(r.IPFrom == "192.168.1.100")
		a.IsTrue(r.IPTo == "192.168.1.200")
		a.IsTrue(r.Contains("192.168.1.100"))
		a.IsTrue(r.Contains("192.168.1.150"))
		a.IsFalse(r.Contains("192.168.2.100"))
	}

	{
		r, err := ParseIPRange("192.168.1.100-192.168.1.200")
		a.IsNil(err)
		a.IsTrue(r.IPFrom == "192.168.1.100")
		a.IsTrue(r.IPTo == "192.168.1.200")
		a.IsTrue(r.Contains("192.168.1.100"))
		a.IsTrue(r.Contains("192.168.1.150"))
		a.IsFalse(r.Contains("192.168.2.100"))
	}

	{
		r, err := ParseIPRange("all")
		a.IsNil(err)
		a.IsTrue(r.Type == IPRangeTypeAll)
		a.IsTrue(r.Contains("192.168.1.100"))
		a.IsTrue(r.Contains("192.168.1.150"))
		a.IsTrue(r.Contains("192.168.2.100"))
	}

	{
		r, err := ParseIPRange("192.168.1.*")
		a.IsNil(err)
		if r != nil {
			a.IsTrue(r.Type == IPRangeTypeWildcard)
			a.IsTrue(r.Contains("192.168.1.100"))
			a.IsFalse(r.Contains("192.168.2.100"))
		}
	}

	{
		r, err := ParseIPRange("192.168.*.*")
		a.IsNil(err)
		if r != nil {
			a.IsTrue(r.Type == IPRangeTypeWildcard)
			a.IsTrue(r.Contains("192.168.1.100"))
			a.IsTrue(r.Contains("192.168.2.100"))
		}
	}
}

func BenchmarkIPRangeConfig_Contains(b *testing.B) {
	r, err := ParseIPRange("192.168.1.*")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		_ = r.Contains("192.168.1.100")
	}
}
