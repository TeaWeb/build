package teacache

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/timers"
	"github.com/iwind/TeaGo/types"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"path/filepath"
	"strings"
	"time"
)

// Leveldb缓存管理器
type LevelDBManager struct {
	Manager

	Capacity float64       // 容量
	Life     time.Duration // 有效期

	dir    string // 数据库所在目录
	db     *leveldb.DB
	looper *timers.Looper
}

func NewLevelDBManager() *LevelDBManager {
	manager := &LevelDBManager{}
	manager.looper = timers.Loop(30*time.Minute, func(looper *timers.Looper) {
		manager.CleanExpired()
	})
	return manager
}

func (this *LevelDBManager) SetOptions(options map[string]interface{}) {
	if this.Life <= 0 {
		this.Life = 1800 * time.Second
	}

	dir, found := options["dir"]
	if found {
		this.dir = types.String(dir)
	}

	if !filepath.IsAbs(this.dir) {
		this.dir = Tea.Root + "/" + this.dir
	}

	dirFile := files.NewFile(this.dir)
	if !dirFile.Exists() {
		err := dirFile.MkdirAll()
		if err != nil {
			logs.Error(err)
			return
		}
	}

	db, err := leveldb.OpenFile(this.dir+"/cache.leveldb", nil)
	if err != nil {
		logs.Error(err)
	} else {
		this.db = db
	}
}

func (this *LevelDBManager) Write(key string, data []byte) error {
	if this.db == nil {
		return errors.New("leveldb nil pointer found")
	}
	life := fmt.Sprintf("%d_", time.Now().Unix()+int64(this.Life.Seconds()))
	return this.db.Put(append([]byte("KEY"), []byte(key)...), append([]byte(life), data...), nil)
}

func (this *LevelDBManager) Read(key string) (data []byte, err error) {
	if this.db == nil {
		return nil, errors.New("leveldb nil pointer found")
	}
	newKey := append([]byte("KEY"), []byte(key)...)
	data, err = this.db.Get(newKey, nil)
	if err == leveldb.ErrNotFound {
		return nil, ErrNotFound
	}

	index := bytes.Index(data, []byte("_"))
	if index < -1 {
		this.db.Delete(newKey, nil)
		return nil, ErrNotFound
	}

	timestamp := types.Int64(data[:index])
	if timestamp < time.Now().Unix() {
		this.db.Delete(newKey, nil)
		return nil, ErrNotFound
	}

	data = data[index+1:]

	return
}

// 删除
func (this *LevelDBManager) Delete(key string) error {
	if this.db == nil {
		return errors.New("leveldb nil pointer found")
	}
	return this.db.Delete(append([]byte("KEY"), []byte(key)...), nil)
}

// 删除key前缀
func (this *LevelDBManager) DeletePrefixes(prefixes []string) (int, error) {
	if len(prefixes) == 0 {
		return 0, nil
	}
	if this.db == nil {
		return 0, nil
	}

	it := this.db.NewIterator(util.BytesPrefix([]byte("KEY")), nil)
	defer it.Release()

	count := 0
	for it.Next() {
		keyBytes := it.Key()
		key := string(keyBytes[3:]) // skip "KEY"
		for _, prefix := range prefixes {
			if strings.HasPrefix(key, prefix) || strings.HasPrefix("http://"+key, prefix) || strings.HasPrefix("https://"+key, prefix) {
				err := this.db.Delete(keyBytes, nil)
				if err != nil {
					return count, err
				}
				count++
				break
			}
		}
	}

	return count, nil
}

func (this *LevelDBManager) CleanExpired() error {
	if this.db == nil {
		return nil
	}
	it := this.db.NewIterator(util.BytesPrefix([]byte("KEY")), nil)
	for it.Next() {
		key := it.Key()
		data := it.Value()
		index := bytes.Index(data, []byte("_"))
		if index < -1 {
			if this.db != nil {
				this.db.Delete(key, nil)
			}
			continue
		}

		timestamp := types.Int64(data[:index])
		if timestamp < time.Now().Unix() {
			if this.db != nil {
				this.db.Delete(key, nil)
			}
			continue
		}
	}
	it.Release()
	return nil
}

// 统计
func (this *LevelDBManager) Stat() (size int64, countKeys int, err error) {
	if this.db == nil {
		return
	}
	it := this.db.NewIterator(util.BytesPrefix([]byte("KEY")), nil)
	for it.Next() {
		data := it.Value()
		countKeys++
		size += int64(len(data) + len(it.Key()))
	}
	it.Release()

	return
}

// 清理
func (this *LevelDBManager) Clean() error {
	if this.db == nil {
		return nil
	}
	it := this.db.NewIterator(util.BytesPrefix([]byte("KEY")), nil)
	for it.Next() {
		key := it.Key()

		if this.db != nil {
			this.db.Delete(key, nil)
		}
		continue
	}
	it.Release()
	return nil
}

// 关闭
func (this *LevelDBManager) Close() error {
	if this.looper != nil {
		this.looper.Stop()
		this.looper = nil
	}
	if this.db != nil {
		err := this.db.Close()
		//logs.Println("[cache]close cache policy instance: leveldb")

		this.db = nil
		return err
	}
	return nil
}
