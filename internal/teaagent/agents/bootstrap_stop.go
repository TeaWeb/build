package teaagents

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teaagent/agentconst"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/Tea"
)

// 停止
func onStop() {
	pidFile := Tea.Root + "/logs/pid"
	proc := teautils.CheckPid(pidFile)
	if proc == nil {
		fmt.Println(agentconst.AgentProductName + " agent is not running")
		return
	}

	_ = proc.Kill()
	fmt.Println(agentconst.AgentProductName+" stopped pid:", proc.Pid)

	_ = teautils.DeletePid(pidFile)
}
