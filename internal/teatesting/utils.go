package teatesting

import (
	"os"
	"strings"
)

// 是否是全局测试：go test ./...
func IsGlobal() bool {
	if len(os.Args) == 0 {
		return false
	}
	return strings.HasSuffix(os.Args[0], ".test") || strings.HasSuffix(os.Args[0], ".test.exe")
}
