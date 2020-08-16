package teaagents

import (
	"github.com/TeaWeb/build/internal/teaagent/agentconst"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/processes"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"os"
	"strings"
)

// 检查是否启动新版本
func shouldStartNewVersion() bool {
	if connectConfig.Id == "local" {
		return false
	}
	fileList := files.NewFile(Tea.Root + "/bin/upgrade/").List()
	latestVersion := agentconst.AgentVersion
	for _, f := range fileList {
		filename := f.Name()
		index := strings.Index(filename, "@")
		if index <= 0 {
			continue
		}
		version := strings.Replace(filename[index+1:], ".exe", "", -1)
		if stringutil.VersionCompare(latestVersion, version) < 0 {
			process := processes.NewProcess(Tea.Root+Tea.DS+"bin"+Tea.DS+"upgrade"+Tea.DS+filename, os.Args[1:]...)
			err := process.Start()
			if err != nil {
				logs.Println("[error]", err.Error())
				return false
			}

			err = process.Wait()
			if err != nil {
				logs.Println("[error]", err.Error())
				return false
			}

			return true
		}
	}
	return false
}
