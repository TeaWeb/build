package logbuffer

import (
	"github.com/iwind/TeaGo/logs"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

// 文件Buffer
// 日志文件：PREFIX.INDEX.log 其中INDEX为0-N
type Buffer struct {
	prefix           string
	writingFileIndex int
	readingFileIndex int
	chunkSize        int

	files map[int]*File
}

// 获得FileBuffer对象
func NewBuffer(prefix string) *Buffer {
	buf := &Buffer{
		prefix:           prefix,
		writingFileIndex: 0,
		chunkSize:        64 * (1 << 20), // 64M
		files:            map[int]*File{},
	}
	buf.init()
	return buf
}

// 初始化
func (this *Buffer) init() {
	matches, err := filepath.Glob(this.prefix + ".*.log")
	if err != nil {
		logs.Error(err)
	} else {
		for _, match := range matches {
			if len(match) > 0 {
				err := os.Remove(match)
				if err != nil {
					logs.Error(err)
				}
			}
		}
	}
}

// 写入数据
func (this *Buffer) Write(data []byte) (n int, err error) {
	file, ok := this.files[this.writingFileIndex]
	if !ok {
		file = NewFile(this.filename(this.writingFileIndex))
		this.files[this.writingFileIndex] = file
	}

	n, err = file.Write(data)
	if err != nil {
		return
	}

	if file.Size() > this.chunkSize {
		_ = file.Sync()
		this.writingFileIndex++
	}

	return
}

// 读取数据
func (this *Buffer) Read() (data []byte, err error) {
	file, ok := this.files[this.readingFileIndex]
	if !ok {
		return
	}
	data, err = file.Read()
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		if this.readingFileIndex != this.writingFileIndex {
			// 关闭
			err := file.Close()
			if err != nil {
				logs.Error(err)
			}

			// 删除
			err = file.Delete()
			if err != nil {
				logs.Error(err)
			}

			this.readingFileIndex++
		}
		err = nil
	}
	return
}

// 文件名
func (this *Buffer) filename(index int) string {
	return this.prefix + "." + strconv.Itoa(index) + ".log"
}

// 关闭
func (this *Buffer) Close() error {
	var resultErr error
	for _, file := range this.files {
		err := file.Close()
		if err != nil {
			resultErr = err
		}
	}
	return resultErr
}
