package teaagents

import (
	"flag"
	"fmt"
	"github.com/TeaWeb/build/internal/teaagent/agentconfigs"
	"github.com/iwind/TeaGo/logs"
	"strings"
)

// 初始化
func onInit() {
	flag.Parse()

	masterAddr := ""
	groupKey := ""

	for _, arg := range flag.Args() {
		pieces := strings.SplitN(arg, "=", 2)
		if len(pieces) == 2 {
			piece1 := strings.TrimPrefix(strings.TrimSpace(pieces[0]), "-")
			piece2 := strings.TrimSpace(pieces[1])
			switch piece1 {
			case "master":
				masterAddr = piece2
			case "group":
				groupKey = piece2
			}
		}
	}

	if len(masterAddr) == 0 {
		logs.Println("[error]'master' argument should not be empty")
		return
	}

	if len(groupKey) == 0 {
		logs.Println("[error]'group' argument should not be empty")
		return
	}

	agentConfig, err := agentconfigs.SharedAgentConfig()
	if err != nil {
		agentConfig = &agentconfigs.AgentConfig{
			Master:   masterAddr,
			Id:       "",
			Key:      "",
			GroupKey: groupKey,
		}
	} else {
		agentConfig.Master = masterAddr
		agentConfig.GroupKey = groupKey
	}

	if len(agentConfig.Id) == 0 {
		agentConfig.Id = "ID"
	}
	if len(agentConfig.Key) == 0 {
		agentConfig.Key = "KEY"
	}

	connectConfig = agentConfig
	err = testConnection()
	if err != nil {
		logs.Println("[error]" + err.Error())
		return
	}

	fmt.Println("init successfully, now you can start the agent")
}
