package scheduling

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/maps"
	"math"
	"math/rand"
	"time"
)

// 随机调度算法
type RandomScheduling struct {
	Scheduling

	array []CandidateInterface
	count uint // 实际总的服务器数
}

// 启动
func (this *RandomScheduling) Start() {
	sumWeight := uint(0)
	for _, c := range this.Candidates {
		weight := c.CandidateWeight()
		if weight == 0 {
			weight = 1
		} else if weight > 10000 {
			weight = 10000
		}
		sumWeight += weight
	}

	if sumWeight == 0 {
		return
	}

	for _, c := range this.Candidates {
		weight := c.CandidateWeight()
		if weight == 0 {
			weight = 1
		} else if weight > 10000 {
			weight = 10000
		}
		count := uint(0)
		if sumWeight <= 1000 {
			count = weight
		} else {
			count = uint(math.Round(float64(weight*10000) / float64(sumWeight))) // 1% 产生 100个数据，最多支持10000个服务器
		}
		for i := uint(0); i < count; i++ {
			this.array = append(this.array, c)
		}
		this.count += count
	}

	rand.Seed(time.Now().UnixNano())
}

// 获取下一个候选对象
func (this *RandomScheduling) Next(call *shared.RequestCall) CandidateInterface {
	if this.count == 0 {
		return nil
	}
	if this.count == 1 {
		return this.array[0]
	}
	index := rand.Int() % int(this.count)
	return this.array[index]
}

// 获取简要信息
func (this *RandomScheduling) Summary() maps.Map {
	return maps.Map{
		"code":        "random",
		"name":        "Random随机算法",
		"description": "根据权重设置随机分配后端服务器",
		"networks":    []string{"http", "tcp"},
	}
}
