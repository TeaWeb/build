package configs

import (
	"github.com/iwind/TeaGo/assert"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"testing"
)

func TestSharedAdminConfig(t *testing.T) {
	adminConfig := SharedAdminConfig()
	t.Logf("%#v", adminConfig)
}

func TestAdminConfig_ComparePassword(t *testing.T) {
	a := assert.NewAssertion(t)

	{
		config := AdminConfig{}
		a.IsTrue(config.ComparePassword(stringutil.Md5("123456"), "123456"))
	}

	{
		config := AdminConfig{}
		a.IsTrue(config.ComparePassword(stringutil.Md5("123456"), "clear:123456"))
	}

	{
		config := AdminConfig{}
		a.IsTrue(config.ComparePassword(stringutil.Md5("123456"), "md5:"+stringutil.Md5("123456")))
		a.IsFalse(config.ComparePassword(stringutil.Md5("123456789"), "md5:"+stringutil.Md5("123456")))
	}
}
