package teamemory

import "testing"

func TestList(t *testing.T) {
	l := &List{}

	var e1 *Item = nil
	{
		e := &Item{
			ValueInt64: 1,
		}
		l.Add(e)
		e1 = e
	}

	var e2 *Item = nil
	{
		e := &Item{
			ValueInt64: 2,
		}
		l.Add(e)
		e2 = e
	}

	var e3 *Item = nil
	{
		e := &Item{
			ValueInt64: 3,
		}
		l.Add(e)
		e3 = e
	}

	var e4 *Item = nil
	{
		e := &Item{
			ValueInt64: 4,
		}
		l.Add(e)
		e4 = e
	}

	l.Remove(e1)
	//l.Remove(e2)
	//l.Remove(e3)
	l.Remove(e4)

	for e := l.head; e != nil; e = e.Next {
		t.Log(e.ValueInt64)
	}

	t.Log("e1, e2, e3, e4, head, end:", e1, e2, e3, e4)
	if l.head != nil {
		t.Log("head:", l.head.ValueInt64)
	} else {
		t.Log("head: nil")
	}
	if l.end != nil {
		t.Log("end:", l.end.ValueInt64)
	} else {
		t.Log("end: nil")
	}
}
