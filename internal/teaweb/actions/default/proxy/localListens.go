package proxy

import (
	"fmt"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/shirou/gopsutil/net"
	"os/exec"
	"runtime"
	"strings"
)

type LocalListensAction actions.Action

// 本地正在监听的地址
func (this *LocalListensAction) RunPost(params struct{}) {

	// 本机正在监听的地址
	result := []maps.Map{}
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		// ps?
		ps, err := exec.LookPath("ps")
		if err == nil {
			connections, _ := net.Connections("tcp")
			for _, conn := range connections {
				if conn.Status != "LISTEN" && len(conn.Raddr.IP) > 0 || conn.Raddr.Port > 0 || conn.Pid == 0 {
					continue
				}

				cmd := exec.Command(ps, "-p", fmt.Sprintf("%d", conn.Pid), "-o", "comm=")
				output, err := cmd.Output()
				if err != nil {
					continue
				}
				cmdName := strings.TrimSpace(string(output))
				if len(cmdName) == 0 {
					continue
				}
				if strings.Contains(cmdName, "/") {
					index := strings.LastIndex(cmdName, "/")
					cmdName = cmdName[index+1:]
				}
				addr := ""
				if conn.Laddr.IP == "*" {
					addr += "127.0.0.1"
				} else {
					addr += conn.Laddr.IP
				}
				addr += ":" + fmt.Sprintf("%d", conn.Laddr.Port)

				// 是否已经存在
				found := false
				for _, r := range result {
					if r.GetString("addr") == addr {
						found = true
						break
					}
				}
				if found {
					continue
				}

				result = append(result, maps.Map{
					"name": cmdName,
					"addr": addr,
				})
			}
		}
	}
	this.Data["result"] = result

	this.Success()
}
