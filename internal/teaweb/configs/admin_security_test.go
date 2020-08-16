package configs

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestAdminSecurity_AllowIP(t *testing.T) {
	a := assert.NewAssertion(t)

	{
		security := new(AdminSecurity)
		a.IsNil(security.Validate())
		a.IsTrue(security.AllowIP("192.168.2.40"))
	}

	{
		security := new(AdminSecurity)
		security.Deny = []string{"192.168.2.40"}
		a.IsNil(security.Validate())
		a.IsFalse(security.AllowIP("192.168.2.40"))
	}

	{
		security := new(AdminSecurity)
		security.Deny = []string{"192.168.2.1/24"}
		a.IsNil(security.Validate())
		a.IsFalse(security.AllowIP("192.168.2.40"))
	}

	{
		security := new(AdminSecurity)
		security.Deny = []string{"192.168.3.1 / 24"}
		a.IsNil(security.Validate())
		a.IsTrue(security.AllowIP("192.168.2.40"))
	}

	{
		security := new(AdminSecurity)
		security.Deny = []string{"192.168.2.40,192.168.2.42"}
		a.IsNil(security.Validate())
		a.IsFalse(security.AllowIP("192.168.2.40"))
	}

	{
		security := new(AdminSecurity)
		security.Deny = []string{"192.168.2.41,192.168.2.42", "all"}
		a.IsNil(security.Validate())
		a.IsFalse(security.AllowIP("192.168.2.40"))
	}

	{
		security := new(AdminSecurity)
		security.Deny = []string{"0.0.0.0"}
		a.IsNil(security.Validate())
		a.IsFalse(security.AllowIP("192.168.2.40"))
	}

	{
		security := new(AdminSecurity)
		security.Deny = []string{}
		security.Allow = []string{"192.168.2.1"}
		a.IsNil(security.Validate())
		a.IsFalse(security.AllowIP("192.168.2.40"))
	}

	{
		security := new(AdminSecurity)
		security.Deny = []string{}
		security.Allow = []string{"192.168.2.1/24"}
		a.IsNil(security.Validate())
		a.IsTrue(security.AllowIP("192.168.2.40"))
		a.IsFalse(security.AllowIP("192.168.1.100"))
	}

	{
		security := new(AdminSecurity)
		security.Deny = []string{}
		security.Allow = []string{"192.168.2.1,192.168.2.100"}
		a.IsNil(security.Validate())
		a.IsTrue(security.AllowIP("192.168.2.40"))
		a.IsFalse(security.AllowIP("192.168.1.100"))
	}

	{
		security := new(AdminSecurity)
		security.Deny = []string{}
		security.Allow = []string{"all"}
		a.IsNil(security.Validate())
		a.IsTrue(security.AllowIP("192.168.2.40"))
	}

	{
		security := new(AdminSecurity)
		security.Deny = []string{}
		security.Allow = []string{"0.0.0.0"}
		a.IsNil(security.Validate())
		a.IsTrue(security.AllowIP("192.168.2.40"))
	}

	{
		security := new(AdminSecurity)
		security.Deny = []string{}
		security.Allow = []string{}
		a.IsNil(security.Validate())
		a.IsTrue(security.AllowIP("192.168.2.40"))
	}
}

func BenchmarkAdminSecurity_AllowIP(b *testing.B) {
	security := new(AdminSecurity)
	security.Deny = []string{}
	security.Allow = []string{"192.168.2.1/24"}
	_ = security.Validate()

	for i := 0; i < b.N; i++ {
		_ = security.AllowIP("192.168.2.40")
		_ = security.AllowIP("192.168.1.100")
	}
}
