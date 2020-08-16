package teaconfigs

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestFastcgiConfig_Address(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()

	f := NewFastcgiConfig()

	{
		f.Pass = "127.0.0.1"
		err := f.Validate()
		a.IsNil(err)
		a.IsTrue(f.Network() == "tcp")
		a.IsTrue(f.Address() == "127.0.0.1:9000")
	}

	{
		f.Pass = "9000"
		err := f.Validate()
		a.IsNil(err)
		a.IsTrue(f.Network() == "tcp")
		a.IsTrue(f.Address() == "127.0.0.1:9000")
	}

	{
		f.Pass = "192.168.1.1:9000"
		err := f.Validate()
		a.IsNil(err)
		a.IsTrue(f.Network() == "tcp")
		a.IsTrue(f.Address() == "192.168.1.1:9000")
	}

	{
		f.Pass = ":9000"
		err := f.Validate()
		a.IsNil(err)
		a.IsTrue(f.Network() == "tcp")
		a.IsTrue(f.Address() == "127.0.0.1:9000")
	}

	{
		f.Pass = "unix:/tmp/php-fpm.sock"
		err := f.Validate()
		a.IsNil(err)
		a.IsTrue(f.Network() == "unix")
		a.IsTrue(f.Address() == "/tmp/php-fpm.sock")
	}

	{
		f.Pass = "/tmp/php-fpm.sock"
		err := f.Validate()
		a.IsNil(err)
		a.IsTrue(f.Network() == "unix")
		a.IsTrue(f.Address() == "/tmp/php-fpm.sock")
	}

	{
		f.Pass = "/tmp/php-fpm"
		err := f.Validate()
		a.IsNil(err)
		a.IsTrue(f.Network() == "unix")
		a.IsTrue(f.Address() == "/tmp/php-fpm")
	}

	{
		f.Pass = "../tmp/php-fpm"
		err := f.Validate()
		a.IsNil(err)
		a.IsTrue(f.Network() == "unix")
		a.IsTrue(f.Address() == "../tmp/php-fpm")
	}

	{
		f.Pass = "192.168.1."
		err := f.Validate()
		a.IsNotNil(err)
	}
}
