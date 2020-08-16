package teaagents

import (
	"github.com/TeaWeb/build/internal/teaagent/agentconst"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/processes"
	"github.com/iwind/TeaGo/timers"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"time"
)

// 检查更新
func checkNewVersion() {
	if runningAgent.Id == "local" {
		return
	}
	timers.Loop(120*time.Second, func(looper *timers.Looper) {
		if !runningAgent.AutoUpdates {
			return
		}

		//logs.Println("check new version")
		req, err := http.NewRequest(http.MethodGet, connectConfig.Master+"/api/agent/upgrade", nil)
		if err != nil {
			logs.Println("error:", err.Error())
			return
		}

		req.Header.Set("User-Agent", "TeaWeb-Agent/"+agentconst.AgentVersion)
		req.Header.Set("Tea-Agent-Id", connectConfig.Id)
		req.Header.Set("Tea-Agent-Key", connectConfig.Key)
		req.Header.Set("Tea-Agent-Version", agentconst.AgentVersion)
		req.Header.Set("Tea-Agent-Os", runtime.GOOS)
		req.Header.Set("Tea-Agent-Arch", runtime.GOARCH)

		resp, err := LongHTTPClient.Do(req)
		if err != nil {
			logs.Println("error:", err.Error())
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			logs.Println("error:status code not", http.StatusOK)
			return
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logs.Println("error:", err.Error())
			return
		}

		if len(data) > 1024 {
			logs.Println("start to upgrade")

			dir := Tea.Root + Tea.DS + "bin" + Tea.DS + "upgrade"
			dirFile := files.NewFile(dir)
			if !dirFile.Exists() {
				err := dirFile.Mkdir()
				if err != nil {
					logs.Println("error:", err.Error())
					return
				}
			}

			newVersion := resp.Header.Get("Tea-Agent-Version")
			filename := "teaweb-agent@" + newVersion
			if runtime.GOOS == "windows" {
				filename = "teaweb-agent@" + newVersion + ".exe"
			}
			file := files.NewFile(dir + "/" + filename)
			err = file.Write(data)
			if err != nil {
				logs.Println("error:", err.Error())
				return
			}

			err = file.Chmod(0777)
			if err != nil {
				logs.Println("error:", err.Error())
				return
			}

			// 停止当前
			if db != nil {
				db.Close()
				db = nil
			}

			// status server
			if statusServer != nil {
				statusServer.Shutdown()
			}

			// 启动
			logs.Println("start new version")
			proc := processes.NewProcess(dir+Tea.DS+filename, os.Args[1:]...)
			err = proc.StartBackground()
			if err != nil {
				logs.Println("error:", err.Error())
				return
			}

			logs.Println("exit to switch agent to latest version")
			time.Sleep(1 * time.Second)
			os.Exit(0)
		}
	})
}
