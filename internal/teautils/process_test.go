package teautils

import (
	"github.com/iwind/TeaGo/Tea"
	"testing"
)

func TestCheckPid(t *testing.T) {
	t.Log(CheckPid(Tea.Root + "/bin/pid"))
}
