package cmd

import (
	"fmt"
	"github.com/TeaWeb/build/internal/teacluster"
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"io"
	logger "log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

var sharedShell *WebShell = nil
var workerProcess *os.Process

// 命令行相关封装
type WebShell struct {
	ShouldStop bool
}

// 获取新对象
func NewWebShell() *WebShell {
	sharedShell = &WebShell{}
	return sharedShell
}

// 启动
func (this *WebShell) Start(server *TeaGo.Server) {
	// 重置ROOT
	this.resetRoot()

	// 执行参数
	if this.execArgs(os.Stdout) {
		this.ShouldStop = true
		return
	}

	// 信号
	teautils.ListenSignal(func(sig os.Signal) {
		logs.Println("[signal]" + sig.String())

		if sig == syscall.SIGHUP { // 重置
			configs.SharedAdminConfig().Reset()
		} else if sig == syscall.Signal(0x1e /**syscall.SIGUSR1**/) { // 刷新代理状态
			err := teaproxy.SharedManager.Restart()
			if err != nil {
				logs.Println("[error]" + err.Error())
			} else {
				proxyutils.FinishChange()
			}
		} else if sig == syscall.Signal(0x1f /**syscall.SIGUSR2**/) { // 同步
			node := teaconfigs.SharedNodeConfig()
			if node == nil {
				logs.Println("[cluster]not a node yet")
				return
			}

			if node.IsMaster() {
				logs.Println("[cluster]push items")
				teacluster.SharedManager.BuildSum()
				teacluster.SharedManager.PushItems()
			} else {
				logs.Println("[cluster]pull items")
				teacluster.SharedManager.BuildSum()
				teacluster.SharedManager.PullItems()
			}
		} else {
			if sig == syscall.SIGINT { // 终止进程
				if server != nil {
					server.Stop()
					time.Sleep(1 * time.Second)
				}
			}

			os.Exit(0)
		}
	}, syscall.SIGINT, syscall.SIGHUP, syscall.Signal(0x1e /**syscall.SIGUSR1**/), syscall.Signal(0x1f /**syscall.SIGUSR2**/), syscall.SIGTERM)
}

// 重置Root
func (this *WebShell) resetRoot() {
	if !Tea.IsTesting() {
		exePath, err := os.Executable()
		if err != nil {
			exePath = os.Args[0]
		}
		link, err := filepath.EvalSymlinks(exePath)
		if err == nil {
			exePath = link
		}
		fullPath, err := filepath.Abs(exePath)
		if err == nil {
			Tea.UpdateRoot(filepath.Dir(filepath.Dir(fullPath)))
		}
	}
	Tea.SetPublicDir(teautils.WebRoot() + Tea.DS + "public")
	Tea.SetViewsDir(teautils.WebRoot() + Tea.DS + "views")
	Tea.SetTmpDir(teautils.WebRoot() + Tea.DS + "tmp")
}

// 检查命令行参数
func (this *WebShell) execArgs(writer io.Writer) bool {
	if len(os.Args) == 1 {
		// 检查是否已经启动
		proc := this.checkPid()
		if proc != nil {
			this.write(writer, teaconst.TeaProductName+" is already running, pid:", proc.Pid)
			return true
		}

		return false
	}

	args := os.Args[1:]
	arg0 := ""
	if len(args) > 0 {
		arg0 = args[0]
	}
	if this.hasArg(arg0, "?", "help", "-help", "h", "-h") { // 帮助
		return this.ExecHelp(writer)
	}
	if this.hasArg(arg0, "-v", "version", "-version") { // 版本号
		return this.ExecVersion(writer)
	}
	if this.hasArg(arg0, "start") { // 启动
		return this.ExecStart(writer)
	}
	if this.hasArg(arg0, "stop") { // 停止
		return this.ExecStop(os.Stdout)
	}
	if this.hasArg(arg0, "reload") { // 重新加载代理配置
		return this.ExecReload(writer)
	}
	if this.hasArg(arg0, "restart") { // 重启
		return this.ExecRestart(writer)
	}
	if this.hasArg(arg0, "reset") { // 重置
		return this.ExecReset(writer)
	}
	if this.hasArg(arg0, "status") { // 状态
		return this.ExecStatus(writer)
	}
	if this.hasArg(arg0, "sync") { // 同步
		return this.ExecSync(writer)
	}
	if this.hasArg(arg0, "service") && runtime.GOOS == "windows" { // Windows服务
		return this.ExecService(writer)
	}
	if this.hasArg(arg0, "pprof") {
		return this.ExecPprof(writer)
	}
	if this.hasArg(arg0, "master") {
		return this.ExecMasterProcess(writer)
	}
	if this.hasArg(arg0, "worker") {
		return this.ExecWorker(writer)
	}

	if len(args) > 0 {
		this.write(writer, "Unknown command option '"+strings.Join(args, " ")+"', run './bin/"+teaconst.TeaProcessName+" -h' to lookup the usage.")
		return true
	}
	return false
}

// 帮助
func (this *WebShell) ExecHelp(writer io.Writer) bool {
	this.write(writer, teaconst.TeaProductName+" v"+teaconst.TeaVersion)
	this.write(writer, "Usage:", "\n   ./bin/"+teaconst.TeaProcessName+" [option]")
	this.write(writer, "")
	this.write(writer, "Options:")
	this.write(writer, fmt.Sprintf("  %-20s%s", "-h", ": print this help"))
	this.write(writer, fmt.Sprintf("  %-20s%s", "-v", ": print version"))
	this.write(writer, fmt.Sprintf("  %-20s%s", "start", ": start the server in background"))
	this.write(writer, fmt.Sprintf("  %-20s%s", "stop", ": stop the server"))
	this.write(writer, fmt.Sprintf("  %-20s%s", "reload", ": reload all proxy servers configs"))
	this.write(writer, fmt.Sprintf("  %-20s%s", "restart", ": restart the server"))
	this.write(writer, fmt.Sprintf("  %-20s%s", "reset", ": reset the server locker status"))
	this.write(writer, fmt.Sprintf("  %-20s%s", "status", ": print server status"))
	this.write(writer, fmt.Sprintf("  %-20s%s", "sync", ": sync config files with cluster"))
	this.write(writer, fmt.Sprintf("  %-20s%s", "pprof [address]", ": start pprof server"))
	this.write(writer, "")
	this.write(writer, "To run the server in foreground:", "\n   ./bin/"+teaconst.TeaProcessName+"\n")

	return true
}

// 版本号
func (this *WebShell) ExecVersion(writer io.Writer) bool {
	this.write(writer, teaconst.TeaProductName+" v"+teaconst.TeaVersion, "(build: "+runtime.Version(), runtime.GOOS, runtime.GOARCH+")")
	return true
}

// 启动
func (this *WebShell) ExecStart(writer io.Writer) bool {
	proc := this.checkPid()
	if proc != nil {
		this.write(writer, teaconst.TeaProductName+" already started, pid:", proc.Pid)
		return true
	}

	cmd := exec.Command(os.Args[0], "master")
	err := cmd.Start()
	if err != nil {
		this.write(writer, teaconst.TeaProductName+"  start failed:", err.Error())
		return true
	}
	this.write(writer, teaconst.TeaProductName+" started ok, pid:", cmd.Process.Pid)

	return true
}

// 停止
func (this *WebShell) ExecStop(writer io.Writer) bool {
	proc := this.checkPid()
	if proc == nil {
		this.write(writer, teaconst.TeaProductName+" not started yet")
		return true
	}

	// 通知关闭
	_ = teautils.NotifySignal(proc, syscall.SIGINT)
	time.Sleep(1 * time.Second)

	// 停止进程
	_ = proc.Kill()

	// 在Windows上经常不能及时释放资源
	_ = teautils.DeletePid(Tea.Root + "/bin/pid")
	this.write(writer, teaconst.TeaProductName+" stopped ok, pid:", proc.Pid)

	return true
}

// 重载代理配置
func (this *WebShell) ExecReload(writer io.Writer) bool {
	proc := this.checkPid()
	if proc == nil {
		this.write(writer, teaconst.TeaProductName+" not started yet")
		return true
	}
	err := teautils.NotifySignal(proc, syscall.Signal(0x1e /**syscall.SIGUSR1**/))
	if err != nil {
		this.write(writer, "[ERROR]"+err.Error())
		return true
	}
	this.write(writer, "reload success")
	return true
}

// 重启
func (this *WebShell) ExecRestart(writer io.Writer) bool {
	proc := this.checkPid()
	if proc != nil {
		// 通知关闭
		_ = teautils.NotifySignal(proc, syscall.SIGINT)
		time.Sleep(1 * time.Second)

		_ = proc.Kill()

		// 等待进程结束
		time.Sleep(100 * time.Millisecond)
	}

	cmd := exec.Command(os.Args[0], "master")
	err := cmd.Start()
	if err != nil {
		this.write(writer, teaconst.TeaProductName+" restart failed:", err.Error())
		return true
	}
	this.write(writer, teaconst.TeaProductName+" restarted ok, new pid:", cmd.Process.Pid)

	return true
}

// 重置
func (this *WebShell) ExecReset(writer io.Writer) bool {
	proc := this.checkPid()
	if proc == nil {
		this.write(writer, teaconst.TeaProductName+" not started yet")
		return true
	}
	err := teautils.NotifySignal(proc, syscall.SIGHUP)
	if err != nil {
		this.write(writer, "[ERROR]"+err.Error())
		return true
	}
	this.write(writer, "reset success")
	return true
}

// 状态
func (this *WebShell) ExecStatus(writer io.Writer) bool {
	proc := this.checkPid()
	if proc == nil {
		this.write(writer, teaconst.TeaProductName+" not started yet")
	} else {
		this.write(writer, teaconst.TeaProductName+" is running, pid: "+fmt.Sprintf("%d", proc.Pid))
	}
	return true
}

// 同步
func (this *WebShell) ExecSync(writer io.Writer) bool {
	proc := this.checkPid()
	if proc == nil {
		this.write(writer, teaconst.TeaProductName+" not started yet")
	} else {
		err := teautils.NotifySignal(proc, syscall.Signal(0x1f /**syscall.SIGUSR2**/))
		if err != nil {
			this.write(writer, "[ERROR]"+err.Error())
			return true
		}
		this.write(writer, "signal sent successfully")
	}
	return true
}

// 启动pprof
func (this *WebShell) ExecPprof(writer io.Writer) bool {
	addr := ":6060"
	if len(os.Args) == 3 {
		addr = os.Args[2]
	}
	this.write(writer, "===\nstart pprof server '"+addr+"'\n===")
	go func() {
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			this.write(writer, "[error]"+err.Error())
		}
	}()

	return false
}

// 启动守护进程
func (this *WebShell) ExecMasterProcess(writer io.Writer) bool {
	// 当前PID文件句柄
	err := this.writePid()
	if err != nil {
		logs.Println("[error]write pid file failed: '" + err.Error() + "'")
		return true
	}

	logWriter, _ := os.OpenFile(Tea.Root+"/logs/master.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if logWriter != nil {
		logger.SetOutput(logWriter)
	}

	teautils.ListenSignal(func(sig os.Signal) {
		logger.Println("[signal]" + sig.String())

		if sig == syscall.SIGHUP { // 重置
			if workerProcess != nil {
				_ = teautils.NotifySignal(workerProcess, sig)
			}
		} else if sig == syscall.Signal(0x1e /**syscall.SIGUSR1**/) { // 刷新代理状态
			if workerProcess != nil {
				_ = teautils.NotifySignal(workerProcess, sig)
			}
		} else if sig == syscall.Signal(0x1f /**syscall.SIGUSR2**/) { // 同步
			if workerProcess != nil {
				_ = teautils.NotifySignal(workerProcess, sig)
			}
		} else {
			if workerProcess != nil {
				_ = workerProcess.Kill()
			}

			// 删除PID
			err := teautils.DeletePid(Tea.Root + "/bin/pid")
			if err != nil {
				logs.Error(err)
			}
			os.Exit(0)
		}
	}, syscall.SIGINT, syscall.SIGHUP, syscall.Signal(0x1e /**syscall.SIGUSR1**/), syscall.Signal(0x1f /**syscall.SIGUSR2**/), syscall.SIGTERM)

	// 生成子进程，当前进程变成守护进程
	for {
		cmd := exec.Command(os.Args[0], "worker")
		err := cmd.Start()
		if err != nil {
			this.write(writer, teaconst.TeaProductName+"  start failed:", err.Error())
			break
		}
		workerProcess = cmd.Process
		_ = cmd.Wait()
	}
	return true
}

// 在后台执行
func (this *WebShell) ExecWorker(writer io.Writer) bool {
	ppid := os.Getppid()

	go func() {
		for {
			// 检查父进程是否已消失
			if os.Getppid() != ppid {
				os.Exit(0)
				return
			}
			parentProcess, _ := os.FindProcess(ppid)
			if parentProcess == nil {
				os.Exit(0)
				return
			}
			time.Sleep(5 * time.Second)
		}
	}()
	return false
}

// 写入PID
func (this *WebShell) writePid() error {
	return teautils.WritePid(Tea.Root + Tea.DS + "bin" + Tea.DS + "pid")
}

// 检查PID
func (this *WebShell) checkPid() *os.Process {
	return teautils.CheckPid(Tea.Root + "/bin/pid")
}

// 写入string到writer
func (this *WebShell) write(writer io.Writer, args ...interface{}) {
	_, _ = fmt.Fprintln(writer, args...)
}

// 判断命令
func (this *WebShell) hasArg(arg string, value ...string) bool {
	return lists.ContainsAny(value, arg)
}
