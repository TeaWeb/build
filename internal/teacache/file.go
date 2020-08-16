package teacache

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/timers"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/string"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var bytePool = teautils.NewBytePool(40960, 2048)

// 文件缓存管理器
type FileManager struct {
	Manager

	Capacity float64       // 容量
	Life     time.Duration // 有效期

	looper      *timers.Looper
	dir         string
	autoCreate  bool // 是否自动创建
	writeLocker sync.RWMutex
}

func NewFileManager() *FileManager {
	manager := &FileManager{}

	// 删除过期
	manager.looper = timers.Loop(30*time.Minute, func(looper *timers.Looper) {
		if len(manager.dir) == 0 {
			return
		}
		dirFile := files.NewFile(manager.dir)
		if !dirFile.IsDir() {
			return
		}
		for _, dirFile1 := range dirFile.List() {
			if !dirFile1.IsDir() {
				continue
			}
			for _, dirFile2 := range dirFile1.List() {
				for _, file := range dirFile2.List() {
					if file.Ext() != ".cache" {
						continue
					}
					reader, err := file.Reader()
					if err != nil {
						logs.Error(err)
						continue
					}
					data := reader.Read(10)
					if len(data) != 10 {
						_ = reader.Close()
						continue
					}
					timestamp := types.Int64(string(data))
					_ = reader.Close()
					if timestamp < time.Now().Unix()-100 { // 超时100秒以上的
						err := file.Delete()
						if err != nil {
							logs.Error(err)
						}
					}
				}

				time.Sleep(500 * time.Millisecond)
			}
		}
	})

	return manager
}

func (this *FileManager) SetOptions(options map[string]interface{}) {
	if this.Life <= 0 {
		this.Life = 1800 * time.Second
	}

	dir, ok := options["dir"]
	if ok {
		this.dir = types.String(dir)
		if !filepath.IsAbs(this.dir) {
			this.dir = Tea.Root + Tea.DS + this.dir
		}
	}

	autoCreate, ok := options["autoCreate"]
	if ok {
		this.autoCreate = types.Bool(autoCreate)
	}
}

// 写入
// 内容格式 timestamp | key length (20 bytes) | key | data (n bytes) |
func (this *FileManager) Write(key string, data []byte) error {
	if len(this.dir) == 0 {
		return errors.New("cache dir should not be empty")
	}

	this.writeLocker.Lock()
	defer this.writeLocker.Unlock()

	dirFile := files.NewFile(this.dir)
	if !dirFile.IsDir() {
		// 自动创建
		if this.autoCreate {
			err := dirFile.MkdirAll()
			if err != nil {
				return errors.New("can not create cache dir: " + err.Error())
			}
		} else {
			return errors.New("cache dir should be a valid dir")
		}
	}

	md5 := stringutil.Md5(key)
	newDir := files.NewFile(this.dir + Tea.DS + md5[:2] + Tea.DS + md5[2:4])
	if !newDir.Exists() {
		err := newDir.MkdirAll()
		if err != nil {
			return err
		}
	}

	path := newDir.Path() + Tea.DS + md5 + ".cache"
	fp, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_SYNC, 0666)
	if err != nil {
		return err
	}
	defer func() {
		_ = fp.Close()

		if err != nil {
			_ = os.Remove(path)
		}
	}()

	// 头部加入有效期
	var life = int64(this.Life.Seconds())
	if life <= 0 {
		life = 30 * 86400
	} else if life >= 365*86400 { // 最大值限制
		life = 365 * 86400
	}
	_, err = fp.WriteString(fmt.Sprintf("%d", time.Now().Unix()+life))
	if err != nil {
		return err
	}

	_, err = fp.WriteString("|")
	if err != nil {
		return err
	}

	_, err = fp.WriteString(strconv.Itoa(len(key)))
	if err != nil {
		return err
	}

	_, err = fp.WriteString("|")
	if err != nil {
		return err
	}

	_, err = fp.WriteString(key)
	if err != nil {
		return err
	}

	_, err = fp.WriteString("|")
	if err != nil {
		return err
	}

	_, err = fp.Write(data)

	return err
}

// 读取
func (this *FileManager) Read(key string) (data []byte, err error) {
	md5 := stringutil.Md5(key)
	fp, err := os.Open(this.dir + Tea.DS + md5[:2] + Tea.DS + md5[2:4] + Tea.DS + md5 + ".cache")
	if err != nil {
		return nil, ErrNotFound
	}
	defer func() {
		_ = fp.Close()
	}()

	this.writeLocker.RLock()
	defer this.writeLocker.RUnlock()

	r := bufio.NewReader(fp)
	key, ok := this.readHead(r)
	if !ok {
		return nil, ErrNotFound
	}

	buf := bytePool.Get()
	defer bytePool.Put(buf)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			data = append(data, buf[:n]...)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, ErrNotFound
		}
	}
	return data, nil
}

func (this *FileManager) readHead(r *bufio.Reader) (key string, ok bool) {
	timestampBytes, err := r.ReadBytes('|')
	if err != nil {
		return
	}
	if len(timestampBytes) > 12 {
		return
	}
	if types.Int64(string(timestampBytes[:len(timestampBytes)-1])) < time.Now().Unix() {
		return
	}

	keyLengthBytes, err := r.ReadBytes('|')
	if err != nil {
		return
	}

	keyLength := types.Int(string(keyLengthBytes[:len(keyLengthBytes)-1]))
	if keyLength > 0 {
		keyBytes := []byte{}
		for i := 0; i < keyLength; i++ {
			b, err := r.ReadByte()
			if err != nil {
				return
			}
			keyBytes = append(keyBytes, b)
		}

		key = string(keyBytes)
	}
	_, err = r.ReadByte()
	if err != nil {
		return
	}

	ok = true

	return
}

// 删除
func (this *FileManager) Delete(key string) error {
	if len(this.dir) == 0 {
		return errors.New("cache dir should not be empty")
	}

	this.writeLocker.Lock()
	defer this.writeLocker.Unlock()

	dirFile := files.NewFile(this.dir)
	if !dirFile.IsDir() {
		return errors.New("cache dir should be a valid dir")
	}

	md5 := stringutil.Md5(key)
	newDir := files.NewFile(this.dir + Tea.DS + md5[:2] + Tea.DS + md5[2:4])
	if !newDir.Exists() {
		return nil
	}

	newFile := files.NewFile(newDir.Path() + Tea.DS + md5 + ".cache")
	if !newFile.Exists() {
		return nil
	}
	return newFile.Delete()
}

// 删除key前缀
func (this *FileManager) DeletePrefixes(prefixes []string) (int, error) {
	if len(prefixes) == 0 {
		return 0, nil
	}
	// 检查目录是否存在
	info, err := os.Stat(this.dir)
	if err != nil || !info.IsDir() {
		return 0, nil
	}

	count := 0
	err = filepath.Walk(this.dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".cache") {
			return nil
		}

		fp, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() {
			_ = fp.Close()
		}()

		key, ok := this.readHead(bufio.NewReader(fp))
		if !ok {
			return nil
		}
		for _, prefix := range prefixes {
			if strings.HasPrefix(key, prefix) || strings.HasPrefix("http://"+key, prefix) || strings.HasPrefix("https://"+key, prefix) {
				_ = fp.Close()
				count++
				return os.Remove(path)
			}
		}

		return nil
	})

	return count, err
}

// 统计
func (this *FileManager) Stat() (size int64, countKeys int, err error) {
	// 检查目录是否存在
	info, err := os.Stat(this.dir)
	if err != nil || !info.IsDir() {
		return
	}

	err = filepath.Walk(this.dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".cache") {
			return nil
		}
		size += info.Size()
		countKeys++

		return nil
	})

	if err != nil {
		logs.Error(err)
	}

	return
}

// 清理
func (this *FileManager) Clean() error {
	dirReg := regexp.MustCompile("^[0-9a-f]{2}$")
	for _, file := range files.NewFile(this.dir).List() {
		if !file.IsDir() {
			continue
		}
		if !dirReg.MatchString(file.Name()) {
			continue
		}

		err := file.DeleteAll()
		if err != nil {
			logs.Error(err)
		}
	}
	return nil
}

func (this *FileManager) Close() error {
	//logs.Println("[cache]close cache policy instance: file")
	if this.looper != nil {
		this.looper.Stop()
		this.looper = nil
	}

	// TODO 删除所有文件和目录

	return nil
}
