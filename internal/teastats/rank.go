package teastats

import (
	"encoding/json"
	"fmt"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/types"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"sync"
)

// 排行计算
type Rank struct {
	size       int
	top        map[string]int
	buffer     map[string]int
	bufferSize int
	min        int
	minKey     string

	locker sync.RWMutex

	isLoading bool
	dirtyKeys map[string]int
}

// 排行值
type RankValue struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

// 获取新对象
func NewRank(topSize int, bufferSize int) *Rank {
	return &Rank{
		size:       topSize,
		bufferSize: bufferSize,
		top:        map[string]int{},
		buffer:     map[string]int{},
		dirtyKeys:  map[string]int{},
	}
}

// 添加键值
func (this *Rank) Add(key string) {
	this.locker.Lock()
	defer this.locker.Unlock()

	v, ok := this.top[key]
	if ok {
		newV := v + 1
		this.top[key] = newV

		if !this.isLoading {
			this.dirtyKeys[key] = newV
		}

		if this.min == 0 {
			this.min = newV
			this.minKey = key
		} else if this.minKey == key {
			this.min = 0
			for k, v := range this.top {
				if v < this.min {
					this.min = v
					this.minKey = k
				}
			}
		}

		return
	}

	v, ok = this.buffer[key]
	newV := v + 1

	if !this.isLoading {
		this.dirtyKeys[key] = newV
	}

	if newV >= this.min { // 移到top中
		if ok {
			delete(this.buffer, key)
		}
		this.top[key] = newV

		if len(this.top) > this.size {
			this.min = 0
			for k, v := range this.top {
				if v < this.min || this.min == 0 {
					this.min = v
					this.minKey = k
				}
			}

			// 将min key转到buffer中
			delete(this.top, this.minKey)
			if len(this.buffer) < this.bufferSize {
				this.buffer[this.minKey] = this.min
			}

			this.min = 0
			this.minKey = ""

			this.min = 0
			for k, v := range this.top {
				if v < this.min || this.min == 0 {
					this.min = v
					this.minKey = k
				}
			}
		} else if this.min == 0 {
			this.min = newV
			this.minKey = key
		}
	} else if len(this.top) < this.size {
		this.top[key] = newV
	} else {
		if ok || len(this.buffer) < this.bufferSize {
			this.buffer[key] = newV
		}
	}
}

// 获取排行结果
func (this *Rank) Top() []*RankValue {
	this.locker.RLock()
	defer this.locker.RUnlock()

	values := []*RankValue{}
	for k, v := range this.top {
		values = append(values, &RankValue{
			Key:   k,
			Value: v,
		})
	}
	lists.Sort(values, func(i int, j int) bool {
		return values[i].Value > values[j].Value
	})
	return values
}

// 保存到LevelDB
func (this *Rank) Save(db *leveldb.DB, prefix string) {
	tx, err := db.OpenTransaction()
	if err != nil {
		return
	}

	this.locker.Lock()
	keysCopy := this.dirtyKeys
	this.dirtyKeys = map[string]int{}
	topData, err := json.Marshal(this.top)
	this.locker.Unlock()

	for k, v := range keysCopy {
		err1 := tx.Put([]byte(prefix+"@item@"+k), []byte(fmt.Sprintf("%d", v)), nil)
		if err1 != nil {
			logs.Error(err1)
		}
	}

	if err != nil {
		logs.Error(err)
	} else {
		err1 := tx.Put([]byte(prefix+"@top"), topData, nil)
		if err1 != nil {
			logs.Error(err1)
		}
	}

	err = tx.Commit()
	if err != nil {
		logs.Error(err)
	}
}

// 从LevelDB中加载
func (this *Rank) Load(db *leveldb.DB, prefix string) {
	this.locker.Lock()
	defer this.locker.Unlock()

	this.isLoading = true

	top, err := db.Get([]byte(prefix+"@top"), nil)
	if err == nil && top != nil {
		m := map[string]int{}
		err = json.Unmarshal(top, &m)
		if err == nil {
			this.top = m

			for k, v := range m {
				if v < this.min || this.min == 0 {
					this.min = v
					this.minKey = k
				}
			}
		}
	}

	it := db.NewIterator(util.BytesPrefix([]byte(prefix+"@item@")), nil)
	prefixLen := len(prefix + "@item@")
	i := 0
	for it.Next() {
		k := string(it.Key()[prefixLen:])
		if _, ok := this.top[k]; ok {
			continue
		}

		if i >= this.bufferSize {
			break
		}

		v := types.Int(string(it.Value()))
		this.buffer[k] = v
		i++
	}
	it.Release()
	this.isLoading = false
}

// 重置
func (this *Rank) Reset() {
	this.locker.Lock()
	defer this.locker.Unlock()

	this.top = map[string]int{}
	this.buffer = map[string]int{}
	this.min = 0
	this.minKey = ""
	this.dirtyKeys = map[string]int{}
}
