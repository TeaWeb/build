// +build windows

package teaagents

import (
	"encoding/base64"
	"golang.org/x/sys/windows/registry"
	"runtime"
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
	key, _ := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	productName, _, _ := key.GetStringValue("ProductName")
	key.Close()

	if len(productName) > 0 {
		return productName
	}
	return "Windows"
}
