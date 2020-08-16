package teaagents

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/Tea"
	"strconv"
)

// lookup status
func onStatus() {
	proc := teautils.CheckPid(Tea.Root + "/logs/pid")
	if proc == nil {
		fmt.Println("Agent not started yet")
		return
	}

	fmt.Println("Agent is running, pid:", strconv.Itoa(proc.Pid))
}
