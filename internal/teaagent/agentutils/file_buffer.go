package agentutils

import (
	"bufio"
	"github.com/iwind/TeaGo/logs"
	"os"
	"sync"
)

// 文件缓冲
type FileBuffer struct {
	MaxLines int

	debug bool

	path string

	reader *os.File
	writer *os.File

	file1 string
	file2 string

	isPositive bool
	hasWritten bool

	locker   sync.Mutex
	isClosed bool
}

// 创建新文件缓冲
func NewFileBuffer(path string) *FileBuffer {
	file1 := path + ".1.buf"
	file2 := path + ".2.buf"

	// 删除已有数据
	_, err := os.Stat(file1)
	if err == nil {
		err = os.Remove(file1)
		if err != nil {
			logs.Error(err)
		}
	}

	_, err = os.Stat(file2)
	if err == nil {
		err = os.Remove(file2)
		if err != nil {
			logs.Error(err)
		}
	}

	return &FileBuffer{
		path:  path,
		file1: file1,
		file2: file2,
	}
}

// 是否开启调试模式
func (this *FileBuffer) Debug() {
	this.debug = true
}

// 写入数据
func (this *FileBuffer) Write(data []byte) error {
	this.locker.Lock()

	if this.writer == nil {
		this.locker.Unlock()
		return nil
	}

	this.hasWritten = true

	_, err := this.writer.Write(data)
	if err == nil {
		err = this.writer.Sync()
	}

	this.locker.Unlock()
	return err
}

// 读取数据
func (this *FileBuffer) Read(f func(data []byte)) {
	if !this.reset() {
		if this.debug {
			logs.Println("skip")
		}
		return
	}

	if this.reader == nil {
		return
	}

	reader := bufio.NewReader(this.reader)
	i := 0
	for {
		line, err := reader.ReadBytes('\n')
		if len(line) > 0 {
			f(line[:len(line)-1])
			i++
		}
		if this.MaxLines > 0 && i > this.MaxLines {
			break
		}
		if err != nil {
			break
		}
	}

	return
}

// 是否有更多数据
func (this *FileBuffer) Next() bool {
	return !this.isClosed
}

// 关闭
func (this *FileBuffer) Close() {
	this.locker.Lock()

	this.close()
	this.isClosed = true

	this.locker.Unlock()
}

// 重置
func (this *FileBuffer) reset() (next bool) {
	next = true

	this.locker.Lock()
	defer this.locker.Unlock()

	// 如果没有写入，则不交换
	if this.reader != nil && this.writer != nil && !this.hasWritten {
		next = false
		return
	}

	this.hasWritten = false

	this.close()

	this.isPositive = !this.isPositive

	file1 := this.file1
	file2 := this.file2
	if this.isPositive {
		file1 = this.file2
		file2 = this.file1
	}

	if this.debug {
		logs.Println("reader", file1)
	}
	reader, err := os.OpenFile(file1, os.O_CREATE|os.O_RDONLY, 0777)
	if err != nil {
		logs.Error(err)
		return
	}
	this.reader = reader

	if this.debug {
		logs.Println("writer", file2)
	}
	writer, err := os.OpenFile(file2, os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		logs.Error(err)
		return
	}
	this.writer = writer
	return
}

// 关闭
func (this *FileBuffer) close() {
	if this.reader != nil {
		err := this.reader.Close()
		if err != nil {
			logs.Error(err)
		}
	}

	if this.writer != nil {
		err := this.writer.Close()
		if err != nil {
			logs.Error(err)
		}
	}
}
