package log

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"sync"
)

type RuntimeAction actions.Action

var runtimeLogOffset = int64(0)
var runtimeLogLocker = sync.Mutex{}

// 系统日志
func (this *RuntimeAction) Run(params struct{}) {
	runtimeLogLocker.Lock()
	runtimeLogOffset = 0
	runtimeLogLocker.Unlock()

	this.Show()
}

// 读取数据
func (this *RuntimeAction) RunPost() {
	runtimeLogLocker.Lock()
	runtimeLogLocker.Unlock()

	reader, err := files.NewReader(Tea.LogFile("teaweb.log"))
	if err != nil {
		return
	}
	defer reader.Close()

	if runtimeLogOffset == 0 {
		n, err := reader.Seek(int64(-4096), files.WhenceEnd)
		if err != nil {
			runtimeLogOffset = 0
		} else {
			runtimeLogOffset = n

			// 往前找换行符
			if n > 0 {
				for {
					newOffset, err := reader.Seek(-1, files.WhenceCurrent)
					if err != nil {
						runtimeLogOffset = newOffset
						break
					} else {
						char := string(reader.ReadByte())

						if char == "\n" || char == "\r" {
							runtimeLogOffset = newOffset + 1
							break
						}

						reader.Seek(-1, files.WhenceCurrent)
					}
				}
			}
		}
	} else {
		_, err := reader.Seek(runtimeLogOffset, files.WhenceStart)
		if err != nil {
			return
		}
	}

	data := reader.Read(4096)
	runtimeLogOffset += int64(len(data))

	this.Data["data"] = string(data)

	this.Success()
}
