// +build darwin

package teautils

import (
	"syscall"
)

// set resource limit
func SetRLimit(limit uint64) error {
	rLimit := &syscall.Rlimit{}
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, rLimit)
	if err != nil {
		return err
	}

	if rLimit.Cur < limit {
		rLimit.Cur = limit
	}
	if rLimit.Max < limit {
		rLimit.Max = limit
	}
	return syscall.Setrlimit(syscall.RLIMIT_NOFILE, rLimit)
}

// set best resource limit value
func SetSuitableRLimit() {
	SetRLimit(4096 * 100) // 1M=100Files
}
