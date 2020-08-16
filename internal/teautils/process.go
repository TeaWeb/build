package teautils

import (
	"fmt"
	"github.com/iwind/TeaGo/types"
	"io/ioutil"
	"os"
	"runtime"
)

var pidFileList = []*os.File{}

// 检查Pid
func CheckPid(path string) *os.Process {
	// windows上打开的文件是不能删除的
	if runtime.GOOS == "windows" {
		if os.Remove(path) == nil {
			return nil
		}
	}

	file, err := os.Open(path)
	if err != nil {
		return nil
	}

	defer func() {
		_ = file.Close()
	}()

	// 是否能取得Lock
	err = LockFile(file)
	if err == nil {
		_ = UnlockFile(file)
		return nil
	}

	pidBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil
	}
	pid := types.Int(string(pidBytes))

	if pid <= 0 {
		return nil
	}

	proc, _ := os.FindProcess(pid)
	return proc
}

// 写入Pid
func WritePid(path string) error {
	fp, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_RDONLY, 0666)
	if err != nil {
		return err
	}

	if runtime.GOOS != "windows" {
		err = LockFile(fp)
		if err != nil {
			return err
		}
	}
	pidFileList = append(pidFileList, fp) // hold the file pointers

	_, err = fp.WriteString(fmt.Sprintf("%d", os.Getpid()))
	if err != nil {
		return err
	}

	return nil
}

// 写入Ppid
func WritePpid(path string) error {
	fp, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_RDONLY, 0666)
	if err != nil {
		return err
	}

	if runtime.GOOS != "windows" {
		err = LockFile(fp)
		if err != nil {
			return err
		}
	}
	pidFileList = append(pidFileList, fp) // hold the file pointers

	_, err = fp.WriteString(fmt.Sprintf("%d", os.Getppid()))
	if err != nil {
		return err
	}

	return nil
}

// 删除Pid
func DeletePid(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, fp := range pidFileList {
		_ = UnlockFile(fp)
		_ = fp.Close()
	}
	return os.Remove(path)
}
