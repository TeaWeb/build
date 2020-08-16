package teaproxy

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

// 公用的静态文件分发器
var ShareStaticDelivery = NewStaticDelivery()

// 静态文件分发器
type StaticDelivery struct {
	cacheMap map[string]*StaticFileCache // key => cache
	locker   sync.RWMutex
	capacity int
	life     int64
}

// 静态文件缓存
type StaticFileCache struct {
	path      string
	createdAt int64
	content   []byte
}

// 获取新对象
func NewStaticDelivery() *StaticDelivery {
	delivery := &StaticDelivery{
		capacity: 10000,
		life:     1800,
	}
	delivery.init()
	return delivery
}

// 初始化
func (this *StaticDelivery) init() {
	this.cacheMap = map[string]*StaticFileCache{}
}

// 读取
func (this *StaticDelivery) Read(path string, stat os.FileInfo) (reader io.Reader, shouldClose bool, err error) {
	if stat.Size() > 10*1024 { // <10K
		reader, err = os.OpenFile(path, os.O_RDONLY, 0444)
		shouldClose = true
		return
	}
	modifiedAt := stat.ModTime().Unix()
	key := path + "_" + fmt.Sprintf("%d_%d", stat.Size(), modifiedAt)

	this.locker.RLock()
	cache, found := this.cacheMap[key]
	if found {
		this.locker.RUnlock()
		return bytes.NewBuffer(cache.content), false, nil
	}
	this.locker.RUnlock()

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, false, err
	}

	this.locker.Lock()
	defer this.locker.Unlock()

	// 写入缓存
	reader = bytes.NewBuffer(content)
	timestamp := time.Now().Unix()
	this.cacheMap[key] = &StaticFileCache{
		path:      path,
		createdAt: timestamp,
		content:   content,
	}

	// 清理
	if len(this.cacheMap) > this.capacity {
		for k, v := range this.cacheMap {
			if timestamp-v.createdAt > this.life {
				delete(this.cacheMap, k)
			}
		}
	}

	return
}
