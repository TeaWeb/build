package teaagents

import (
	"github.com/TeaWeb/build/internal/teaagent/agentconst"
	"github.com/TeaWeb/build/internal/teaagent/agentutils"
)

// 打印帮助
func printHelp() {
	agentutils.NewCommandHelp().
		Product(agentconst.AgentProductName).
		Version(agentconst.AgentVersion).
		Usage("./bin/"+agentconst.AgentProcessName+" [options]").
		Option("-h", "show this help").
		Option("-v|version", "show agent version").
		Option("start", "start agent in background").
		Option("stop", "stop running agent").
		Option("restart", "restart the agent").
		Option("status", "lookup agent status").
		Option("run [TASK ID]", "run task").
		Option("run [ITEM ID]", "run app item").
		Option("init -master=[MASTER SERVER] -group=[GROUP KEY]", "register agent to master server and specified group").
		Append("To run agent in foreground\n   bin/teaweb-agent").
		Print()
}
