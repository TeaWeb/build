package teaconfigs

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestPageConfig_Match(t *testing.T) {
	a := assert.NewAssertion(t)

	{
		page := NewPageConfig()
		page.Status = []string{"200"}
		page.Validate()
		a.IsTrue(page.Match(200))
		a.IsFalse(page.Match(201))
	}

	{
		page := NewPageConfig()
		page.Status = []string{"4xx", "5xx"}
		page.Validate()
		a.IsFalse(page.Match(200))
		a.IsTrue(page.Match(401))
		a.IsTrue(page.Match(404))
		a.IsTrue(page.Match(500))
		a.IsTrue(page.Match(505))
	}
}
