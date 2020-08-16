package agents

import "testing"

func TestFindDataSourceInstance(t *testing.T) {
	t.Logf("%p", FindDataSourceInstance("file", map[string]interface{}{}))
	t.Logf("%p", FindDataSourceInstance("file", map[string]interface{}{}))
	t.Logf("%p", FindDataSourceInstance("file", map[string]interface{}{}))
}
