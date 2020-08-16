// +build !windows

package teaagents

import (
	"encoding/base64"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/maps"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

var (
	currentOSName       string // 当前OS名称
	currentOSNameBase64 string // 当前OS名称Base64 encode的结果
)

// 获取系统发行版本信息
func retrieveOSName() string {
	if len(currentOSName) == 0 {
		currentOSName = retrieveOSNameInternal()
		if len(currentOSName) == 0 {
			currentOSName = runtime.GOOS
		}
	}
	return currentOSName
}

// 获取系统发行版本信息Base64结果
func retrieveOSNameBase64() string {
	if len(currentOSNameBase64) == 0 {
		currentOSNameBase64 = base64.StdEncoding.EncodeToString([]byte(retrieveOSName()))
	}
	return currentOSNameBase64
}

// 内部实际函数
func retrieveOSNameInternal() string {
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("sw_vers", "-productVersion")
		data, err := cmd.CombinedOutput()
		if err != nil {
			return "Mac OS X"
		}
		return "Mac OS X " + strings.TrimSpace(string(data))
	}

	if runtime.GOOS == "linux" {
		{
			osFile := files.NewFile("/etc/os-release")
			if osFile.Exists() {
				s, err := osFile.ReadAllString()
				if err == nil && len(s) > 0 {
					m := maps.Map{}
					for _, field := range strings.Split(s, "\n") {
						pieces := strings.SplitN(field, "=", 2)
						if len(pieces) != 2 {
							continue
						}
						m[pieces[0]] = strings.Trim(pieces[1], "\"")
					}
					name := m.GetString("NAME")
					version := m.GetString("VERSION_ID")
					if len(name) > 0 {
						return name + " " + version
					}
				}
			}
		}

		{
			etcFile := files.NewFile("/etc/redhat-release")
			if etcFile.Exists() {
				s, err := etcFile.ReadAllString()
				if err == nil && len(s) > 0 {
					return strings.TrimSpace(shortenOSName(s))
				}
			}
		}

		{
			etcFile := files.NewFile("/etc/issue")
			if etcFile.Exists() {
				s, err := etcFile.ReadAllString()
				if err == nil && len(s) > 0 {
					s = strings.Replace(s, "\\n", "", -1)
					s = strings.Replace(s, "\\l", "", -1)

					return strings.TrimSpace(shortenOSName(s))
				}
			}
		}

		{
			etcFile := files.NewFile("/etc/issue.net")
			if etcFile.Exists() {
				s, err := etcFile.ReadAllString()
				if err == nil && len(s) > 0 {
					s = strings.Replace(s, "\\n", "", -1)
					s = strings.Replace(s, "\\l", "", -1)
					s = strings.Replace(s, " GNU/Linux ", " ", -1)
					return strings.TrimSpace(shortenOSName(s))
				}
			}
		}
	} else if runtime.GOOS == "freebsd" {
		cmd := exec.Command("uname", "-a")
		data, err := cmd.CombinedOutput()
		if err != nil {
			return "FreeBSD"
		}

		result := regexp.MustCompile("FreeBSD\\s+([0-9.]+)").FindStringSubmatch(string(data))
		if len(result) == 0 {
			return "FreeBSD"
		}
		return "FreeBSD " + result[1]
	}

	return ""
}

// 简化OS的名称
func shortenOSName(osName string) string {
	osName = strings.Replace(osName, "CentOS Linux", "CentOS", -1)
	osName = strings.Replace(osName, "Red Hat Enterprise Linux Server", "RHEL", -1)
	osName = strings.Replace(osName, "SUSE Linux Enterprise Server", "SLES", -1)
	return osName
}
