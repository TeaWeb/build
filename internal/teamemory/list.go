package teamemory

type List struct {
	head *Item
	end  *Item
}

func NewList() *List {
	return &List{}
}

func (this *List) Add(item *Item) {
	if item == nil {
		return
	}
	if this.end != nil {
		this.end.Next = item
		item.Prev = this.end
		item.Next = nil
	}
	this.end = item
	if this.head == nil {
		this.head = item
	}
}

func (this *List) Remove(item *Item) {
	if item == nil {
		return
	}
	if item.Prev != nil {
		item.Prev.Next = item.Next
	}
	if item.Next != nil {
		item.Next.Prev = item.Prev
	}
	if item == this.head {
		this.head = item.Next
	}
	if item == this.end {
		this.end = item.Prev
	}

	item.Prev = nil
	item.Next = nil
}

func (this *List) Len() int {
	l := 0
	for e := this.head; e != nil; e = e.Next {
		l ++
	}
	return l
}

func (this *List) Range(f func(item *Item) (goNext bool)) {
	for e := this.head; e != nil; e = e.Next {
		goNext := f(e)
		if !goNext {
			break
		}
	}
}

func (this *List) Reset() {
	this.head = nil
	this.end = nil
}
