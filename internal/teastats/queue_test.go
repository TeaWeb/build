package teastats

import (
	"github.com/iwind/TeaGo/assert"
	"github.com/iwind/TeaGo/maps"
	"testing"
)

func TestQueue_Add(t *testing.T) {

}

func TestQueue_equalStrings(t *testing.T) {
	queue := new(Queue)

	a := assert.NewAssertion(t)
	a.IsTrue(queue.equalStrings([]string{}, []string{}))
	a.IsTrue(queue.equalStrings([]string{"a"}, []string{"a"}))
	a.IsTrue(queue.equalStrings([]string{"a", "b", "c"}, []string{"c", "b", "a"}))
	a.IsFalse(queue.equalStrings([]string{"a"}, []string{}))
	a.IsFalse(queue.equalStrings([]string{"a", "b", "c"}, []string{"c", "a"}))
}

func TestQueue_Increase(t *testing.T) {
	m := maps.Map{
		"count": 123,
	}
	m1 := maps.Map{
		"count": 234.0,
	}

	queue := &Queue{}
	t.Log(queue.increase(m, m1))
}
