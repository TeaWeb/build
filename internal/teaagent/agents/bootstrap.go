package teaagents

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teaagent/agentconfigs"
	"github.com/TeaWeb/build/internal/teaagent/agentconst"
	"github.com/TeaWeb/build/internal/teaagent/agentutils"
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var connectConfig *agentconfigs.AgentConfig = nil
var runningAgent *agents.AgentConfig = nil
var runningTasks = map[string]*Task{} // task id => task
var runningTasksLocker = sync.Mutex{}
var runningItems = map[string]*Item{} // item id => task
var runningItemsLocker = sync.Mutex{}
var isBooting = true
var connectionIsBroken = false

// 启动
func Start() {
	// 当前ROOT
	if !Tea.IsTesting() {
		exePath := agentutils.Executable()
		if strings.Contains(filepath.Base(exePath), "@") { // 是不是升级的文件
			Tea.UpdateRoot(filepath.Dir(filepath.Dir(filepath.Dir(exePath))))
		} else {
			Tea.UpdateRoot(filepath.Dir(filepath.Dir(exePath)))
		}
	}

	// 帮助
	if lists.ContainsAny(os.Args, "h", "-h", "help", "-help") {
		printHelp()
		return
	}

	// 版本号
	if lists.ContainsAny(os.Args, "version", "-v") {
		fmt.Println("v" + agentconst.AgentVersion)
		return
	}

	// 初始化
	if lists.ContainsString(os.Args, "init") {
		onInit()
		return
	}

	if len(os.Args) == 1 {
		writePid()
	}

	// 连接配置
	{
		config, err := agentconfigs.SharedAgentConfig()
		if err != nil {
			logs.Println("start failed: " + err.Error())
			return
		}
		connectConfig = config
	}

	// 检查新版本
	if shouldStartNewVersion() {
		return
	}

	// 启动
	if lists.ContainsString(os.Args, "start") {
		onStart()
		return
	}

	// 停止
	if lists.ContainsString(os.Args, "stop") {
		onStop()
		return
	}

	// 重启
	if lists.ContainsString(os.Args, "restart") {
		onStop()
		onStart()
		return
	}

	// 查看状态
	if lists.ContainsString(os.Args, "status") {
		onStatus()
		return
	}

	// 运行某个脚本
	if lists.ContainsAny(os.Args, "run") {
		runTaskOrItem()
		return
	}

	// 测试连接
	if lists.ContainsAny(os.Args, "test", "-t") {
		err := testConnection()
		if err != nil {
			logs.Println("error:", err.Error())
		} else {
			logs.Println("connection to master is ok")
		}
		return
	}

	// 日志
	if lists.ContainsAny(os.Args, "background", "-d") {
		writePid()

		logDir := files.NewFile(Tea.Root + "/logs")
		if !logDir.IsDir() {
			err := logDir.Mkdir()
			if err != nil {
				logs.Println(err.Error())
			}
		}

		fp, err := os.OpenFile(Tea.Root+"/logs/run.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(fp)
		} else {
			logs.Println(err)
		}
	}

	// Windows服务
	if lists.ContainsAny(os.Args, "service") && runtime.GOOS == "windows" {
		manager := agentutils.NewServiceManager("TeaWeb Agent", "TeaWeb Agent Manager")
		manager.Run()
	}

	logs.Println("agent starting ...")

	// 启动监听端口
	if connectConfig.Id != "local" && runtime.GOOS != "windows" {
		go startListening()
	}

	// 下载配置
	{
		// 遇到网络错误多次尝试
		countTries := 10
		for i := 0; i < countTries; i++ {
			err := downloadConfig()
			if err != nil {
				logs.Println("start failed: " + err.Error())

				if i < countTries-1 {
					time.Sleep(5 * time.Second)
					continue
				} else {
					return
				}
			}
			break
		}
	}

	// 启动
	logs.Println("agent boot tasks ...")
	bootTasks()
	isBooting = false

	// 定时
	logs.Println("agent schedule tasks ...")
	_ = scheduleTasks()

	// 监控项数据
	logs.Println("agent schedule items ...")
	_ = scheduleItems()

	// 检测Apps
	logs.Println("agent detect tasks ...")
	detectApps()

	// 检查更新
	checkNewVersion()

	// 推送日志
	go pushEvents()

	// 同步配置
	countTries := 0
	for {
		err := pullEvents()
		if err != nil {
			countTries++
			logs.Println("pull error:", err.Error())

			if countTries < 5 {
				time.Sleep(3 * time.Second)
			} else {
				time.Sleep(30 * time.Second)
			}
		} else {
			countTries = 0
		}
	}
}

// 初始化连接
func initConnection() {
	detectApps()
}

// 启动监听端口欧
var statusServer *Server = nil

func startListening() {
	statusServer = NewServer()
	err := statusServer.Start()
	if err != nil {
		logs.Error(err)
	}
}

// 检测App
func detectApps() {
	// 暂时不做任何事情
}

// 查找任务
func findTaskName(taskId string) string {
	if runningAgent == nil {
		return ""
	}
	_, task := runningAgent.FindTask(taskId)
	if task == nil {
		return ""
	}
	return task.Name
}

// 写入Pid
func writePid() {
	err := teautils.WritePid(Tea.Root + "/logs/pid")
	if err != nil {
		logs.Println(err.Error())
	}
}
