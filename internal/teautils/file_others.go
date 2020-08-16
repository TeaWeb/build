// +build !windows

package teautils

import (
	"os"
	"syscall"
)

// lock file
func LockFile(fp *os.File) error {
	return syscall.Flock(int(fp.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
}

func UnlockFile(fp *os.File) error {
	return syscall.Flock(int(fp.Fd()), syscall.LOCK_UN)
}
