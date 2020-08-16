package teastats

import (
	"fmt"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/timers"
	"github.com/iwind/TeaGo/types"
	"github.com/syndtr/goleveldb/leveldb"
	"strings"
	"time"
)

// KV存储引擎
type KVStorage struct {
	db *leveldb.DB
}

var sharedKV *KVStorage = nil

// 获取新对象
func NewKVStorage(filename string) *KVStorage {
	db, err := leveldb.OpenFile(Tea.LogFile(filename), nil)
	if err != nil {
		logs.Error(err)
		return nil
	}
	kv := &KVStorage{
		db: db,
	}
	go func() {
		timers.Loop(24*time.Hour, func(looper *timers.Looper) {
			kv.compact()
		})
	}()
	return kv
}

// 设置键值对
func (this *KVStorage) Set(key string, value string, life time.Duration) error {
	// key: [life end timestamp]_[value]
	timestamp := time.Now().Unix() + int64(life.Seconds())
	return this.db.Put([]byte(key), []byte(fmt.Sprintf("%d", timestamp)+"_"+value), nil)
}

// 获取键对应的值
func (this *KVStorage) Get(key string) (string, error) {
	value, err := this.db.Get([]byte(key), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return "", nil
		}
		return "", err
	}
	valueString := string(value)
	index := strings.Index(valueString, "_")
	if index < 0 {
		return valueString, nil
	}
	timestamp := types.Int64(valueString[:index])
	if timestamp < time.Now().Unix() {
		return "", nil
	}
	return valueString[index+1:], nil
}

// 检查key是否存在，并不检查是否过期
func (this *KVStorage) Has(key string) (bool, error) {
	return this.db.Has([]byte(key), nil)
}

// 关闭
func (this *KVStorage) Close() error {
	return this.db.Close()
}

// 清理
func (this *KVStorage) compact() {
	// TODO 需要实现
}
