package agentutils

import (
	"os"
	"path/filepath"
)

// find executable file
func Executable() string {
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
		return fullPath
	}
	return exePath
}
