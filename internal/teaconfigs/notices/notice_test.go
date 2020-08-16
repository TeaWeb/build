package notices

import "testing"

func TestNotice_Hash(t *testing.T) {
	notice := NewNotice()
	notice.Message = "Hello, World"
	notice.Hash()
	t.Log(notice.MessageHash)

	notice.Message = "Hello, World"
	notice.Hash()
	t.Log(notice.MessageHash)

	notice.Message = "Hello, World2"
	notice.Hash()
	t.Log(notice.MessageHash)
}
