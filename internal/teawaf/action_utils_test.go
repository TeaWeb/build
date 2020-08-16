package teawaf

import (
	"github.com/iwind/TeaGo/assert"
	"github.com/iwind/TeaGo/maps"
	"runtime"
	"testing"
)

func TestFindActionInstance(t *testing.T) {
	a := assert.NewAssertion(t)

	t.Logf("ActionBlock: %p", FindActionInstance(ActionBlock, nil))
	t.Logf("ActionBlock: %p", FindActionInstance(ActionBlock, nil))
	t.Logf("ActionGoGroup: %p", FindActionInstance(ActionGoGroup, nil))
	t.Logf("ActionGoGroup: %p", FindActionInstance(ActionGoGroup, nil))
	t.Logf("ActionGoSet: %p", FindActionInstance(ActionGoSet, nil))
	t.Logf("ActionGoSet: %p", FindActionInstance(ActionGoSet, nil))
	t.Logf("ActionGoSet: %#v", FindActionInstance(ActionGoSet, maps.Map{"groupId": "a", "setId": "b",}))

	a.IsTrue(FindActionInstance(ActionGoSet, nil) != FindActionInstance(ActionGoSet, nil))
}

func BenchmarkFindActionInstance(b *testing.B) {
	runtime.GOMAXPROCS(1)
	for i := 0; i < b.N; i++ {
		FindActionInstance(ActionGoSet, nil)
	}
}
