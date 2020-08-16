package tealogs

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/assert"
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestDecodePolicyStorage(t *testing.T) {
	a := assert.NewAssertion(t)

	policy := new(teaconfigs.AccessLogStoragePolicy)
	a.IsNil(DecodePolicyStorage(policy))

	policy.Type = StorageTypeFile
	policy.Options = map[string]interface{}{
		"format": "json",
		"path":   "/tmp/log${date}",
	}
	storage := DecodePolicyStorage(policy)
	if storage == nil {
		t.Fatal("decode failed")
	}
	logs.PrintAsJSON(storage, t)
}
