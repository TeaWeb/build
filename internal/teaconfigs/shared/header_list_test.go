package shared

import (
	"fmt"
	"testing"
)

func TestHeaderList_FormatHeaders(t *testing.T) {
	list := &HeaderList{}

	for i := 0; i < 5; i++ {
		list.AddRequestHeader(&HeaderConfig{
			On:    true,
			Name:  "A" + fmt.Sprintf("%d", i),
			Value: "ABCDEFGHIJ${name}KLM${hello}NEFGHIJILKKKk",
		})
	}

	err := list.ValidateHeaders()
	if err != nil {
		t.Fatal(err)
	}
}
