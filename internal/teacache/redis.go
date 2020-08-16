package teacache

import (
	"fmt"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
	"strings"
	"time"
)

// 内存缓存管理器
type RedisManager struct {
	Manager

	Capacity float64       // 容量
	Life     time.Duration // 有效期

	Network  string
	Host     string
	Port     int
	Password string
	Sock     string

	client *redis.Client
}

func NewRedisManager() *RedisManager {
	m := &RedisManager{}
	return m
}

func (this *RedisManager) SetOptions(options map[string]interface{}) {
	if this.Life <= 0 {
		this.Life = 1800 * time.Second
	}

	m := maps.NewMap(options)
	this.Network = m.GetString("network")
	this.Host = m.GetString("host")
	this.Port = m.GetInt("port")
	this.Password = m.GetString("password")
	this.Sock = m.GetString("sock")

	if len(this.Network) == 0 {
		this.Network = "tcp"
	}

	addr := ""
	if this.Network == "tcp" {
		if this.Port > 0 {
			addr = fmt.Sprintf("%s:%d", this.Host, this.Port)
		} else {
			addr = this.Host + ":6379"
		}
	} else if this.Network == "sock" {
		addr = this.Sock
	}

	if this.client != nil {
		_ = this.client.Close()
	}

	this.client = redis.NewClient(&redis.Options{
		Network:      this.Network,
		Addr:         addr,
		Password:     this.Password,
		DialTimeout:  10 * time.Second, // TODO 换成可配置
		ReadTimeout:  10 * time.Second, // TODO 换成可配置
		WriteTimeout: 10 * time.Second, // TODO 换成可配置
		TLSConfig:    nil,              // TODO 支持TLS
	})
}

func (this *RedisManager) Write(key string, data []byte) error {
	cmd := this.client.Set(context.Background(), "TEA_CACHE_"+this.id+key, string(data), this.Life)
	return cmd.Err()
}

func (this *RedisManager) Read(key string) (data []byte, err error) {
	cmd := this.client.Get(context.Background(), "TEA_CACHE_"+this.id+key)
	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			return nil, ErrNotFound
		}
		logs.Printf("%#v", cmd.Err())
		return nil, cmd.Err()
	}
	return []byte(cmd.Val()), nil
}

// 删除
func (this *RedisManager) Delete(key string) error {
	cmd := this.client.Del(context.Background(), "TEA_CACHE_"+this.id+key)
	return cmd.Err()
}

// 删除key前缀
func (this *RedisManager) DeletePrefixes(prefixes []string) (int, error) {
	if len(prefixes) == 0 {
		return 0, nil
	}

	cursor := uint64(0)
	var err error
	loopCount := 0
	count := 0
	keyPrefix := "TEA_CACHE_" + this.Id()
	keyPrefixLength := len(keyPrefix)
	for {
		loopCount++

		var keys []string
		keys, cursor, err = this.client.Scan(context.Background(), cursor, keyPrefix+"*", 10000).Result()
		if err != nil {
			return count, err
		}
		if len(keys) > 0 {
			for _, key := range keys {
				realKey := key[keyPrefixLength:]
				for _, prefix := range prefixes {
					if strings.HasPrefix(realKey, prefix) || strings.HasPrefix("http://"+realKey, prefix) || strings.HasPrefix("https://"+realKey, prefix) {
						err1 := this.client.Del(context.Background(), key).Err()
						if err1 != nil {
							err = err1
							break
						}
						count++
						break
					}
				}
			}
		}

		// 防止单个操作时间过长
		if loopCount > 10000 {
			break
		}

		if cursor == 0 {
			break
		}
	}

	return count, nil
}

// 统计
func (this *RedisManager) Stat() (size int64, countKeys int, err error) {
	cursor := uint64(0)
	loopCount := 0
	for {
		loopCount++

		var keys []string
		keys, cursor, err = this.client.Scan(context.Background(), cursor, "TEA_CACHE_"+this.Id()+"*", 10000).Result()
		if err != nil {
			return
		}
		if len(keys) > 0 {
			countKeys += len(keys)
			for _, key := range keys {
				val, _ := this.client.Get(context.Background(), key).Bytes()
				size += int64(len(val))
			}
		}

		// 防止单个操作时间过长
		if loopCount > 10000 {
			break
		}

		if cursor == 0 {
			break
		}
	}
	return
}

// 清理
func (this *RedisManager) Clean() error {
	cursor := uint64(0)
	var err error
	loopCount := 0
	for {
		loopCount++

		var keys []string
		keys, cursor, err = this.client.Scan(context.Background(), cursor, "TEA_CACHE_"+this.Id()+"*", 10000).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			for _, key := range keys {
				err1 := this.client.Del(context.Background(), key).Err()
				if err1 != nil {
					err = err1
					break
				}
			}
		}

		// 防止单个操作时间过长
		if loopCount > 10000 {
			break
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

// 关闭
func (this *RedisManager) Close() error {
	if this.client != nil {
		//logs.Println("[cache]close cache policy instance: redis")

		err := this.client.Close()
		this.client = nil

		return err
	}

	return nil
}
