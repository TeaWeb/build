package teadb

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestCurrentDB(t *testing.T) {
	db := SharedDB()
	t.Log(db)
}

func TestInitTable(t *testing.T) {
	a := assert.NewAssertion(t)
	a.IsFalse(isInitializedTable("a"))
	a.IsTrue(isInitializedTable("a"))
	a.IsTrue(isInitializedTable("a"))
}
