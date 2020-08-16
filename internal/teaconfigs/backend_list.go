package teaconfigs

import (
	"github.com/TeaWeb/build/internal/teaconfigs/scheduling"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/lists"
	"sync"
)

// BackendList接口定义
type BackendListInterface interface {
	// 校验
	ValidateBackends() error

	// 添加Backend
	AddBackend(backend *BackendConfig)

	// 所有的Backends
	AllBackends() []*BackendConfig

	// 删除后端服务器
	DeleteBackend(backendId string)

	// 查找后端服务器
	FindBackend(backendId string) *BackendConfig

	// 调度算法
	SchedulingConfig() *SchedulingConfig

	// 设置调度算法
	SetSchedulingConfig(scheduling *SchedulingConfig)

	// 重新建立调度
	SetupScheduling(isBackup bool)
}

// BackendList定义
type BackendList struct {
	Backends    []*BackendConfig  `yaml:"backends" json:"backends"`
	Scheduling  *SchedulingConfig `yaml:"scheduling" json:"scheduling"` // 调度算法选项
	hasBackends bool

	schedulingIsBackup bool
	schedulingObject   scheduling.SchedulingInterface
	schedulingLocker   sync.Mutex
}

// 校验
func (this *BackendList) ValidateBackends() error {
	this.hasBackends = len(this.Backends) > 0

	for _, backend := range this.Backends {
		err := backend.Validate()
		if err != nil {
			return err
		}
	}

	// scheduling
	this.SetupScheduling(false)

	return nil
}

// 添加Backend
func (this *BackendList) AddBackend(backend *BackendConfig) {
	this.Backends = append(this.Backends, backend)
}

// 所有的Backends
func (this *BackendList) AllBackends() []*BackendConfig {
	return this.Backends
}

// 删除后端服务器
func (this *BackendList) DeleteBackend(backendId string) {
	result := []*BackendConfig{}
	for _, backend := range this.Backends {
		if backend.Id == backendId {
			continue
		}
		result = append(result, backend)
	}
	this.Backends = result
}

// 删除一组后端服务器
func (this *BackendList) DeleteBackends(backendIds []string) {
	if len(backendIds) == 0 {
		return
	}
	result := []*BackendConfig{}
	for _, backend := range this.Backends {
		if lists.ContainsString(backendIds, backend.Id) {
			continue
		}
		result = append(result, backend)
	}
	this.Backends = result
}

// 根据ID查找后端服务器
func (this *BackendList) FindBackend(backendId string) *BackendConfig {
	for _, backend := range this.Backends {
		if backend.Id == backendId {
			return backend
		}
	}
	return nil
}

// 取得下一个可用的后端服务
func (this *BackendList) NextBackend(call *shared.RequestCall) *BackendConfig {
	this.schedulingLocker.Lock()
	defer this.schedulingLocker.Unlock()

	if this.schedulingObject == nil {
		return nil
	}

	if this.Scheduling != nil && call != nil && call.Options != nil {
		for k, v := range this.Scheduling.Options {
			call.Options[k] = v
		}
	}

	candidate := this.schedulingObject.Next(call)
	if candidate == nil {
		// 启用备用服务器
		if !this.schedulingIsBackup {
			this.SetupScheduling(true)

			candidate = this.schedulingObject.Next(call)
			if candidate == nil {
				return nil
			}
		}

		if candidate == nil {
			return nil
		}
	}

	return candidate.(*BackendConfig)
}

// 设置调度算法
func (this *BackendList) SetupScheduling(isBackup bool) {
	if !isBackup {
		this.schedulingLocker.Lock()
		defer this.schedulingLocker.Unlock()
	}
	this.schedulingIsBackup = isBackup

	if this.Scheduling == nil {
		this.schedulingObject = &scheduling.RandomScheduling{}
	} else {
		typeCode := this.Scheduling.Code
		s := scheduling.FindSchedulingType(typeCode)
		if s == nil {
			this.Scheduling = nil
			this.schedulingObject = &scheduling.RandomScheduling{}
		} else {
			this.schedulingObject = s["instance"].(scheduling.SchedulingInterface)
		}
	}

	for _, backend := range this.Backends {
		if backend.On && !backend.IsDown {
			if isBackup && backend.IsBackup {
				this.schedulingObject.Add(backend)
			} else if !isBackup && !backend.IsBackup {
				this.schedulingObject.Add(backend)
			}
		}
	}

	this.schedulingObject.Start()
}

// 调度算法
func (this *BackendList) SchedulingConfig() *SchedulingConfig {
	return this.Scheduling
}

// 设置调度算法
func (this *BackendList) SetSchedulingConfig(scheduling *SchedulingConfig) {
	this.Scheduling = scheduling
}

// 判断是否有后端服务器
func (this *BackendList) HasBackends() bool {
	return this.hasBackends
}

// 克隆
func (this *BackendList) CloneBackendList() *BackendList {
	newBackendList := new(BackendList)
	newBackendList.Backends = this.Backends
	newBackendList.Scheduling = this.Scheduling
	newBackendList.hasBackends = this.hasBackends
	newBackendList.SetupScheduling(false)
	return newBackendList
}
