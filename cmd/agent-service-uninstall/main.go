package main

import (
	"github.com/TeaWeb/build/internal/teaagent/agentutils"
	"log"
)

// 卸载服务
func main() {
	log.Println("uninstalling ...")
	manager := agentutils.NewServiceManager("TeaWeb Agent", "TeaWeb Agent Manager")
	err := manager.Uninstall()
	if err != nil {
		log.Println("ERROR: " + err.Error())
		manager.Close()
		manager.PauseWindow()
		return
	}

	log.Println("uninstalled service successfully")
	log.Println("done.")

	manager.PauseWindow()
}
