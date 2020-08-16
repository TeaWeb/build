package api

// API限制
type APILimit struct {
	Concurrent    uint               `yaml:"concurrent" json:"concurrent"` // 并发数
	RequestLimits []*APIRequestLimit `yaml:"request" json:"request"`       // 请求数限制 TODO
	DataLimits    []*APIDataLimit    `yaml:"data" json:"data"`             // 数据量限制 TODO

	concurrentChan     chan bool // 并发管道
	concurrentIsLocked bool      // 是否正在锁中
}

// 获取新的对象
func NewAPILimit() *APILimit {
	return &APILimit{}
}

// 校验
func (this *APILimit) Validate() error {
	// request
	for _, reqLimit := range this.RequestLimits {
		err := reqLimit.Validate()
		if err != nil {
			return err
		}
	}

	// data
	for _, dataLimit := range this.DataLimits {
		err := dataLimit.Validate()
		if err != nil {
			return err
		}
	}

	// 管道
	if this.Concurrent > 0 {
		this.concurrentChan = make(chan bool, this.Concurrent)
	}

	return nil
}

// 等待并发
func (this *APILimit) Begin() {
	// concurrent
	if this.Concurrent > 0 {
		this.concurrentChan <- true
	}

	// request

	// data
}

func (this *APILimit) Done() {
	// concurrent
	if this.Concurrent > 0 {
		<-this.concurrentChan
	}

	// request

	// data
}
