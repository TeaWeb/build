package teaconfigs

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestGzipConfig_MatchContentType(t *testing.T) {
	a := assert.NewAssertion(t)

	{
		gzip := &GzipConfig{}
		a.IsNil(gzip.Validate())
		a.IsTrue(gzip.MatchContentType("text/html"))
	}

	{
		gzip := &GzipConfig{}
		a.IsNil(gzip.Validate())
		a.IsTrue(gzip.MatchContentType("text/html; charset=utf-8"))
	}

	{
		gzip := &GzipConfig{}
		gzip.MimeTypes = []string{"text/*"}
		a.IsNil(gzip.Validate())
		a.IsTrue(gzip.MatchContentType("text/html; charset=utf-8"))
	}

	{
		gzip := &GzipConfig{}
		gzip.MimeTypes = []string{"text/*"}
		a.IsNil(gzip.Validate())
		a.IsFalse(gzip.MatchContentType("application/json; charset=utf-8"))
	}

	{
		gzip := &GzipConfig{}
		gzip.MimeTypes = []string{"text/*", "image/*"}
		a.IsNil(gzip.Validate())
		a.IsTrue(gzip.MatchContentType("image/jpeg; charset=utf-8"))
	}
}
