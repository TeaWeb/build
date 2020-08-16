package tealogs

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/logs"
	"sync"
)

// 存储策略相关
var (
	policyMap       = map[string]*teaconfigs.AccessLogStoragePolicy{} // policyId => policy
	storageMap      = map[string]StorageInterface{}                   // policyId => StorageInterface
	storageNamesMap = map[string]string{}                             // policyId => policy name
	storageLocker   = sync.Mutex{}
)

// 通过策略ID查找策略
func FindPolicy(policyId string) *teaconfigs.AccessLogStoragePolicy {
	storageLocker.Lock()
	defer storageLocker.Unlock()

	policy, ok := policyMap[policyId]
	if ok {
		return policy
	}

	policy = teaconfigs.NewAccessLogStoragePolicyFromId(policyId)
	if policy != nil {
		err := policy.Validate()
		if err != nil {
			logs.Error(err)
		}
	}
	policyMap[policyId] = policy
	return policy
}

// 通过策略ID查找存储
func FindPolicyStorage(policyId string) StorageInterface {
	storageLocker.Lock()
	defer storageLocker.Unlock()

	storage, ok := storageMap[policyId]
	if ok {
		return storage
	}

	policy := teaconfigs.NewAccessLogStoragePolicyFromId(policyId)
	if policy == nil || !policy.On {
		storageMap[policyId] = nil
		return nil
	}

	storageNamesMap[policyId] = policy.Name

	storage = DecodePolicyStorage(policy)
	if storage != nil {
		err := storage.Start()
		if err != nil {
			logs.Println("access log storage '"+policyId+"/"+policy.Name+"' start failed:", err.Error())
			storage = nil
		}
	}
	storageMap[policyId] = storage

	return storage
}

// 清除策略相关信息
func ResetPolicyStorage(policyId string) {
	storageLocker.Lock()
	storage, ok := storageMap[policyId]
	if ok {
		delete(policyMap, policyId)
		delete(storageMap, policyId)
		delete(storageNamesMap, policyId)
	}
	storageLocker.Unlock()

	if storage != nil {
		_ = storage.Close()
	}
}

// 清除所有策略相关信息
func ResetAllPolicies() {
	storageLocker.Lock()
	policyMap = map[string]*teaconfigs.AccessLogStoragePolicy{}
	storageMap = map[string]StorageInterface{}
	storageNamesMap = map[string]string{}
	storageLocker.Unlock()
}

// 解析策略中的存储对象
func DecodePolicyStorage(policy *teaconfigs.AccessLogStoragePolicy) StorageInterface {
	if policy == nil {
		return nil
	}

	var instance StorageInterface = nil
	switch policy.Type {
	case StorageTypeFile:
		instance = new(FileStorage)
	case StorageTypeES:
		instance = new(ESStorage)
	case StorageTypeMySQL:
		instance = new(MySQLStorage)
	case StorageTypeTCP:
		instance = new(TCPStorage)
	case StorageTypeSyslog:
		instance = new(SyslogStorage)
	case StorageTypeCommand:
		instance = new(CommandStorage)
	}
	if instance == nil {
		return nil
	}

	err := teautils.MapToObjectJSON(policy.Options, instance)
	if err != nil {
		logs.Error(err)
	}

	return instance
}

// 查找策略名称
func FindPolicyName(policyId string) string {
	storageLocker.Lock()
	name, _ := storageNamesMap[policyId]
	storageLocker.Unlock()

	return name
}
