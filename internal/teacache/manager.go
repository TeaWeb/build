package teacache

import (
	"errors"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/utils/string"
	"regexp"
	"time"
)

var ErrNotFound = errors.New("cache not found")

// 缓存管理接口
type ManagerInterface interface {
	// 设置ID
	SetId(id string)

	// 写入
	Write(key string, data []byte) error

	// 读取
	Read(key string) (data []byte, err error)

	// 删除
	Delete(key string) error

	// 删除key前缀
	DeletePrefixes(prefixes []string) (int, error)

	// 设置选项
	SetOptions(options map[string]interface{})

	// 统计
	Stat() (size int64, countKeys int, err error)

	// 清理
	Clean() error

	// 关闭
	Close() error
}

type Manager struct {
	id string
}

// 设置ID
func (this *Manager) SetId(id string) {
	this.id = id
}

// 获取ID
func (this *Manager) Id() string {
	return this.id
}

// 获取新的管理对象
func NewManagerFromConfig(config *shared.CachePolicy) ManagerInterface {
	if len(config.Id) == 0 {
		filename := config.Filename
		match := regexp.MustCompile("^cache\\.policy\\.(\\w+)\\.conf$").FindStringSubmatch(filename)
		if len(match) > 0 {
			config.Id = match[1]
		}
	}

	switch config.Type {
	case "memory":
		m := NewMemoryManager()
		m.Life, _ = time.ParseDuration(config.Life)
		m.Capacity, _ = stringutil.ParseFileSize(config.Capacity)
		m.SetOptions(config.Options)
		m.SetId(config.Id)

		if len(config.Filename) > 0 {
			cachePolicyMap[config.Filename] = m
		}

		return m
	case "file":
		m := NewFileManager()
		m.Life, _ = time.ParseDuration(config.Life)
		m.Capacity, _ = stringutil.ParseFileSize(config.Capacity)
		m.SetOptions(config.Options)
		m.SetId(config.Id)
		return m
	case "redis":
		m := NewRedisManager()
		m.Life, _ = time.ParseDuration(config.Life)
		m.Capacity, _ = stringutil.ParseFileSize(config.Capacity)
		m.SetOptions(config.Options)
		m.SetId(config.Id)
		return m
	case "leveldb":
		m := NewLevelDBManager()
		m.Life, _ = time.ParseDuration(config.Life)
		m.Capacity, _ = stringutil.ParseFileSize(config.Capacity)
		m.SetOptions(config.Options)
		m.SetId(config.Id)
		return m
	}

	return nil
}
