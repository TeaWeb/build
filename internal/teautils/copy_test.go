package teautils

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestCopyStructObject(t *testing.T) {
	type Book struct {
		Name   string
		Price  int
		Year   int
		Author string
		press  string
	}

	book1 := &Book{
		Name:   "Hello Golang",
		Price:  100,
		Year:   2020,
		Author: "Liu",
		press:  "Beijing",
	}
	book2 := new(Book)
	CopyStructObject(book2, book1)
	logs.PrintAsJSON(book2, t)
	logs.PrintAsJSON(book1, t)
}
